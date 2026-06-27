package main

import (
    "context"
    "net/http"
    "strconv"

    "sqlc-demo/db"
    "github.com/gin-gonic/gin"
    "github.com/jackc/pgx/v5/pgxpool"
)

func main() {
    // Database connection
    connString := "postgres://u0_a283@localhost:5432/mydb"
    conn, err := pgxpool.New(context.Background(), connString)
    if err != nil {
        panic(err)
    }
    defer conn.Close()

    queries := db.New(conn)
    ctx := context.Background()

    r := gin.Default()

    // Health check
    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })

    // User routes
    r.GET("/users", func(c *gin.Context) {
        users, err := queries.ListUsers(ctx)
        if err != nil {
            c.JSON(500, gin.H{"error": err.Error()})
            return
        }
        c.JSON(200, users)
    })

    r.GET("/users/:id", func(c *gin.Context) {
        id, err := strconv.Atoi(c.Param("id"))
        if err != nil {
            c.JSON(400, gin.H{"error": "Invalid user ID"})
            return
        }
        user, err := queries.GetUser(ctx, int32(id))
        if err != nil {
            c.JSON(404, gin.H{"error": "User not found"})
            return
        }
        c.JSON(200, user)
    })

    r.POST("/users", func(c *gin.Context) {
        var params db.CreateUserParams
        if err := c.BindJSON(&params); err != nil {
            c.JSON(400, gin.H{"error": err.Error()})
            return
        }
        if err := queries.CreateUser(ctx, params); err != nil {
            c.JSON(500, gin.H{"error": err.Error()})
            return
        }
        c.JSON(201, gin.H{"message": "User created"})
    })

    r.GET("/users/:id/posts", func(c *gin.Context) {
        id, err := strconv.Atoi(c.Param("id"))
        if err != nil {
            c.JSON(400, gin.H{"error": "Invalid user ID"})
            return
        }
        posts, err := queries.GetUserPosts(ctx, int32(id))
        if err != nil {
            c.JSON(500, gin.H{"error": err.Error()})
            return
        }
        c.JSON(200, posts)
    })

    // Post routes
    r.GET("/posts/:id/comments", func(c *gin.Context) {
        id, err := strconv.Atoi(c.Param("id"))
        if err != nil {
            c.JSON(400, gin.H{"error": "Invalid post ID"})
            return
        }
        comments, err := queries.GetPostComments(ctx, int32(id))
        if err != nil {
            c.JSON(500, gin.H{"error": err.Error()})
            return
        }
        c.JSON(200, comments)
    })

    r.POST("/comments", func(c *gin.Context) {
        var params db.CreateCommentParams
        if err := c.BindJSON(&params); err != nil {
            c.JSON(400, gin.H{"error": err.Error()})
            return
        }
        if err := queries.CreateComment(ctx, params); err != nil {
            c.JSON(500, gin.H{"error": err.Error()})
            return
        }
        c.JSON(201, gin.H{"message": "Comment added"})
    })

    // Start server
    r.Run(":8080")
}
