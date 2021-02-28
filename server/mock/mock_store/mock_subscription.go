// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/mattermost/mattermost-plugin-apps/server/store (interfaces: Subscription)

// Package mock_store is a generated GoMock package.
package mock_store

import (
	gomock "github.com/golang/mock/gomock"
	apps "github.com/mattermost/mattermost-plugin-apps/apps"
	reflect "reflect"
)

// MockSubscription is a mock of Subscription interface
type MockSubscription struct {
	ctrl     *gomock.Controller
	recorder *MockSubscriptionMockRecorder
}

// MockSubscriptionMockRecorder is the mock recorder for MockSubscription
type MockSubscriptionMockRecorder struct {
	mock *MockSubscription
}

// NewMockSubscription creates a new mock instance
func NewMockSubscription(ctrl *gomock.Controller) *MockSubscription {
	mock := &MockSubscription{ctrl: ctrl}
	mock.recorder = &MockSubscriptionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSubscription) EXPECT() *MockSubscriptionMockRecorder {
	return m.recorder
}

// Delete mocks base method
func (m *MockSubscription) Delete(arg0 *apps.Subscription) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockSubscriptionMockRecorder) Delete(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockSubscription)(nil).Delete), arg0)
}

// Get mocks base method
func (m *MockSubscription) Get(arg0 apps.Subject, arg1, arg2 string) ([]*apps.Subscription, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1, arg2)
	ret0, _ := ret[0].([]*apps.Subscription)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockSubscriptionMockRecorder) Get(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockSubscription)(nil).Get), arg0, arg1, arg2)
}

// Save mocks base method
func (m *MockSubscription) Save(arg0 *apps.Subscription) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save
func (mr *MockSubscriptionMockRecorder) Save(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockSubscription)(nil).Save), arg0)
}
