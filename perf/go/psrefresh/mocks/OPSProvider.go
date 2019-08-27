// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	context "context"

	btts "go.skia.org/infra/perf/go/btts"

	mock "github.com/stretchr/testify/mock"

	paramtools "go.skia.org/infra/go/paramtools"
)

// OPSProvider is an autogenerated mock type for the OPSProvider type
type OPSProvider struct {
	mock.Mock
}

// GetLatestTile provides a mock function with given fields:
func (_m *OPSProvider) GetLatestTile() (btts.TileKey, error) {
	ret := _m.Called()

	var r0 btts.TileKey
	if rf, ok := ret.Get(0).(func() btts.TileKey); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(btts.TileKey)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetOrderedParamSet provides a mock function with given fields: ctx, tileKey
func (_m *OPSProvider) GetOrderedParamSet(ctx context.Context, tileKey btts.TileKey) (*paramtools.OrderedParamSet, error) {
	ret := _m.Called(ctx, tileKey)

	var r0 *paramtools.OrderedParamSet
	if rf, ok := ret.Get(0).(func(context.Context, btts.TileKey) *paramtools.OrderedParamSet); ok {
		r0 = rf(ctx, tileKey)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*paramtools.OrderedParamSet)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, btts.TileKey) error); ok {
		r1 = rf(ctx, tileKey)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}