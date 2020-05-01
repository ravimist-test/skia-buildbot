// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	tiling "go.skia.org/infra/golden/go/tiling"

	time "time"

	tracestore "go.skia.org/infra/golden/go/tracestore"
)

// TraceStore is an autogenerated mock type for the TraceStore type
type TraceStore struct {
	mock.Mock
}

// GetDenseTile provides a mock function with given fields: ctx, nCommits
func (_m *TraceStore) GetDenseTile(ctx context.Context, nCommits int) (*tiling.Tile, []*tiling.Commit, error) {
	ret := _m.Called(ctx, nCommits)

	var r0 *tiling.Tile
	if rf, ok := ret.Get(0).(func(context.Context, int) *tiling.Tile); ok {
		r0 = rf(ctx, nCommits)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*tiling.Tile)
		}
	}

	var r1 []*tiling.Commit
	if rf, ok := ret.Get(1).(func(context.Context, int) []*tiling.Commit); ok {
		r1 = rf(ctx, nCommits)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([]*tiling.Commit)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, int) error); ok {
		r2 = rf(ctx, nCommits)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetTile provides a mock function with given fields: ctx, nCommits
func (_m *TraceStore) GetTile(ctx context.Context, nCommits int) (*tiling.Tile, []*tiling.Commit, error) {
	ret := _m.Called(ctx, nCommits)

	var r0 *tiling.Tile
	if rf, ok := ret.Get(0).(func(context.Context, int) *tiling.Tile); ok {
		r0 = rf(ctx, nCommits)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*tiling.Tile)
		}
	}

	var r1 []*tiling.Commit
	if rf, ok := ret.Get(1).(func(context.Context, int) []*tiling.Commit); ok {
		r1 = rf(ctx, nCommits)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([]*tiling.Commit)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, int) error); ok {
		r2 = rf(ctx, nCommits)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// Put provides a mock function with given fields: ctx, commitHash, entries, ts
func (_m *TraceStore) Put(ctx context.Context, commitHash string, entries []*tracestore.Entry, ts time.Time) error {
	ret := _m.Called(ctx, commitHash, entries, ts)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, []*tracestore.Entry, time.Time) error); ok {
		r0 = rf(ctx, commitHash, entries, ts)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
