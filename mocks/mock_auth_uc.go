// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/gosom/bookclub (interfaces: AuthUseCases,JWTProvider)
//
// Generated by this command:
//
//	mockgen -package mocks -destination ./mocks/mock_auth_uc.go . AuthUseCases,JWTProvider
//
// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	bookclub "github.com/gosom/bookclub"
	gomock "go.uber.org/mock/gomock"
)

// MockAuthUseCases is a mock of AuthUseCases interface.
type MockAuthUseCases struct {
	ctrl     *gomock.Controller
	recorder *MockAuthUseCasesMockRecorder
}

// MockAuthUseCasesMockRecorder is the mock recorder for MockAuthUseCases.
type MockAuthUseCasesMockRecorder struct {
	mock *MockAuthUseCases
}

// NewMockAuthUseCases creates a new mock instance.
func NewMockAuthUseCases(ctrl *gomock.Controller) *MockAuthUseCases {
	mock := &MockAuthUseCases{ctrl: ctrl}
	mock.recorder = &MockAuthUseCasesMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthUseCases) EXPECT() *MockAuthUseCasesMockRecorder {
	return m.recorder
}

// Login mocks base method.
func (m *MockAuthUseCases) Login(arg0 context.Context, arg1 bookclub.LoginParams) (string, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Login indicates an expected call of Login.
func (mr *MockAuthUseCasesMockRecorder) Login(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockAuthUseCases)(nil).Login), arg0, arg1)
}

// Refresh mocks base method.
func (m *MockAuthUseCases) Refresh(arg0 context.Context, arg1 string) (string, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Refresh", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Refresh indicates an expected call of Refresh.
func (mr *MockAuthUseCasesMockRecorder) Refresh(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Refresh", reflect.TypeOf((*MockAuthUseCases)(nil).Refresh), arg0, arg1)
}

// MockJWTProvider is a mock of JWTProvider interface.
type MockJWTProvider struct {
	ctrl     *gomock.Controller
	recorder *MockJWTProviderMockRecorder
}

// MockJWTProviderMockRecorder is the mock recorder for MockJWTProvider.
type MockJWTProviderMockRecorder struct {
	mock *MockJWTProvider
}

// NewMockJWTProvider creates a new mock instance.
func NewMockJWTProvider(ctrl *gomock.Controller) *MockJWTProvider {
	mock := &MockJWTProvider{ctrl: ctrl}
	mock.recorder = &MockJWTProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockJWTProvider) EXPECT() *MockJWTProviderMockRecorder {
	return m.recorder
}

// GenerateRefreshToken mocks base method.
func (m *MockJWTProvider) GenerateRefreshToken(arg0 context.Context, arg1 bookclub.JWTGenerateParams) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateRefreshToken", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateRefreshToken indicates an expected call of GenerateRefreshToken.
func (mr *MockJWTProviderMockRecorder) GenerateRefreshToken(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateRefreshToken", reflect.TypeOf((*MockJWTProvider)(nil).GenerateRefreshToken), arg0, arg1)
}

// GenerateToken mocks base method.
func (m *MockJWTProvider) GenerateToken(arg0 context.Context, arg1 bookclub.JWTGenerateParams) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateToken", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateToken indicates an expected call of GenerateToken.
func (mr *MockJWTProviderMockRecorder) GenerateToken(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateToken", reflect.TypeOf((*MockJWTProvider)(nil).GenerateToken), arg0, arg1)
}

// ValidateToken mocks base method.
func (m *MockJWTProvider) ValidateToken(arg0 context.Context, arg1 string) (bookclub.JWTClaims, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateToken", arg0, arg1)
	ret0, _ := ret[0].(bookclub.JWTClaims)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateToken indicates an expected call of ValidateToken.
func (mr *MockJWTProviderMockRecorder) ValidateToken(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateToken", reflect.TypeOf((*MockJWTProvider)(nil).ValidateToken), arg0, arg1)
}
