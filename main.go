package main

import (
    "context"
    "fmt"
    "log"

    "sqlc-demo/db"
    "github.com/jackc/pgx/v5/pgxpool"
)

func main() {
    // PostgreSQL connection string
    connString := "postgres://u0_a283@localhost:5432/mydb"
    
    // Connect to PostgreSQL
    conn, err := pgxpool.New(context.Background(), connString)
    if err != nil {
        log.Fatal("❌ Failed to connect to database:", err)
    }
    defer conn.Close()

    ctx := context.Background()

    // Test connection
    err = conn.Ping(ctx)
    if err != nil {
        log.Fatal("❌ Failed to ping database:", err)
    }
    fmt.Println("✅ Connected to PostgreSQL")

    // Create queries
    queries := db.New(conn)

    // 1. Create a user
    err = queries.CreateUser(ctx, db.CreateUserParams{
        Name:  "John Doe",
        Email: "john@example.com",
    })
    if err != nil {
        log.Fatal("❌ Failed to create user:", err)
    }
    fmt.Println("✅ User created successfully")

    // 2. Get the user
    user, err := queries.GetUser(ctx, 1)
    if err != nil {
        log.Fatal("❌ Failed to get user:", err)
    }
    fmt.Printf("👤 User found: ID=%d, Name=%s, Email=%s\n", user.ID, user.Name, user.Email)
    
    // Format the timestamp
    createdStr := "Unknown"
    if user.CreatedAt.Valid {
        createdStr = user.CreatedAt.Time.Format("2006-01-02 15:04:05")
    }
    fmt.Printf("   Created at: %s\n", createdStr)

    // 3. Create sample posts
    createSamplePosts(conn, user.ID)

    // 4. Get user's posts
    posts, err := queries.GetUserPosts(ctx, user.ID)
    if err != nil {
        log.Fatal("❌ Failed to get posts:", err)
    }
    fmt.Printf("📝 User has %d posts:\n", len(posts))
    for _, post := range posts {
        published := "no"
        if post.Published != nil && *post.Published {
            published = "yes"
        }
        fmt.Printf("   - %s (published: %s)\n", post.Title, published)
        
        // Add a comment to each post
        if post.ID == 1 {
            err = queries.CreateComment(ctx, db.CreateCommentParams{
                PostID:  post.ID,
                UserID:  user.ID,
                Content: "Great post! Keep it up!",
            })
            if err != nil {
                log.Printf("⚠️ Failed to add comment: %v", err)
            } else {
                fmt.Println("   💬 Added comment to post")
            }
        }
    }

    // 5. Get comments for first post
    comments, err := queries.GetPostComments(ctx, 1)
    if err != nil {
        log.Fatal("❌ Failed to get comments:", err)
    }
    fmt.Printf("💬 Comments on post 1: %d\n", len(comments))
    for _, comment := range comments {
        fmt.Printf("   - %s (user %d)\n", comment.Content, comment.UserID)
    }

    fmt.Println("✅ All operations completed successfully!")
}

func createSamplePosts(conn *pgxpool.Pool, userID int32) {
    posts := []struct {
        title     string
        content   string
        published bool
    }{
        {"My First Post", "This is the content of my first post.", true},
        {"Learning Go with SQLC", "SQLC makes database integration so easy!", false},
        {"Termux on Android", "Running Go applications on mobile is awesome.", true},
    }

    ctx := context.Background()
    for _, post := range posts {
        _, err := conn.Exec(ctx,
            "INSERT INTO posts (user_id, title, content, published) VALUES ($1, $2, $3, $4)",
            userID, post.title, post.content, post.published,
        )
        if err != nil {
            log.Printf("⚠️ Failed to create post '%s': %v", post.title, err)
        }
    }
    fmt.Println("✅ Sample posts created")
}
