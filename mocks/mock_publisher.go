package mocks

import (
	"financial-backend/pkg/config"
	"github.com/stretchr/testify/mock"
)

// MockPublisher is a mock for config.Publisher
type MockPublisher struct {
	mock.Mock
}

// RegisterHandler mocks the RegisterHandler method
func (m *MockPublisher) RegisterHandler(h config.Handler) {
	m.Called(h)
}

// Publish mocks the Publish method
func (m *MockPublisher) Publish(event config.Event) {
	m.Called(event)
}
