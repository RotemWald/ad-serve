package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var (
	db        *gorm.DB
	redisPool *redis.Pool
)

func main() {
	db = NewDatabaseConnection("sqlite3", "spotim.db")
	defer db.Close()
	db.AutoMigrate(&AdConfiguration{})

	redisPool = NewRedisPool(80, 12000, "localhost:6379")

	r := gin.Default()
	r.GET("/config", GetConfig)
	r.GET("/configs", GetConfigs)
	r.POST("/configs", CreateConfig)
	r.PUT("/configs/:id", UpdateConfig)
	r.DELETE("/configs/:id", DeleteConfig)
	r.Run()
}
