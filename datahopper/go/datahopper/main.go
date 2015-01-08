/*
	Pulls data from multiple sources and funnels into InfluxDB.
*/

package main

import (
	"flag"
	"path"
	"time"

	influxdb "github.com/influxdb/influxdb/client"
	"github.com/rcrowley/go-metrics"
	"github.com/skia-dev/glog"
	"skia.googlesource.com/buildbot.git/datahopper/go/autoroll_ingest"
	"skia.googlesource.com/buildbot.git/go/buildbot"
	"skia.googlesource.com/buildbot.git/go/common"
	"skia.googlesource.com/buildbot.git/go/database"
	"skia.googlesource.com/buildbot.git/go/gitinfo"
	"skia.googlesource.com/buildbot.git/go/metadata"
)

const (
	SKIA_REPO = "https://skia.googlesource.com/skia"
)

// flags
var (
	graphiteServer   = flag.String("graphite_server", "localhost:2003", "Where is Graphite metrics ingestion server running.")
	influxDbHost     = flag.String("influxdb_host", "localhost:8086", "The InfluxDB hostname.")
	influxDbName     = flag.String("influxdb_name", "root", "The InfluxDB username.")
	influxDbPassword = flag.String("influxdb_password", "root", "The InfluxDB password.")
	influxDbDatabase = flag.String("influxdb_database", "", "The InfluxDB database.")
	workdir          = flag.String("workdir", ".", "Working directory used by data processors.")
	local            = flag.Bool("local", false, "Running locally if true. As opposed to in production.")
)

func main() {
	// Setup DB flags.
	database.SetupFlags(buildbot.PROD_DB_HOST, buildbot.PROD_DB_PORT, database.USER_RW, buildbot.PROD_DB_NAME)

	// Global init to initialize glog and parse arguments.
	common.InitWithMetrics("datahopper", graphiteServer)

	// Prepare the InfluxDB credentials. Load from metadata if appropriate.
	if !*local {
		*influxDbName = metadata.Must(metadata.ProjectGet(metadata.INFLUXDB_NAME))
		*influxDbPassword = metadata.Must(metadata.ProjectGet(metadata.INFLUXDB_PASSWORD))
	}
	dbClient, err := influxdb.New(&influxdb.ClientConfig{
		Host:       *influxDbHost,
		Username:   *influxDbName,
		Password:   *influxDbPassword,
		Database:   *influxDbDatabase,
		HttpClient: nil,
		IsSecure:   false,
		IsUDP:      false,
	})
	if err != nil {
		glog.Fatalf("Failed to initialize InfluxDB client: %s", err)
	}

	// Data generation goroutines.
	go autoroll_ingest.LoadAutoRollData(dbClient, *workdir)
	go func() {
		// Initialize the buildbot database.
		conf, err := database.ConfigFromFlagsAndMetadata(*local, buildbot.MigrationSteps())
		if err := buildbot.InitDB(conf); err != nil {
			glog.Fatal(err)
		}
		// Create the Git repo.
		skiaRepo, err := gitinfo.CloneOrUpdate(SKIA_REPO, path.Join(*workdir, "buildbot_git", "skia"), true)
		if err != nil {
			glog.Fatal(err)
		}
		// Ingest data in a loop.
		for _ = range time.Tick(30 * time.Second) {
			skiaRepo.Update(true, true)
			glog.Info("Ingesting builds.")
			if err := buildbot.IngestNewBuilds(skiaRepo); err != nil {
				glog.Errorf("Failed to ingest new builds: %v", err)
			}
		}
	}()
	// Measure buildbot data ingestion progress.
	totalGuage := metrics.GetOrRegisterGauge("buildbot.builds.total", metrics.DefaultRegistry)
	ingestGuage := metrics.GetOrRegisterGauge("buildbot.builds.ingested", metrics.DefaultRegistry)
	go func() {
		for _ = range time.Tick(common.SAMPLE_PERIOD) {
			totalBuilds, err := buildbot.NumTotalBuilds()
			if err != nil {
				glog.Error(err)
				continue
			}
			ingestedBuilds, err := buildbot.NumIngestedBuilds()
			if err != nil {
				glog.Error(err)
				continue
			}
			totalGuage.Update(int64(totalBuilds))
			ingestGuage.Update(int64(ingestedBuilds))
		}
	}()

	// Wait while the above goroutines generate data.
	select {}
}
