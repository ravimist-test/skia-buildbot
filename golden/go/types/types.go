package types

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"

	"go.skia.org/infra/go/paramtools"
	"go.skia.org/infra/go/tiling"
)

func init() {
	// Register *GoldenTrace in gob so that it can be used as a
	// concrete type for Trace when writing and reading Tiles in gobs.
	gob.Register(&GoldenTrace{})
}

const (
	// Primary key field that uniquely identifies a key.
	PRIMARY_KEY_FIELD = "name"

	// Field that contains the corpus identifier.
	CORPUS_FIELD = "source_type"

	// MAXIMUM_NAME_LENGTH is the maximum length in bytes a test name can be.
	MAXIMUM_NAME_LENGTH = 256
)

// Label for classifying digests.
type Label int

const (
	// Classifications for observed digests.
	UNTRIAGED Label = iota // == 0
	POSITIVE
	NEGATIVE
)

// String representation for Labels. The order must match order above.
var labelStringRepresentation = []string{
	"untriaged",
	"positive",
	"negative",
}

func (l Label) String() string {
	return labelStringRepresentation[l]
}

var labels = map[string]Label{
	"untriaged": UNTRIAGED,
	"positive":  POSITIVE,
	"negative":  NEGATIVE,
}

func LabelFromString(s string) Label {
	if l, ok := labels[s]; ok {
		return l
	}
	return UNTRIAGED
}

// ValidLabel returns true if the given label is a valid label string.
func ValidLabel(s string) bool {
	_, ok := labels[s]
	return ok
}

// Stores the digests and their associated labels.
// Note: The name of the test is assumed to be handled by the client of this
// type. Most likely in the keys of a map.
type TestClassification map[string]Label

func (tc *TestClassification) DeepCopy() TestClassification {
	result := make(map[string]Label, len(*tc))
	for k, v := range *tc {
		result[k] = v
	}
	return result
}

// ComplexTile contains an enriched version of a tile loaded through the ingestion process.
// It provides ways to handle sparse tiles, where many commits of the underlying raw tile
// contain no data and therefore removed.
// In either case (sparse or dense tile) it offers two versions of the tile.
// one with all ignored traces and one without the ignored traces.
// In addition it also contains the ignore rules and information about the larger "sparse" tile
// if the tiles at hand were condensed from a larger tile.
type ComplexTile struct {
	// tile is the current tile without ignored traces.
	tile *tiling.Tile

	// tileWithIgnores is the current tile containing all available data.
	tileWithIgnores *tiling.Tile

	// ignoreRules contains the rules used to created the TileWithIgnores.
	ignoreRules paramtools.ParamMatcher

	// irRevision is the (monotonically increasing) revision of the ignore rules.
	irRevision int64

	// sparseCommits are all the commits that were used condense the underlying tile.
	sparseCommits []*tiling.Commit

	// cards captures the cardinality of each commit in sparse tile, meaning how many data points
	// each commit contains.
	cardinalities []int

	// filled contains the number of commits that are non-empty.
	filled int
}

func NewComplexTile(completeTile *tiling.Tile) *ComplexTile {
	return &ComplexTile{
		tile:            completeTile,
		tileWithIgnores: completeTile,
	}
}

// SetIgnoreRules adds ignore rules to the tile and a sub-tile with the ignores removed.
// In other words this function assumes that original tile has been filtered by the
// ignore rules that are being passed.
func (c *ComplexTile) SetIgnoreRules(reducedTile *tiling.Tile, ignoreRules paramtools.ParamMatcher, irRev int64) *ComplexTile {
	c.tile = reducedTile
	c.irRevision = irRev
	c.ignoreRules = ignoreRules
	return c
}

// SetSparse sets sparsity information about this tile.
func (c *ComplexTile) SetSparse(sparseCommits []*tiling.Commit, cardinalities []int) *ComplexTile {
	// Make sure we always have valid values sparce commits.
	if len(sparseCommits) == 0 {
		sparseCommits = c.tileWithIgnores.Commits
	}

	filled := len(c.tileWithIgnores.Commits)
	if len(cardinalities) == 0 {
		cardinalities = make([]int, len(sparseCommits))
		for idx := range cardinalities {
			cardinalities[idx] = len(c.tileWithIgnores.Traces)
		}
	} else {
		for _, card := range cardinalities {
			if card > 0 {
				filled++
			}
		}
	}

	commitsLen := tiling.LastCommitIndex(sparseCommits) + 1
	if commitsLen < len(sparseCommits) {
		sparseCommits = sparseCommits[:commitsLen]
		cardinalities = cardinalities[:commitsLen]
	}
	c.sparseCommits = sparseCommits
	c.cardinalities = cardinalities
	c.filled = filled
	return c
}

