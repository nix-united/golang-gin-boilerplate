// Code generated by MockGen. DO NOT EDIT.
// Source: post_handler.go
//
// Generated by this command:
//
//	mockgen -source=post_handler.go -destination=post_handler_mock_test.go -package=handler_test -typed=true
//

// Package handler_test is a generated GoMock package.
package handler_test

import (
	reflect "reflect"

	model "github.com/nix-united/golang-gin-boilerplate/internal/model"
	gomock "go.uber.org/mock/gomock"
)

// MockpostService is a mock of postService interface.
type MockpostService struct {
	ctrl     *gomock.Controller
	recorder *MockpostServiceMockRecorder
}

// MockpostServiceMockRecorder is the mock recorder for MockpostService.
type MockpostServiceMockRecorder struct {
	mock *MockpostService
}

// NewMockpostService creates a new mock instance.
func NewMockpostService(ctrl *gomock.Controller) *MockpostService {
	mock := &MockpostService{ctrl: ctrl}
	mock.recorder = &MockpostServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockpostService) EXPECT() *MockpostServiceMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockpostService) Create(post *model.Post) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", post)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockpostServiceMockRecorder) Create(post any) *MockpostServiceCreateCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockpostService)(nil).Create), post)
	return &MockpostServiceCreateCall{Call: call}
}

// MockpostServiceCreateCall wrap *gomock.Call
type MockpostServiceCreateCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockpostServiceCreateCall) Return(arg0 error) *MockpostServiceCreateCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockpostServiceCreateCall) Do(f func(*model.Post) error) *MockpostServiceCreateCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockpostServiceCreateCall) DoAndReturn(f func(*model.Post) error) *MockpostServiceCreateCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// CreatePost mocks base method.
func (m *MockpostService) CreatePost(title, content string, userID uint) (*model.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePost", title, content, userID)
	ret0, _ := ret[0].(*model.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreatePost indicates an expected call of CreatePost.
func (mr *MockpostServiceMockRecorder) CreatePost(title, content, userID any) *MockpostServiceCreatePostCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePost", reflect.TypeOf((*MockpostService)(nil).CreatePost), title, content, userID)
	return &MockpostServiceCreatePostCall{Call: call}
}

// MockpostServiceCreatePostCall wrap *gomock.Call
type MockpostServiceCreatePostCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockpostServiceCreatePostCall) Return(arg0 *model.Post, arg1 error) *MockpostServiceCreatePostCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockpostServiceCreatePostCall) Do(f func(string, string, uint) (*model.Post, error)) *MockpostServiceCreatePostCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockpostServiceCreatePostCall) DoAndReturn(f func(string, string, uint) (*model.Post, error)) *MockpostServiceCreatePostCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Delete mocks base method.
func (m *MockpostService) Delete(post *model.Post) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", post)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockpostServiceMockRecorder) Delete(post any) *MockpostServiceDeleteCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockpostService)(nil).Delete), post)
	return &MockpostServiceDeleteCall{Call: call}
}

// MockpostServiceDeleteCall wrap *gomock.Call
type MockpostServiceDeleteCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockpostServiceDeleteCall) Return(arg0 error) *MockpostServiceDeleteCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockpostServiceDeleteCall) Do(f func(*model.Post) error) *MockpostServiceDeleteCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockpostServiceDeleteCall) DoAndReturn(f func(*model.Post) error) *MockpostServiceDeleteCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetAll mocks base method.
func (m *MockpostService) GetAll() ([]model.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll")
	ret0, _ := ret[0].([]model.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockpostServiceMockRecorder) GetAll() *MockpostServiceGetAllCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockpostService)(nil).GetAll))
	return &MockpostServiceGetAllCall{Call: call}
}

// MockpostServiceGetAllCall wrap *gomock.Call
type MockpostServiceGetAllCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockpostServiceGetAllCall) Return(arg0 []model.Post, arg1 error) *MockpostServiceGetAllCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockpostServiceGetAllCall) Do(f func() ([]model.Post, error)) *MockpostServiceGetAllCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockpostServiceGetAllCall) DoAndReturn(f func() ([]model.Post, error)) *MockpostServiceGetAllCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetByID mocks base method.
func (m *MockpostService) GetByID(id int) (*model.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", id)
	ret0, _ := ret[0].(*model.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockpostServiceMockRecorder) GetByID(id any) *MockpostServiceGetByIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockpostService)(nil).GetByID), id)
	return &MockpostServiceGetByIDCall{Call: call}
}

// MockpostServiceGetByIDCall wrap *gomock.Call
type MockpostServiceGetByIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockpostServiceGetByIDCall) Return(arg0 *model.Post, arg1 error) *MockpostServiceGetByIDCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockpostServiceGetByIDCall) Do(f func(int) (*model.Post, error)) *MockpostServiceGetByIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockpostServiceGetByIDCall) DoAndReturn(f func(int) (*model.Post, error)) *MockpostServiceGetByIDCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Save mocks base method.
func (m *MockpostService) Save(post *model.Post) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", post)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockpostServiceMockRecorder) Save(post any) *MockpostServiceSaveCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockpostService)(nil).Save), post)
	return &MockpostServiceSaveCall{Call: call}
}

// MockpostServiceSaveCall wrap *gomock.Call
type MockpostServiceSaveCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockpostServiceSaveCall) Return(arg0 error) *MockpostServiceSaveCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockpostServiceSaveCall) Do(f func(*model.Post) error) *MockpostServiceSaveCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockpostServiceSaveCall) DoAndReturn(f func(*model.Post) error) *MockpostServiceSaveCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
