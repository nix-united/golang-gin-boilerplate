// Code generated by MockGen. DO NOT EDIT.
// Source: auth_handler.go
//
// Generated by this command:
//
//	mockgen -source=auth_handler.go -destination=auth_handler_mock_test.go -package=handler_test -typed=true
//

// Package handler_test is a generated GoMock package.
package handler_test

import (
	reflect "reflect"

	request "github.com/nix-united/golang-gin-boilerplate/internal/request"
	gomock "go.uber.org/mock/gomock"
)

// MockuserService is a mock of userService interface.
type MockuserService struct {
	ctrl     *gomock.Controller
	recorder *MockuserServiceMockRecorder
}

// MockuserServiceMockRecorder is the mock recorder for MockuserService.
type MockuserServiceMockRecorder struct {
	mock *MockuserService
}

// NewMockuserService creates a new mock instance.
func NewMockuserService(ctrl *gomock.Controller) *MockuserService {
	mock := &MockuserService{ctrl: ctrl}
	mock.recorder = &MockuserServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockuserService) EXPECT() *MockuserServiceMockRecorder {
	return m.recorder
}

// CreateUser mocks base method.
func (m *MockuserService) CreateUser(req request.RegisterRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", req)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockuserServiceMockRecorder) CreateUser(req any) *MockuserServiceCreateUserCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockuserService)(nil).CreateUser), req)
	return &MockuserServiceCreateUserCall{Call: call}
}

// MockuserServiceCreateUserCall wrap *gomock.Call
type MockuserServiceCreateUserCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockuserServiceCreateUserCall) Return(arg0 error) *MockuserServiceCreateUserCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockuserServiceCreateUserCall) Do(f func(request.RegisterRequest) error) *MockuserServiceCreateUserCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockuserServiceCreateUserCall) DoAndReturn(f func(request.RegisterRequest) error) *MockuserServiceCreateUserCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
