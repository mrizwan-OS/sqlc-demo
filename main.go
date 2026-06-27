package main

import (
    "context"
    "fmt"
    "log"

    "sqlc-demo/db"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/jackc/pgx/v5/pgconn"
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

    // 1. Check if user exists, create if not
    var userID int32 = 1
    user, err := queries.GetUser(ctx, userID)
    
    if err != nil {
        // User doesn't exist, create one
        fmt.Println("📝 Creating new user...")
        err = queries.CreateUser(ctx, db.CreateUserParams{
            Name:  "John Doe",
            Email: "john@example.com",
        })
        if err != nil {
            // Check if it's a duplicate key error
            if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
                fmt.Println("⚠️ User already exists, skipping creation")
            } else {
                log.Fatal("❌ Failed to create user:", err)
            }
        } else {
            fmt.Println("✅ User created successfully")
        }
        
        // Get the user again
        user, err = queries.GetUser(ctx, userID)
        if err != nil {
            log.Fatal("❌ Failed to get user:", err)
        }
    } else {
        fmt.Printf("👤 User already exists: %s\n", user.Name)
    }

    fmt.Printf("👤 User found: ID=%d, Name=%s, Email=%s\n", user.ID, user.Name, user.Email)
    
    // Format the timestamp
    createdStr := "Unknown"
    if user.CreatedAt.Valid {
        createdStr = user.CreatedAt.Time.Format("2006-01-02 15:04:05")
    }
    fmt.Printf("   Created at: %s\n", createdStr)

    // 2. Create sample posts (check if posts exist first)
    posts, err := queries.GetUserPosts(ctx, user.ID)
    if err != nil {
        log.Printf("⚠️ Failed to get posts: %v", err)
    }
    
    if len(posts) == 0 {
        fmt.Println("📝 Creating sample posts...")
        createSamplePosts(conn, user.ID)
        posts, err = queries.GetUserPosts(ctx, user.ID)
        if err != nil {
            log.Fatal("❌ Failed to get posts:", err)
        }
    }

    fmt.Printf("📝 User has %d posts:\n", len(posts))
    for _, post := range posts {
        published := "no"
        if post.Published != nil && *post.Published {
            published = "yes"
        }
        fmt.Printf("   - %s (published: %s)\n", post.Title, published)
        
        // Add a comment to each post if not exists
        if post.ID == 1 {
            comments, err := queries.GetPostComments(ctx, post.ID)
            if err != nil {
                log.Printf("⚠️ Failed to get comments: %v", err)
            }
            
            if len(comments) == 0 {
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
    }

    // 3. Get comments for first post
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
