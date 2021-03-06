// Automatically generated by MockGen. DO NOT EDIT!
// Source: github.com/topfreegames/pitaya/cluster (interfaces: ServiceDiscovery)

package mocks

import (
	gomock "github.com/golang/mock/gomock"
	cluster "github.com/topfreegames/pitaya/cluster"
)

// Mock of ServiceDiscovery interface
type MockServiceDiscovery struct {
	ctrl     *gomock.Controller
	recorder *_MockServiceDiscoveryRecorder
}

// Recorder for MockServiceDiscovery (not exported)
type _MockServiceDiscoveryRecorder struct {
	mock *MockServiceDiscovery
}

func NewMockServiceDiscovery(ctrl *gomock.Controller) *MockServiceDiscovery {
	mock := &MockServiceDiscovery{ctrl: ctrl}
	mock.recorder = &_MockServiceDiscoveryRecorder{mock}
	return mock
}

func (_m *MockServiceDiscovery) EXPECT() *_MockServiceDiscoveryRecorder {
	return _m.recorder
}

func (_m *MockServiceDiscovery) AfterInit() {
	_m.ctrl.Call(_m, "AfterInit")
}

func (_mr *_MockServiceDiscoveryRecorder) AfterInit() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "AfterInit")
}

func (_m *MockServiceDiscovery) BeforeShutdown() {
	_m.ctrl.Call(_m, "BeforeShutdown")
}

func (_mr *_MockServiceDiscoveryRecorder) BeforeShutdown() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "BeforeShutdown")
}

func (_m *MockServiceDiscovery) GetServer(_param0 string) (*cluster.Server, error) {
	ret := _m.ctrl.Call(_m, "GetServer", _param0)
	ret0, _ := ret[0].(*cluster.Server)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockServiceDiscoveryRecorder) GetServer(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetServer", arg0)
}

func (_m *MockServiceDiscovery) GetServersByType(_param0 string) (map[string]*cluster.Server, error) {
	ret := _m.ctrl.Call(_m, "GetServersByType", _param0)
	ret0, _ := ret[0].(map[string]*cluster.Server)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockServiceDiscoveryRecorder) GetServersByType(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetServersByType", arg0)
}

func (_m *MockServiceDiscovery) Init() error {
	ret := _m.ctrl.Call(_m, "Init")
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockServiceDiscoveryRecorder) Init() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Init")
}

func (_m *MockServiceDiscovery) Shutdown() error {
	ret := _m.ctrl.Call(_m, "Shutdown")
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockServiceDiscoveryRecorder) Shutdown() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Shutdown")
}

func (_m *MockServiceDiscovery) SyncServers() error {
	ret := _m.ctrl.Call(_m, "SyncServers")
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockServiceDiscoveryRecorder) SyncServers() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SyncServers")
}
