// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/cloudquery/cloudquery/plugins/source/aws/client (interfaces: KinesisClient)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	kinesis "github.com/aws/aws-sdk-go-v2/service/kinesis"
	gomock "github.com/golang/mock/gomock"
)

// MockKinesisClient is a mock of KinesisClient interface.
type MockKinesisClient struct {
	ctrl     *gomock.Controller
	recorder *MockKinesisClientMockRecorder
}

// MockKinesisClientMockRecorder is the mock recorder for MockKinesisClient.
type MockKinesisClientMockRecorder struct {
	mock *MockKinesisClient
}

// NewMockKinesisClient creates a new mock instance.
func NewMockKinesisClient(ctrl *gomock.Controller) *MockKinesisClient {
	mock := &MockKinesisClient{ctrl: ctrl}
	mock.recorder = &MockKinesisClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockKinesisClient) EXPECT() *MockKinesisClientMockRecorder {
	return m.recorder
}

// DescribeStreamSummary mocks base method.
func (m *MockKinesisClient) DescribeStreamSummary(arg0 context.Context, arg1 *kinesis.DescribeStreamSummaryInput, arg2 ...func(*kinesis.Options)) (*kinesis.DescribeStreamSummaryOutput, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DescribeStreamSummary", varargs...)
	ret0, _ := ret[0].(*kinesis.DescribeStreamSummaryOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DescribeStreamSummary indicates an expected call of DescribeStreamSummary.
func (mr *MockKinesisClientMockRecorder) DescribeStreamSummary(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DescribeStreamSummary", reflect.TypeOf((*MockKinesisClient)(nil).DescribeStreamSummary), varargs...)
}

// ListStreams mocks base method.
func (m *MockKinesisClient) ListStreams(arg0 context.Context, arg1 *kinesis.ListStreamsInput, arg2 ...func(*kinesis.Options)) (*kinesis.ListStreamsOutput, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListStreams", varargs...)
	ret0, _ := ret[0].(*kinesis.ListStreamsOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListStreams indicates an expected call of ListStreams.
func (mr *MockKinesisClientMockRecorder) ListStreams(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListStreams", reflect.TypeOf((*MockKinesisClient)(nil).ListStreams), varargs...)
}

// ListTagsForStream mocks base method.
func (m *MockKinesisClient) ListTagsForStream(arg0 context.Context, arg1 *kinesis.ListTagsForStreamInput, arg2 ...func(*kinesis.Options)) (*kinesis.ListTagsForStreamOutput, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListTagsForStream", varargs...)
	ret0, _ := ret[0].(*kinesis.ListTagsForStreamOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListTagsForStream indicates an expected call of ListTagsForStream.
func (mr *MockKinesisClientMockRecorder) ListTagsForStream(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListTagsForStream", reflect.TypeOf((*MockKinesisClient)(nil).ListTagsForStream), varargs...)
}
