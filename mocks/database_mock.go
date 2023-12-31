// Code generated by MockGen. DO NOT EDIT.
// Source: internal/db/db.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"
	time "time"
	models "user-segmentation-service/internal/models"

	gomock "github.com/golang/mock/gomock"
)

// MockInterface is a mock of Interface interface.
type MockInterface struct {
	ctrl     *gomock.Controller
	recorder *MockInterfaceMockRecorder
}

// MockInterfaceMockRecorder is the mock recorder for MockInterface.
type MockInterfaceMockRecorder struct {
	mock *MockInterface
}

// NewMockInterface creates a new mock instance.
func NewMockInterface(ctrl *gomock.Controller) *MockInterface {
	mock := &MockInterface{ctrl: ctrl}
	mock.recorder = &MockInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockInterface) EXPECT() *MockInterfaceMockRecorder {
	return m.recorder
}

// CreateSegment mocks base method.
func (m *MockInterface) CreateSegment(slug string, randomPercentage float64, expirationDate time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSegment", slug, randomPercentage, expirationDate)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateSegment indicates an expected call of CreateSegment.
func (mr *MockInterfaceMockRecorder) CreateSegment(slug, randomPercentage, expirationDate interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSegment", reflect.TypeOf((*MockInterface)(nil).CreateSegment), slug, randomPercentage, expirationDate)
}

// CreateUser mocks base method.
func (m *MockInterface) CreateUser(name string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", name)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockInterfaceMockRecorder) CreateUser(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockInterface)(nil).CreateUser), name)
}

// DeleteSegment mocks base method.
func (m *MockInterface) DeleteSegment(slug string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSegment", slug)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteSegment indicates an expected call of DeleteSegment.
func (mr *MockInterfaceMockRecorder) DeleteSegment(slug interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSegment", reflect.TypeOf((*MockInterface)(nil).DeleteSegment), slug)
}

// DeleteUser mocks base method.
func (m *MockInterface) DeleteUser(userID int) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUser", userID)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteUser indicates an expected call of DeleteUser.
func (mr *MockInterfaceMockRecorder) DeleteUser(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUser", reflect.TypeOf((*MockInterface)(nil).DeleteUser), userID)
}

// GetUserReport mocks base method.
func (m *MockInterface) GetUserReport(userID int, yearMonth string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserReport", userID, yearMonth)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserReport indicates an expected call of GetUserReport.
func (mr *MockInterfaceMockRecorder) GetUserReport(userID, yearMonth interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserReport", reflect.TypeOf((*MockInterface)(nil).GetUserReport), userID, yearMonth)
}

// GetUserSegments mocks base method.
func (m *MockInterface) GetUserSegments(userID int) (int, []string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserSegments", userID)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].([]string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetUserSegments indicates an expected call of GetUserSegments.
func (mr *MockInterfaceMockRecorder) GetUserSegments(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserSegments", reflect.TypeOf((*MockInterface)(nil).GetUserSegments), userID)
}

// UpdateUserSegments mocks base method.
func (m *MockInterface) UpdateUserSegments(userID int, addList []models.Segment, removeList []string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserSegments", userID, addList, removeList)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateUserSegments indicates an expected call of UpdateUserSegments.
func (mr *MockInterfaceMockRecorder) UpdateUserSegments(userID, addList, removeList interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserSegments", reflect.TypeOf((*MockInterface)(nil).UpdateUserSegments), userID, addList, removeList)
}
