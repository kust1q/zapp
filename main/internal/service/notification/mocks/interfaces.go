
// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	entity "github.com/kust1q/Zapp/main/internal/domain/entity"
	gomock "go.uber.org/mock/gomock"
)

// MockhubProvider is a mock of hubProvider interface.
type MockhubProvider struct {
	ctrl     *gomock.Controller
	recorder *MockhubProviderMockRecorder
	isgomock struct{}
}

// MockhubProviderMockRecorder is the mock recorder for MockhubProvider.
type MockhubProviderMockRecorder struct {
	mock *MockhubProvider
}

// NewMockhubProvider creates a new mock instance.
func NewMockhubProvider(ctrl *gomock.Controller) *MockhubProvider {
	mock := &MockhubProvider{ctrl: ctrl}
	mock.recorder = &MockhubProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockhubProvider) EXPECT() *MockhubProviderMockRecorder {
	return m.recorder
}

// SendNotification mocks base method.
func (m *MockhubProvider) SendNotification(notification *entity.Notification) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SendNotification", notification)
}

// SendNotification indicates an expected call of SendNotification.
func (mr *MockhubProviderMockRecorder) SendNotification(notification any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendNotification", reflect.TypeOf((*MockhubProvider)(nil).SendNotification), notification)
}

// Mockdb is a mock of db interface.
type Mockdb struct {
	ctrl     *gomock.Controller
	recorder *MockdbMockRecorder
	isgomock struct{}
}

// MockdbMockRecorder is the mock recorder for Mockdb.
type MockdbMockRecorder struct {
	mock *Mockdb
}

// NewMockdb creates a new mock instance.
func NewMockdb(ctrl *gomock.Controller) *Mockdb {
	mock := &Mockdb{ctrl: ctrl}
	mock.recorder = &MockdbMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockdb) EXPECT() *MockdbMockRecorder {
	return m.recorder
}

// GetTweetById mocks base method.
func (m *Mockdb) GetTweetById(ctx context.Context, tweetID int) (*entity.Tweet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTweetById", ctx, tweetID)
	ret0, _ := ret[0].(*entity.Tweet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTweetById indicates an expected call of GetTweetById.
func (mr *MockdbMockRecorder) GetTweetById(ctx, tweetID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTweetById", reflect.TypeOf((*Mockdb)(nil).GetTweetById), ctx, tweetID)
}

// GetUserByID mocks base method.
func (m *Mockdb) GetUserByID(ctx context.Context, userID int) (*entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByID", ctx, userID)
	ret0, _ := ret[0].(*entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByID indicates an expected call of GetUserByID.
func (mr *MockdbMockRecorder) GetUserByID(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByID", reflect.TypeOf((*Mockdb)(nil).GetUserByID), ctx, userID)
}