func (c *ComplexTile) FilledCommits() int {
	return c.filled
}

// ensureSparseInfo is a helper function that fills in the sparsity information if it wasn't set.
func (c *ComplexTile) ensureSparseInfo() {
	if c != nil {
		if c.sparseCommits == nil || c.cardinalities == nil {
			c.SetSparse(nil, nil)
		}
	}
}

// FromSame returns true if the given complex tile was derived from the same tile as the one
// provided and if none of the other parameters changed, especially the ignore revision.
func (c *ComplexTile) FromSame(completeTile *tiling.Tile, ignoreRev int64) bool {
	return c != nil &&
		c.tileWithIgnores != nil &&
		c.tileWithIgnores == completeTile &&
		c.tile != nil &&
		c.irRevision == ignoreRev
}

// DataCommits returns all commits that contain data. In some busy repos, there are commits that
// don't get tested directly because the commits are batched in with others. DataCommits
// is a way to get just the commits where some data has been ingested.
func (c *ComplexTile) DataCommits() []*tiling.Commit {
	return c.tileWithIgnores.Commits
}

// AllCommits returns all commits that were processed to get the data commits.
// Its first commit should match the first commit returned when calling DataCommits.
func (c *ComplexTile) AllCommits() []*tiling.Commit {
	return c.sparseCommits
}

// GetTile returns a simple tile either with or without ignored traces depending on the argument.
func (c *ComplexTile) GetTile(includeIgnores bool) *tiling.Tile {
	if includeIgnores {
		return c.tileWithIgnores
	}
	return c.tile
}

// IgnoreRules returns the ignore rules for this tile.
func (c *ComplexTile) IgnoreRules() paramtools.ParamMatcher {
	return c.ignoreRules
}

const (
	// No digest available.
	MISSING_DIGEST = ""
)

// GoldenTrace represents all the Digests of a single test across a series
// of Commits. GoldenTrace implements the Trace interface.
type GoldenTrace struct {
	// The JSON keys are named this way to maintain backwards compatibility
	// with JSON already written to disk.
	Keys    map[string]string `json:"Params_"`
	Digests []string          `json:"Values"`
}

func (g *GoldenTrace) Params() map[string]string {
	return g.Keys
}

func (g *GoldenTrace) Len() int {
	return len(g.Digests)
}

func (g *GoldenTrace) IsMissing(i int) bool {
	return g.Digests[i] == MISSING_DIGEST
}

func (g *GoldenTrace) DeepCopy() tiling.Trace {
	n := len(g.Digests)
	cp := &GoldenTrace{
		Digests: make([]string, n, n),
		Keys:    make(map[string]string),
	}
	copy(cp.Digests, g.Digests)
	for k, v := range g.Keys {
		cp.Keys[k] = v
	}
	return cp
}

func (g *GoldenTrace) Merge(next tiling.Trace) tiling.Trace {
	nextGold := next.(*GoldenTrace)
	n := len(g.Digests) + len(nextGold.Digests)
	n1 := len(g.Digests)

	merged := NewGoldenTraceN(n)
	merged.Keys = g.Keys
	for k, v := range nextGold.Keys {
		merged.Keys[k] = v
	}
	for i, v := range g.Digests {
		merged.Digests[i] = v
	}
	for i, v := range nextGold.Digests {
		merged.Digests[n1+i] = v
	}
	return merged
}

