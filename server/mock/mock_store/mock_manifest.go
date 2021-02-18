// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/mattermost/mattermost-plugin-apps/server/store (interfaces: Manifest)

// Package mock_store is a generated GoMock package.
package mock_store

import (
	gomock "github.com/golang/mock/gomock"
	apps "github.com/mattermost/mattermost-plugin-apps/apps"
	config "github.com/mattermost/mattermost-plugin-apps/server/config"
	io "io"
	reflect "reflect"
)

// MockManifest is a mock of Manifest interface
type MockManifest struct {
	ctrl     *gomock.Controller
	recorder *MockManifestMockRecorder
}

// MockManifestMockRecorder is the mock recorder for MockManifest
type MockManifestMockRecorder struct {
	mock *MockManifest
}

// NewMockManifest creates a new mock instance
func NewMockManifest(ctrl *gomock.Controller) *MockManifest {
	mock := &MockManifest{ctrl: ctrl}
	mock.recorder = &MockManifestMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockManifest) EXPECT() *MockManifestMockRecorder {
	return m.recorder
}

// Configure mocks base method
func (m *MockManifest) Configure(arg0 config.Config) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Configure", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Configure indicates an expected call of Configure
func (mr *MockManifestMockRecorder) Configure(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Configure", reflect.TypeOf((*MockManifest)(nil).Configure), arg0)
}

// DeleteLocal mocks base method
func (m *MockManifest) DeleteLocal(arg0 apps.AppID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteLocal", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteLocal indicates an expected call of DeleteLocal
func (mr *MockManifestMockRecorder) DeleteLocal(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteLocal", reflect.TypeOf((*MockManifest)(nil).DeleteLocal), arg0)
}

// Get mocks base method
func (m *MockManifest) Get(arg0 apps.AppID) (*apps.Manifest, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0)
	ret0, _ := ret[0].(*apps.Manifest)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockManifestMockRecorder) Get(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockManifest)(nil).Get), arg0)
}

// Init mocks base method
func (m *MockManifest) Init(arg0 io.Reader, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Init", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Init indicates an expected call of Init
func (mr *MockManifestMockRecorder) Init(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Init", reflect.TypeOf((*MockManifest)(nil).Init), arg0, arg1)
}

// List mocks base method
func (m *MockManifest) List() map[apps.AppID]*apps.Manifest {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List")
	ret0, _ := ret[0].(map[apps.AppID]*apps.Manifest)
	return ret0
}

// List indicates an expected call of List
func (mr *MockManifestMockRecorder) List() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockManifest)(nil).List))
}

// StoreLocal mocks base method
func (m *MockManifest) StoreLocal(arg0 *apps.Manifest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StoreLocal", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// StoreLocal indicates an expected call of StoreLocal
func (mr *MockManifestMockRecorder) StoreLocal(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StoreLocal", reflect.TypeOf((*MockManifest)(nil).StoreLocal), arg0)
}
