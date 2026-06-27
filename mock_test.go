package main

import (
    "context"
    "testing"

    "sqlc-demo/db"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

// MockQuerier is a mock implementation of the Querier interface
type MockQuerier struct {
    mock.Mock
}

func (m *MockQuerier) CreateUser(ctx context.Context, arg db.CreateUserParams) error {
    args := m.Called(ctx, arg)
    return args.Error(0)
}

func (m *MockQuerier) GetUser(ctx context.Context, id int32) (db.User, error) {
    args := m.Called(ctx, id)
    return args.Get(0).(db.User), args.Error(1)
}

func (m *MockQuerier) ListUsers(ctx context.Context) ([]db.User, error) {
    args := m.Called(ctx)
    return args.Get(0).([]db.User), args.Error(1)
}

func (m *MockQuerier) GetUserPosts(ctx context.Context, userID int32) ([]db.Post, error) {
    args := m.Called(ctx, userID)
    return args.Get(0).([]db.Post), args.Error(1)
}

func (m *MockQuerier) UpdateUserEmail(ctx context.Context, arg db.UpdateUserEmailParams) error {
    args := m.Called(ctx, arg)
    return args.Error(0)
}

func (m *MockQuerier) DeleteUser(ctx context.Context, id int32) error {
    args := m.Called(ctx, id)
    return args.Error(0)
}

// Mock tests
func TestMockUserCreation(t *testing.T) {
    mockQuerier := new(MockQuerier)
    
    // Setup expectations
    mockQuerier.On("CreateUser", mock.Anything, db.CreateUserParams{
        Name:  "Test User",
        Email: "test@example.com",
    }).Return(nil)
    
    mockQuerier.On("GetUser", mock.Anything, int32(1)).Return(db.User{
        ID:    1,
        Name:  "Test User",
        Email: "test@example.com",
    }, nil)
    
    // Test with mock
    ctx := context.Background()
    err := mockQuerier.CreateUser(ctx, db.CreateUserParams{
        Name:  "Test User",
        Email: "test@example.com",
    })
    assert.NoError(t, err)
    
    user, err := mockQuerier.GetUser(ctx, 1)
    assert.NoError(t, err)
    assert.Equal(t, "Test User", user.Name)
    assert.Equal(t, "test@example.com", user.Email)
    
    // Assert expectations were met
    mockQuerier.AssertExpectations(t)
}
