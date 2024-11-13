// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/service/todo.go
//
// Generated by this command:
//
//	mockgen -source=./internal/service/todo.go -destination=test/mock/./service/todo.go
//

// Package mock_service is a generated GoMock package.
package mock_service

import (
	context "context"
	reflect "reflect"

	entity "github.com/sherwin-77/golang-todos/internal/entity"
	dto "github.com/sherwin-77/golang-todos/internal/http/dto"
	gomock "go.uber.org/mock/gomock"
)

// MockTodoService is a mock of TodoService interface.
type MockTodoService struct {
	ctrl     *gomock.Controller
	recorder *MockTodoServiceMockRecorder
	isgomock struct{}
}

// MockTodoServiceMockRecorder is the mock recorder for MockTodoService.
type MockTodoServiceMockRecorder struct {
	mock *MockTodoService
}

// NewMockTodoService creates a new mock instance.
func NewMockTodoService(ctrl *gomock.Controller) *MockTodoService {
	mock := &MockTodoService{ctrl: ctrl}
	mock.recorder = &MockTodoServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTodoService) EXPECT() *MockTodoServiceMockRecorder {
	return m.recorder
}

// CreateTodo mocks base method.
func (m *MockTodoService) CreateTodo(ctx context.Context, request dto.TodoRequest, userID string) (*entity.Todo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTodo", ctx, request, userID)
	ret0, _ := ret[0].(*entity.Todo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTodo indicates an expected call of CreateTodo.
func (mr *MockTodoServiceMockRecorder) CreateTodo(ctx, request, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTodo", reflect.TypeOf((*MockTodoService)(nil).CreateTodo), ctx, request, userID)
}

// DeleteTodo mocks base method.
func (m *MockTodoService) DeleteTodo(ctx context.Context, id, userID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTodo", ctx, id, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteTodo indicates an expected call of DeleteTodo.
func (mr *MockTodoServiceMockRecorder) DeleteTodo(ctx, id, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTodo", reflect.TypeOf((*MockTodoService)(nil).DeleteTodo), ctx, id, userID)
}

// GetTodoByID mocks base method.
func (m *MockTodoService) GetTodoByID(ctx context.Context, id, userID string) (*entity.Todo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTodoByID", ctx, id, userID)
	ret0, _ := ret[0].(*entity.Todo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTodoByID indicates an expected call of GetTodoByID.
func (mr *MockTodoServiceMockRecorder) GetTodoByID(ctx, id, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTodoByID", reflect.TypeOf((*MockTodoService)(nil).GetTodoByID), ctx, id, userID)
}

// GetTodosByUserID mocks base method.
func (m *MockTodoService) GetTodosByUserID(ctx context.Context, userID string) ([]entity.Todo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTodosByUserID", ctx, userID)
	ret0, _ := ret[0].([]entity.Todo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTodosByUserID indicates an expected call of GetTodosByUserID.
func (mr *MockTodoServiceMockRecorder) GetTodosByUserID(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTodosByUserID", reflect.TypeOf((*MockTodoService)(nil).GetTodosByUserID), ctx, userID)
}

// UpdateTodo mocks base method.
func (m *MockTodoService) UpdateTodo(ctx context.Context, request dto.UpdateTodoRequest, userID string) (*entity.Todo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTodo", ctx, request, userID)
	ret0, _ := ret[0].(*entity.Todo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateTodo indicates an expected call of UpdateTodo.
func (mr *MockTodoServiceMockRecorder) UpdateTodo(ctx, request, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTodo", reflect.TypeOf((*MockTodoService)(nil).UpdateTodo), ctx, request, userID)
}