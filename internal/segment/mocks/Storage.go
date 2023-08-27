// Code generated by mockery v2.33.0. DO NOT EDIT.

package mocks

import (
	context "context"

	segment "github.com/POMBNK/avito_test_task/internal/segment"
	mock "github.com/stretchr/testify/mock"
)

// Storage is an autogenerated mock type for the Storage type
type Storage struct {
	mock.Mock
}

// AddUserToSegments provides a mock function with given fields: ctx, segmentsUser, segmentName
func (_m *Storage) AddUserToSegments(ctx context.Context, segmentsUser segment.SegmentsUsers, segmentName string) error {
	ret := _m.Called(ctx, segmentsUser, segmentName)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, segment.SegmentsUsers, string) error); ok {
		r0 = rf(ctx, segmentsUser, segmentName)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Create provides a mock function with given fields: ctx, _a1
func (_m *Storage) Create(ctx context.Context, _a1 segment.Segment) (string, error) {
	ret := _m.Called(ctx, _a1)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, segment.Segment) (string, error)); ok {
		return rf(ctx, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, segment.Segment) string); ok {
		r0 = rf(ctx, _a1)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, segment.Segment) error); ok {
		r1 = rf(ctx, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: ctx, _a1
func (_m *Storage) Delete(ctx context.Context, _a1 segment.Segment) error {
	ret := _m.Called(ctx, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, segment.Segment) error); ok {
		r0 = rf(ctx, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteSegmentFromUser provides a mock function with given fields: ctx, segmentsUser, segmentName
func (_m *Storage) DeleteSegmentFromUser(ctx context.Context, segmentsUser segment.SegmentsUsers, segmentName string) error {
	ret := _m.Called(ctx, segmentsUser, segmentName)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, segment.SegmentsUsers, string) error); ok {
		r0 = rf(ctx, segmentsUser, segmentName)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetActiveSegments provides a mock function with given fields: ctx, userID
func (_m *Storage) GetActiveSegments(ctx context.Context, userID string) ([]segment.ActiveSegments, error) {
	ret := _m.Called(ctx, userID)

	var r0 []segment.ActiveSegments
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]segment.ActiveSegments, error)); ok {
		return rf(ctx, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []segment.ActiveSegments); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]segment.ActiveSegments)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserHistoryOptimized provides a mock function with given fields: ctx, userID, timestampz
func (_m *Storage) GetUserHistoryOptimized(ctx context.Context, userID string, timestampz string) ([]segment.BetterCSVReport, error) {
	ret := _m.Called(ctx, userID, timestampz)

	var r0 []segment.BetterCSVReport
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) ([]segment.BetterCSVReport, error)); ok {
		return rf(ctx, userID, timestampz)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) []segment.BetterCSVReport); ok {
		r0 = rf(ctx, userID, timestampz)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]segment.BetterCSVReport)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, userID, timestampz)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserHistoryOriginal provides a mock function with given fields: ctx, userID, timestampz
func (_m *Storage) GetUserHistoryOriginal(ctx context.Context, userID string, timestampz string) ([]segment.CSVReport, error) {
	ret := _m.Called(ctx, userID, timestampz)

	var r0 []segment.CSVReport
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) ([]segment.CSVReport, error)); ok {
		return rf(ctx, userID, timestampz)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) []segment.CSVReport); ok {
		r0 = rf(ctx, userID, timestampz)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]segment.CSVReport)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, userID, timestampz)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsUserExist provides a mock function with given fields: ctx, segmentsUser
func (_m *Storage) IsUserExist(ctx context.Context, segmentsUser segment.SegmentsUsers) error {
	ret := _m.Called(ctx, segmentsUser)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, segment.SegmentsUsers) error); ok {
		r0 = rf(ctx, segmentsUser)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewStorage creates a new instance of Storage. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewStorage(t interface {
	mock.TestingT
	Cleanup(func())
}) *Storage {
	mock := &Storage{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
