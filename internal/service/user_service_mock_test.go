// Code generated by MockGen. DO NOT EDIT.
// Source: user_service.go
//
// Generated by this command:
//
//	mockgen -source=user_service.go -destination=user_service_mock_test.go -package=service_test -typed=true
//

// Package service_test is a generated GoMock package.
package service_test

import (
	reflect "reflect"

	model "github.com/nix-united/golang-gin-boilerplate/internal/model"
	gomock "go.uber.org/mock/gomock"
)

// MockuserRepository is a mock of userRepository interface.
type MockuserRepository struct {
	ctrl     *gomock.Controller
	recorder *MockuserRepositoryMockRecorder
}

// MockuserRepositoryMockRecorder is the mock recorder for MockuserRepository.
type MockuserRepositoryMockRecorder struct {
	mock *MockuserRepository
}

// NewMockuserRepository creates a new mock instance.
func NewMockuserRepository(ctrl *gomock.Controller) *MockuserRepository {
	mock := &MockuserRepository{ctrl: ctrl}
	mock.recorder = &MockuserRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockuserRepository) EXPECT() *MockuserRepositoryMockRecorder {
	return m.recorder
}

// FindUserByEmail mocks base method.
func (m *MockuserRepository) FindUserByEmail(email string) (model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindUserByEmail", email)
	ret0, _ := ret[0].(model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindUserByEmail indicates an expected call of FindUserByEmail.
func (mr *MockuserRepositoryMockRecorder) FindUserByEmail(email any) *MockuserRepositoryFindUserByEmailCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindUserByEmail", reflect.TypeOf((*MockuserRepository)(nil).FindUserByEmail), email)
	return &MockuserRepositoryFindUserByEmailCall{Call: call}
}

// MockuserRepositoryFindUserByEmailCall wrap *gomock.Call
type MockuserRepositoryFindUserByEmailCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockuserRepositoryFindUserByEmailCall) Return(arg0 model.User, arg1 error) *MockuserRepositoryFindUserByEmailCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockuserRepositoryFindUserByEmailCall) Do(f func(string) (model.User, error)) *MockuserRepositoryFindUserByEmailCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockuserRepositoryFindUserByEmailCall) DoAndReturn(f func(string) (model.User, error)) *MockuserRepositoryFindUserByEmailCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// StoreUser mocks base method.
func (m *MockuserRepository) StoreUser(user model.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StoreUser", user)
	ret0, _ := ret[0].(error)
	return ret0
}

// StoreUser indicates an expected call of StoreUser.
func (mr *MockuserRepositoryMockRecorder) StoreUser(user any) *MockuserRepositoryStoreUserCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StoreUser", reflect.TypeOf((*MockuserRepository)(nil).StoreUser), user)
	return &MockuserRepositoryStoreUserCall{Call: call}
}

// MockuserRepositoryStoreUserCall wrap *gomock.Call
type MockuserRepositoryStoreUserCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockuserRepositoryStoreUserCall) Return(arg0 error) *MockuserRepositoryStoreUserCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockuserRepositoryStoreUserCall) Do(f func(model.User) error) *MockuserRepositoryStoreUserCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockuserRepositoryStoreUserCall) DoAndReturn(f func(model.User) error) *MockuserRepositoryStoreUserCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Mockencryptor is a mock of encryptor interface.
type Mockencryptor struct {
	ctrl     *gomock.Controller
	recorder *MockencryptorMockRecorder
}

// MockencryptorMockRecorder is the mock recorder for Mockencryptor.
type MockencryptorMockRecorder struct {
	mock *Mockencryptor
}

// NewMockencryptor creates a new mock instance.
func NewMockencryptor(ctrl *gomock.Controller) *Mockencryptor {
	mock := &Mockencryptor{ctrl: ctrl}
	mock.recorder = &MockencryptorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockencryptor) EXPECT() *MockencryptorMockRecorder {
	return m.recorder
}

// Encrypt mocks base method.
func (m *Mockencryptor) Encrypt(str string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Encrypt", str)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Encrypt indicates an expected call of Encrypt.
func (mr *MockencryptorMockRecorder) Encrypt(str any) *MockencryptorEncryptCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Encrypt", reflect.TypeOf((*Mockencryptor)(nil).Encrypt), str)
	return &MockencryptorEncryptCall{Call: call}
}

// MockencryptorEncryptCall wrap *gomock.Call
type MockencryptorEncryptCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockencryptorEncryptCall) Return(arg0 string, arg1 error) *MockencryptorEncryptCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockencryptorEncryptCall) Do(f func(string) (string, error)) *MockencryptorEncryptCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockencryptorEncryptCall) DoAndReturn(f func(string) (string, error)) *MockencryptorEncryptCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
