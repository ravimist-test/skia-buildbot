// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	dataframe "go.skia.org/infra/perf/go/dataframe"

	paramtools "go.skia.org/infra/go/paramtools"

	query "go.skia.org/infra/go/query"

	time "time"

	types "go.skia.org/infra/perf/go/types"
)

// DataFrameBuilder is an autogenerated mock type for the DataFrameBuilder type
type DataFrameBuilder struct {
	mock.Mock
}

// NewFromKeysAndRange provides a mock function with given fields: ctx, keys, begin, end, downsample, progress
func (_m *DataFrameBuilder) NewFromKeysAndRange(ctx context.Context, keys []string, begin time.Time, end time.Time, downsample bool, progress types.Progress) (*dataframe.DataFrame, error) {
	ret := _m.Called(ctx, keys, begin, end, downsample, progress)

	var r0 *dataframe.DataFrame
	if rf, ok := ret.Get(0).(func(context.Context, []string, time.Time, time.Time, bool, types.Progress) *dataframe.DataFrame); ok {
		r0 = rf(ctx, keys, begin, end, downsample, progress)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dataframe.DataFrame)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []string, time.Time, time.Time, bool, types.Progress) error); ok {
		r1 = rf(ctx, keys, begin, end, downsample, progress)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewFromQueryAndRange provides a mock function with given fields: ctx, begin, end, q, downsample, progress
func (_m *DataFrameBuilder) NewFromQueryAndRange(ctx context.Context, begin time.Time, end time.Time, q *query.Query, downsample bool, progress types.Progress) (*dataframe.DataFrame, error) {
	ret := _m.Called(ctx, begin, end, q, downsample, progress)

	var r0 *dataframe.DataFrame
	if rf, ok := ret.Get(0).(func(context.Context, time.Time, time.Time, *query.Query, bool, types.Progress) *dataframe.DataFrame); ok {
		r0 = rf(ctx, begin, end, q, downsample, progress)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dataframe.DataFrame)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, time.Time, time.Time, *query.Query, bool, types.Progress) error); ok {
		r1 = rf(ctx, begin, end, q, downsample, progress)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewNFromKeys provides a mock function with given fields: ctx, end, keys, n, progress
func (_m *DataFrameBuilder) NewNFromKeys(ctx context.Context, end time.Time, keys []string, n int32, progress types.Progress) (*dataframe.DataFrame, error) {
	ret := _m.Called(ctx, end, keys, n, progress)

	var r0 *dataframe.DataFrame
	if rf, ok := ret.Get(0).(func(context.Context, time.Time, []string, int32, types.Progress) *dataframe.DataFrame); ok {
		r0 = rf(ctx, end, keys, n, progress)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dataframe.DataFrame)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, time.Time, []string, int32, types.Progress) error); ok {
		r1 = rf(ctx, end, keys, n, progress)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewNFromQuery provides a mock function with given fields: ctx, end, q, n, progress
func (_m *DataFrameBuilder) NewNFromQuery(ctx context.Context, end time.Time, q *query.Query, n int32, progress types.Progress) (*dataframe.DataFrame, error) {
	ret := _m.Called(ctx, end, q, n, progress)

	var r0 *dataframe.DataFrame
	if rf, ok := ret.Get(0).(func(context.Context, time.Time, *query.Query, int32, types.Progress) *dataframe.DataFrame); ok {
		r0 = rf(ctx, end, q, n, progress)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dataframe.DataFrame)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, time.Time, *query.Query, int32, types.Progress) error); ok {
		r1 = rf(ctx, end, q, n, progress)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PreflightQuery provides a mock function with given fields: ctx, end, q
func (_m *DataFrameBuilder) PreflightQuery(ctx context.Context, end time.Time, q *query.Query) (int64, paramtools.ParamSet, error) {
	ret := _m.Called(ctx, end, q)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, time.Time, *query.Query) int64); ok {
		r0 = rf(ctx, end, q)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 paramtools.ParamSet
	if rf, ok := ret.Get(1).(func(context.Context, time.Time, *query.Query) paramtools.ParamSet); ok {
		r1 = rf(ctx, end, q)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(paramtools.ParamSet)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, time.Time, *query.Query) error); ok {
		r2 = rf(ctx, end, q)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}
