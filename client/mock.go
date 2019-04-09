package client

import (
	"io"

	"github.com/stretchr/testify/mock"
)

// MockClient is mock implementation of the Client interface
type MockClient struct {
	mock.Mock
}

// DetectFaces is a mock implementation of the interface method
func (m *MockClient) DetectFaces(r io.Reader) (*Response, error) {
	args := m.Called(r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*Response), args.Error(1)
}
