package main

import (
    "context"
    "github.com/gin-gonic/gin"
    "sqlc-demo/db"
    "github.com/jackc/pgx/v5/pgxpool"
)

func main() {
    conn := // connect to db
    queries := db.New(conn)
    
    r := gin.Default()
    r.GET("/users", func(c *gin.Context) {
        users, _ := queries.ListUsers(context.Background())
        c.JSON(200, users)
    })
    r.Run(":8080")
}
