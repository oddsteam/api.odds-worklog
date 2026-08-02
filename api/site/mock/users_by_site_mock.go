package mock_site

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	models "gitlab.odds.team/worklog/api.odds-worklog/business/models"
)

// MockForGettingUsersBySiteID is a mock of ForGettingUsersBySiteID interface.
type MockForGettingUsersBySiteID struct {
	ctrl     *gomock.Controller
	recorder *MockForGettingUsersBySiteIDMockRecorder
}

// MockForGettingUsersBySiteIDMockRecorder is the mock recorder for MockForGettingUsersBySiteID.
type MockForGettingUsersBySiteIDMockRecorder struct {
	mock *MockForGettingUsersBySiteID
}

// NewMockForGettingUsersBySiteID creates a new mock instance.
func NewMockForGettingUsersBySiteID(ctrl *gomock.Controller) *MockForGettingUsersBySiteID {
	mock := &MockForGettingUsersBySiteID{ctrl: ctrl}
	mock.recorder = &MockForGettingUsersBySiteIDMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockForGettingUsersBySiteID) EXPECT() *MockForGettingUsersBySiteIDMockRecorder {
	return m.recorder
}

// GetBySiteID mocks base method.
func (m *MockForGettingUsersBySiteID) GetBySiteID(id string) ([]*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBySiteID", id)
	ret0, _ := ret[0].([]*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBySiteID indicates an expected call of GetBySiteID.
func (mr *MockForGettingUsersBySiteIDMockRecorder) GetBySiteID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBySiteID", reflect.TypeOf((*MockForGettingUsersBySiteID)(nil).GetBySiteID), id)
}
