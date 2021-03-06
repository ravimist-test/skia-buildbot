// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	incident "go.skia.org/infra/am/go/incident"

	silence "go.skia.org/infra/am/go/silence"
)

// APIClient is an autogenerated mock type for the APIClient type
type APIClient struct {
	mock.Mock
}

// GetAlerts provides a mock function with given fields:
func (_m *APIClient) GetAlerts() ([]incident.Incident, error) {
	ret := _m.Called()

	var r0 []incident.Incident
	if rf, ok := ret.Get(0).(func() []incident.Incident); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]incident.Incident)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSilences provides a mock function with given fields:
func (_m *APIClient) GetSilences() ([]silence.Silence, error) {
	ret := _m.Called()

	var r0 []silence.Silence
	if rf, ok := ret.Get(0).(func() []silence.Silence); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]silence.Silence)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
