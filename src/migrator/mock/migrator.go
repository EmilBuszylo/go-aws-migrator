// Code generated by MockGen. DO NOT EDIT.
// Source: ../migrator.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	dynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	gomock "github.com/golang/mock/gomock"
)

// MockAPI is a mock of API interface.
type MockAPI struct {
	ctrl     *gomock.Controller
	recorder *MockAPIMockRecorder
}

// MockAPIMockRecorder is the mock recorder for MockAPI.
type MockAPIMockRecorder struct {
	mock *MockAPI
}

// NewMockAPI creates a new mock instance.
func NewMockAPI(ctrl *gomock.Controller) *MockAPI {
	mock := &MockAPI{ctrl: ctrl}
	mock.recorder = &MockAPIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAPI) EXPECT() *MockAPIMockRecorder {
	return m.recorder
}

// PutItem mocks base method.
func (m *MockAPI) PutItem(arg0 context.Context, arg1 *dynamodb.PutItemInput, arg2 ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "PutItem", varargs...)
	ret0, _ := ret[0].(*dynamodb.PutItemOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PutItem indicates an expected call of PutItem.
func (mr *MockAPIMockRecorder) PutItem(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PutItem", reflect.TypeOf((*MockAPI)(nil).PutItem), varargs...)
}

// Query mocks base method.
func (m *MockAPI) Query(arg0 context.Context, arg1 *dynamodb.QueryInput, arg2 ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Query", varargs...)
	ret0, _ := ret[0].(*dynamodb.QueryOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Query indicates an expected call of Query.
func (mr *MockAPIMockRecorder) Query(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Query", reflect.TypeOf((*MockAPI)(nil).Query), varargs...)
}
