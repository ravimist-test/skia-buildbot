package diffstore

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"go.skia.org/infra/go/skerr"
	"go.skia.org/infra/go/sklog"
	"go.skia.org/infra/go/util"
	"go.skia.org/infra/golden/go/diff"
	"go.skia.org/infra/golden/go/diffstore/common"
	"go.skia.org/infra/golden/go/types"
	"google.golang.org/grpc"
)

// NetDiffStore implements the DiffStore interface and wraps around
// a DiffService client.
type NetDiffStore struct {
	// serviceClient is the gRPC client for the DiffService.
	serviceClient DiffServiceClient

	// diffServerImageAddress is the port where the diff server serves images.
	diffServerImageAddress string

	// codec is used to decode the byte array received from the diff server into
	// a diff metrics map
	codec util.Codec
}

// NewNetDiffStore implements the diff.DiffStore interface via the gRPC-based DiffService.
func NewNetDiffStore(conn *grpc.ClientConn, diffServerImageAddress string) (diff.DiffStore, error) {
	serviceClient := NewDiffServiceClient(conn)
	if _, err := serviceClient.Ping(context.Background(), &Empty{}); err != nil {
		return nil, fmt.Errorf("Could not ping over connection: %s", err)
	}

	return &NetDiffStore{
		serviceClient:          serviceClient,
		diffServerImageAddress: diffServerImageAddress,
		codec:                  util.NewJSONCodec(map[types.Digest]*diff.DiffMetrics{}),
	}, nil
}

// Get, see the diff.DiffStore interface.
func (n *NetDiffStore) Get(ctx context.Context, priority int64, mainDigest types.Digest, rightDigests types.DigestSlice) (map[types.Digest]*diff.DiffMetrics, error) {
	req := &GetDiffsRequest{Priority: priority, MainDigest: string(mainDigest), RightDigests: common.AsStrings(rightDigests)}
	resp, err := n.serviceClient.GetDiffs(ctx, req)
	if err != nil {
		return nil, skerr.Wrap(err)
	}

	data, err := n.codec.Decode(resp.Diffs)
	if err != nil {
		return nil, skerr.Wrapf(err, "Could not decode response")
	}

	diffMetrics := data.(map[types.Digest]*diff.DiffMetrics)
	return diffMetrics, nil
}

// ImageHandler, see the diff.DiffStore interface. The images are expected to be served by the
// the server that implements the backend of the DiffService.
func (n *NetDiffStore) ImageHandler(urlPrefix string) (http.Handler, error) {
	// Set up a proxy to the diffserver images ports.
	// With ingress rules, it would be possible to serve directly to the diffserver
	// and bypass the main server.
	targetURL, err := url.Parse(fmt.Sprintf("http://%s", n.diffServerImageAddress))
	if err != nil {
		return nil, skerr.Wrapf(err, "invalid URL for serving diff images")
	}
	return httputil.NewSingleHostReverseProxy(targetURL), nil
}

// WarmDigests, see the diff.DiffStore interface.
func (n *NetDiffStore) WarmDigests(ctx context.Context, priority int64, digests types.DigestSlice, sync bool) {
	req := &WarmDigestsRequest{Priority: priority, Digests: common.AsStrings(digests), Sync: sync}
	_, err := n.serviceClient.WarmDigests(ctx, req)
	if err != nil {
		sklog.Errorf("Error warming digests: %s", err)
	}
}

// UnavailableDigests, see the diff.DiffStore interface.
func (n *NetDiffStore) UnavailableDigests(ctx context.Context) map[types.Digest]*diff.DigestFailure {
	resp, err := n.serviceClient.UnavailableDigests(ctx, &Empty{})
	if err != nil {
		sklog.Debugf("Could not fetch unavailable digests: %s", err)
		return map[types.Digest]*diff.DigestFailure{}
	}

	ret := make(map[types.Digest]*diff.DigestFailure, len(resp.DigestFailures))
	for k, failure := range resp.DigestFailures {
		ret[types.Digest(k)] = &diff.DigestFailure{
			Digest: types.Digest(failure.Digest),
			Reason: diff.DiffErr(failure.Reason),
			TS:     failure.TS,
		}
	}
	return ret
}

// PurgeDigests, see the diff.DiffStore interface.
func (n *NetDiffStore) PurgeDigests(ctx context.Context, digests types.DigestSlice, purgeGCS bool) error {
	req := &PurgeDigestsRequest{Digests: common.AsStrings(digests), PurgeGCS: purgeGCS}
	_, err := n.serviceClient.PurgeDigests(ctx, req)
	if err != nil {
		return fmt.Errorf("Error purging digests: %s", err)
	}
	return nil
}
