package mockrepo

import (
	"github.com/stretchr/testify/mock"
)

type MockRepo struct {
	mock.Mock
}

func NewMockRepo() *MockRepo {
	return &MockRepo{}
}
