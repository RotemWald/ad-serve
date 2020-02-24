package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB
var err error

type AdConfiguration struct {
	ID        uint   `json:"id"`
	PartnerId uint   `json:"partnerid"`
	AdCode    string `json:"adcode"`
	Component string `json:"component"`
	Platform  string `json:"platform"`
	Location  string `json:"location"`
}

func main() {
	db, err = gorm.Open("sqlite3", "spotim.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	db.AutoMigrate(&AdConfiguration{})

	//db.Create(&AdConfiguration{
	//	PartnerId: 1,
	//	AdCode:    "abcd-1234",
	//	Component: "CONVERSATION",
	//	Platform:  "DESKTOP",
	//	Location:  "USA",
	//})
	//db.Create(&AdConfiguration{
	//	PartnerId: 2,
	//	AdCode:    "zxcv-5678",
	//	Component: "REACTIONS",
	//	Platform:  "NATIVE",
	//	Location:  "IL",
	//})

	r := gin.Default()
	r.GET("/configs", GetConfigs)
	r.POST("/configs", CreateConfig)
	r.PUT("/configs/:id", UpdateConfig)
	r.DELETE("/configs/:id", DeleteConfig)
	r.Run()
}

func GetConfigs(c *gin.Context) {
	var configs []AdConfiguration
	if err := db.Find(&configs).Error; err != nil {
		c.AbortWithStatus(500)
		fmt.Println(err)
	} else {
		c.JSON(200, configs)
	}
}

func CreateConfig(c *gin.Context) {
	var config AdConfiguration
	c.BindJSON(&config)
	db.Create(&config)
	c.JSON(200, config)
}

func UpdateConfig(c *gin.Context) {
	var config AdConfiguration
	id := c.Param("id")
	if err := db.First(&config, id).Error; err != nil {
		c.AbortWithStatus(400)
		fmt.Println(err)
	} else {
		c.BindJSON(&config)
		db.Save(&config)
		c.JSON(200, config)
	}
}

func DeleteConfig(c *gin.Context) {
	var config AdConfiguration
	id := c.Param("id")
	if err := db.First(&config, id).Error; err != nil {
		c.AbortWithStatus(400)
		fmt.Println(err)
	} else {
		db.Delete(&config)
		c.JSON(204, nil)
	}
}
