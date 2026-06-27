package main

import (
    "context"
    "testing"

    "sqlc-demo/db"
    "sqlc-demo/internal/testutil"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestCreateAndGetUser(t *testing.T) {
    // Setup test database
    conn, cleanup := testutil.SetupTestDB(t)
    defer cleanup()
    
    ctx := context.Background()
    
    // Run migrations
    testutil.RunMigrations(t, conn)
    defer testutil.CleanDB(t, conn)
    
    // Create queries
    queries := db.New(conn)
    
    // Create user
    err := queries.CreateUser(ctx, db.CreateUserParams{
        Name:  "Test User",
        Email: "test@example.com",
    })
    require.NoError(t, err, "Failed to create user")
    
    // Get user
    user, err := queries.GetUser(ctx, 1)
    require.NoError(t, err, "Failed to get user")
    
    // Assertions
    assert.Equal(t, int32(1), user.ID, "User ID should be 1")
    assert.Equal(t, "Test User", user.Name, "User name should match")
    assert.Equal(t, "test@example.com", user.Email, "User email should match")
    assert.NotNil(t, user.CreatedAt, "CreatedAt should not be nil")
    assert.True(t, user.CreatedAt.Valid, "CreatedAt should be valid")
}

func TestUpdateUserEmail(t *testing.T) {
    // Setup
    conn, cleanup := testutil.SetupTestDB(t)
    defer cleanup()
    
    ctx := context.Background()
    testutil.RunMigrations(t, conn)
    defer testutil.CleanDB(t, conn)
    
    queries := db.New(conn)
    
    // Create user
    err := queries.CreateUser(ctx, db.CreateUserParams{
        Name:  "John Doe",
        Email: "john@example.com",
    })
    require.NoError(t, err)
    
    // Update email
    err = queries.UpdateUserEmail(ctx, db.UpdateUserEmailParams{
        ID:    1,
        Email: "john.doe@example.com",
    })
    require.NoError(t, err)
    
    // Get updated user
    user, err := queries.GetUser(ctx, 1)
    require.NoError(t, err)
    
    // Assert email was updated
    assert.Equal(t, "john.doe@example.com", user.Email)
}

func TestListUsers(t *testing.T) {
    // Setup
    conn, cleanup := testutil.SetupTestDB(t)
    defer cleanup()
    
    ctx := context.Background()
    testutil.RunMigrations(t, conn)
    defer testutil.CleanDB(t, conn)
    
    queries := db.New(conn)
    
    // Create multiple users
    users := []struct {
        name  string
        email string
    }{
        {"Alice", "alice@example.com"},
        {"Bob", "bob@example.com"},
        {"Charlie", "charlie@example.com"},
    }
    
    for _, u := range users {
        err := queries.CreateUser(ctx, db.CreateUserParams{
            Name:  u.name,
            Email: u.email,
        })
        require.NoError(t, err)
    }
    
    // List users
    allUsers, err := queries.ListUsers(ctx)
    require.NoError(t, err)
    
    // Assert - users are ordered by created_at DESC
    assert.Len(t, allUsers, 3, "Should have 3 users")
    // The first user should be the most recently created (Charlie)
    assert.Equal(t, "Charlie", allUsers[0].Name, "Most recent user should be first")
    assert.Equal(t, "Bob", allUsers[1].Name, "Second user should be Bob")
    assert.Equal(t, "Alice", allUsers[2].Name, "Third user should be Alice")
}

func TestCreateAndGetPosts(t *testing.T) {
    // Setup
    conn, cleanup := testutil.SetupTestDB(t)
    defer cleanup()
    
    ctx := context.Background()
    testutil.RunMigrations(t, conn)
    defer testutil.CleanDB(t, conn)
    
    queries := db.New(conn)
    
    // Create user
    err := queries.CreateUser(ctx, db.CreateUserParams{
        Name:  "Test User",
        Email: "test@example.com",
    })
    require.NoError(t, err)
    
    // Create posts
    posts := []struct {
        userID    int32
        title     string
        content   string
        published bool
    }{
        {1, "First Post", "Content 1", true},
        {1, "Second Post", "Content 2", false},
    }
    
    for _, p := range posts {
        _, err := conn.Exec(ctx,
            "INSERT INTO posts (user_id, title, content, published) VALUES ($1, $2, $3, $4)",
            p.userID, p.title, p.content, p.published,
        )
        require.NoError(t, err)
    }
    
    // Get user posts
    userPosts, err := queries.GetUserPosts(ctx, 1)
    require.NoError(t, err)
    
    // Assert
    assert.Len(t, userPosts, 2, "Should have 2 posts")
    assert.Equal(t, "First Post", userPosts[0].Title)
}
