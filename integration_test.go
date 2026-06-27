package main

import (
    "context"
    "fmt"
    "testing"
    "time"

    "sqlc-demo/db"
    "sqlc-demo/internal/testutil"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/stretchr/testify/suite"
)

// TestSuite is a test suite for integration tests
type TestSuite struct {
    suite.Suite
    conn    *pgxpool.Pool
    queries *db.Queries
    ctx     context.Context
}

func (s *TestSuite) SetupSuite() {
    // Setup test database once for all tests
    conn, _ := testutil.SetupTestDB(s.T())
    s.conn = conn
    s.ctx = context.Background()
    s.queries = db.New(conn)
    
    // Run migrations
    testutil.RunMigrations(s.T(), conn)
}

func (s *TestSuite) TearDownSuite() {
    // Cleanup after all tests
    if s.conn != nil {
        s.conn.Close()
    }
}

func (s *TestSuite) SetupTest() {
    // Clean database before each test
    testutil.CleanDB(s.T(), s.conn)
}

// TestCreateUser tests user creation
func (s *TestSuite) TestCreateUser() {
    err := s.queries.CreateUser(s.ctx, db.CreateUserParams{
        Name:  "Test User",
        Email: "test@example.com",
    })
    require.NoError(s.T(), err)
    
    user, err := s.queries.GetUser(s.ctx, 1)
    require.NoError(s.T(), err)
    
    assert.Equal(s.T(), "Test User", user.Name)
    assert.Equal(s.T(), "test@example.com", user.Email)
}

// TestUserWithPosts tests a user with posts
func (s *TestSuite) TestUserWithPosts() {
    // Create user
    err := s.queries.CreateUser(s.ctx, db.CreateUserParams{
        Name:  "John Doe",
        Email: "john@example.com",
    })
    require.NoError(s.T(), err)
    
    // Create post
    _, err = s.conn.Exec(s.ctx,
        "INSERT INTO posts (user_id, title, content, published) VALUES ($1, $2, $3, $4)",
        1, "Test Post", "Test Content", true,
    )
    require.NoError(s.T(), err)
    
    // Get user posts
    posts, err := s.queries.GetUserPosts(s.ctx, 1)
    require.NoError(s.T(), err)
    
    assert.Len(s.T(), posts, 1)
    assert.Equal(s.T(), "Test Post", posts[0].Title)
}

// TestDuplicateEmail tests unique email constraint
func (s *TestSuite) TestDuplicateEmail() {
    err := s.queries.CreateUser(s.ctx, db.CreateUserParams{
        Name:  "User 1",
        Email: "same@example.com",
    })
    require.NoError(s.T(), err)
    
    err = s.queries.CreateUser(s.ctx, db.CreateUserParams{
        Name:  "User 2",
        Email: "same@example.com",
    })
    assert.Error(s.T(), err, "Should fail due to duplicate email")
}

// TestCommentFlow tests the complete comment flow
func (s *TestSuite) TestCommentFlow() {
    // Create user
    err := s.queries.CreateUser(s.ctx, db.CreateUserParams{
        Name:  "Author",
        Email: "author@example.com",
    })
    require.NoError(s.T(), err)
    
    // Create post
    _, err = s.conn.Exec(s.ctx,
        "INSERT INTO posts (user_id, title, content, published) VALUES ($1, $2, $3, $4)",
        1, "Post Title", "Post Content", true,
    )
    require.NoError(s.T(), err)
    
    // Create comment
    err = s.queries.CreateComment(s.ctx, db.CreateCommentParams{
        PostID:  1,
        UserID:  1,
        Content: "Great post!",
    })
    require.NoError(s.T(), err)
    
    // Get comments
    comments, err := s.queries.GetPostComments(s.ctx, 1)
    require.NoError(s.T(), err)
    
    assert.Len(s.T(), comments, 1)
    assert.Equal(s.T(), "Great post!", comments[0].Content)
}

// TestDatabasePerformance tests database performance with unique emails
func (s *TestSuite) TestDatabasePerformance() {
    // Create 10 users with unique emails
    start := time.Now()
    
    for i := 0; i < 10; i++ {
        email := fmt.Sprintf("user%d@example.com", i)
        err := s.queries.CreateUser(s.ctx, db.CreateUserParams{
            Name:  fmt.Sprintf("User %d", i),
            Email: email,
        })
        require.NoError(s.T(), err, "Failed to create user %d with email %s", i, email)
    }
    
    elapsed := time.Since(start)
    
    // Verify 10 users were created
    users, err := s.queries.ListUsers(s.ctx)
    require.NoError(s.T(), err)
    assert.Len(s.T(), users, 10, "Should have 10 users")
    
    // Performance assertion
    assert.Less(s.T(), elapsed.Seconds(), 5.0, "Should create 10 users within 5 seconds")
    s.T().Logf("✅ Created 10 users in %v", elapsed)
}

// Run test suite
func TestIntegrationSuite(t *testing.T) {
    suite.Run(t, new(TestSuite))
}