func (g *GoldenTrace) Grow(n int, fill tiling.FillType) {
	if n < len(g.Digests) {
		panic(fmt.Sprintf("Grow must take a value (%d) larger than the current Trace size: %d", n, len(g.Digests)))
	}
	delta := n - len(g.Digests)
	newDigests := make([]string, n)

	if fill == tiling.FILL_AFTER {
		copy(newDigests, g.Digests)
		for i := 0; i < delta; i++ {
			newDigests[i+len(g.Digests)] = MISSING_DIGEST
		}
	} else {
		for i := 0; i < delta; i++ {
			newDigests[i] = MISSING_DIGEST
		}
		copy(newDigests[delta:], g.Digests)
	}
	g.Digests = newDigests
}

func (g *GoldenTrace) Trim(begin, end int) error {
	if end < begin || end > g.Len() || begin < 0 {
		return fmt.Errorf("Invalid Trim range [%d, %d) of [0, %d]", begin, end, g.Len())
	}
	n := end - begin
	newDigests := make([]string, n)

	for i := 0; i < n; i++ {
		newDigests[i] = g.Digests[i+begin]
	}
	g.Digests = newDigests
	return nil
}

func (g *GoldenTrace) SetAt(index int, value []byte) error {
	if index < 0 || index > len(g.Digests) {
		return fmt.Errorf("Invalid index: %d", index)
	}
	g.Digests[index] = string(value)
	return nil
}

// LastDigest returns the last digest in the trace (HEAD) or the empty string otherwise.
func (g *GoldenTrace) LastDigest() string {
	if idx := g.LastIndex(); idx >= 0 {
		return g.Digests[idx]
	}
	return ""
}

// LastIndex returns the index of last non-empty value in this trace and -1 if
// if the entire trace is empty.
func (g *GoldenTrace) LastIndex() int {
	for i := len(g.Digests) - 1; i >= 0; i-- {
		if g.Digests[i] != MISSING_DIGEST {
			return i
		}
	}
	return -1
}

// NewGoldenTrace allocates a new Trace set up for the given number of samples.
//
// The Trace Digests are pre-filled in with the missing data sentinel since not
// all tests will be run on all commits.
func NewGoldenTrace() *GoldenTrace {
	return NewGoldenTraceN(tiling.TILE_SIZE)
}

// NewGoldenTraceN allocates a new Trace set up for the given number of samples.
//
// The Trace Digests are pre-filled in with the missing data sentinel since not
// all tests will be run on all commits.
func NewGoldenTraceN(n int) *GoldenTrace {
	g := &GoldenTrace{
		Digests: make([]string, n, n),
		Keys:    make(map[string]string),
	}
	for i := range g.Digests {
		g.Digests[i] = MISSING_DIGEST
	}
	return g
}

func GoldenTraceBuilder(n int) tiling.Trace {
	return NewGoldenTraceN(n)
}

// Same as Tile but instead of Traces we preserve the raw JSON. This is a
// utility struct that is used to parse a tile where we don't know the
// Trace type upfront.
type TileWithRawTraces struct {
	Traces    map[string]json.RawMessage `json:"traces"`
	ParamSet  map[string][]string        `json:"param_set"`
	Commits   []*tiling.Commit           `json:"commits"`
	Scale     int                        `json:"scale"`
	TileIndex int                        `json:"tileIndex"`
}

// TileFromJson parses a tile that has been serialized to JSON.
// traceExample has to be an instance of the Trace implementation
// that needs to be deserialized.
// Note: Instead of the type switch below we could use reflection
// to be truly generic, but it makes the code harder to read and
// currently we only have two types.
func TileFromJson(r io.Reader, traceExample tiling.Trace) (*tiling.Tile, error) {
	factory := func() tiling.Trace { return NewGoldenTrace() }

	// Decode everything, but the traces.
	dec := json.NewDecoder(r)
	var rawTile TileWithRawTraces
	err := dec.Decode(&rawTile)
	if err != nil {
		return nil, err
	}

	// Parse the traces.
	traces := map[string]tiling.Trace{}
	for k, rawJson := range rawTile.Traces {
		newTrace := factory()
		if err = json.Unmarshal(rawJson, newTrace); err != nil {
			return nil, err
		}
		traces[k] = newTrace.(tiling.Trace)
	}

	return &tiling.Tile{
		Traces:    traces,
		ParamSet:  rawTile.ParamSet,
		Commits:   rawTile.Commits,
		Scale:     rawTile.Scale,
		TileIndex: rawTile.Scale,
	}, nil
}
