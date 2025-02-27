package testutils

import "github.com/stretchr/testify/mock"

type MockStorage struct {
	mock.Mock
}

func NewMockStorage() *MockStorage {
	return &MockStorage{}
}

func (m *MockStorage) Post(url string) (string, error) {
	args := m.Called(url)
	return args.String(0), args.Error(1)
}

func (m *MockStorage) Get(shortURL string) (string, error) {
	args := m.Called(shortURL)
	return args.String(0), args.Error(1)
}

func (m *MockStorage) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockStorage) ExpectPost(url string, response string, err error) {
	m.On("Post", url).Return(response, err)
}

func (m *MockStorage) ExpectGet(shortURL string, response string, err error) {
	m.On("Get", shortURL).Return(response, err)
}
