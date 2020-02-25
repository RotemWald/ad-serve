package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
)

func GetConfig(c *gin.Context) {
	configIdToServe := ""
	componentKey := "component:" + c.Query("component")
	platformKey := "platform:" + c.Query("platform")
	locationKey := "location:" + c.Query("location")

	redisConn := redisPool.Get()
	defer redisConn.Close()

	configIds, err := redis.Strings(redisConn.Do("SINTER", componentKey, platformKey, locationKey))
	if err != nil {
		c.AbortWithStatus(500)
		fmt.Println(err)
		return
	}

	if len(configIds) > 0 {
		configIdToServe = configIds[0]
	} else {
		configIds, err = redis.Strings(redisConn.Do("SINTER", componentKey, platformKey))
		if err != nil {
			c.AbortWithStatus(500)
			fmt.Println(err)
			return
		}

		if len(configIds) > 0 {
			configIdToServe = configIds[0]
		} else {
			configIds, err = redis.Strings(redisConn.Do("SINTER", componentKey, locationKey))
			if err != nil {
				c.AbortWithStatus(500)
				fmt.Println(err)
				return
			}

			if len(configIds) > 0 {
				configIdToServe = configIds[0]
			} else {
				configIds, err = redis.Strings(redisConn.Do("SMEMBERS", componentKey))
				if err != nil {
					c.AbortWithStatus(500)
					fmt.Println(err)
					return
				}

				if len(configIds) > 0 {
					configIdToServe = configIds[0]
				}
			}
		}
	}

	if configIdToServe != "" {
		var config AdConfiguration
		db.First(&config, configIdToServe)
		c.JSON(200, config)
	} else {
		c.JSON(204, nil)
	}
}

func GetConfigs(c *gin.Context) {
	var configs []AdConfiguration
	if err := db.Find(&configs).Error; err != nil {
		c.AbortWithStatus(500)
		fmt.Println(err)
		return
	}
	c.JSON(200, configs)
}

func CreateConfig(c *gin.Context) {
	var config AdConfiguration
	c.BindJSON(&config)
	db.Create(&config)

	redisConn := redisPool.Get()
	defer redisConn.Close()

	_, err := redisConn.Do("SADD", "component:"+config.Component, config.ID)
	if err != nil {
		c.AbortWithStatus(500)
		fmt.Println(err)
		return
	}
	_, err = redisConn.Do("SADD", "platform:"+config.Platform, config.ID)
	if err != nil {
		c.AbortWithStatus(500)
		fmt.Println(err)
		return
	}
	_, err = redisConn.Do("SADD", "location:"+config.Location, config.ID)
	if err != nil {
		c.AbortWithStatus(500)
		fmt.Println(err)
		return
	}

	c.JSON(200, config)
}

func UpdateConfig(c *gin.Context) {
	var config AdConfiguration
	id := c.Param("id")
	if err := db.First(&config, id).Error; err != nil {
		c.AbortWithStatus(400)
		fmt.Println(err)
		return
	}

	redisConn := redisPool.Get()
	defer redisConn.Close()

	_, err := redisConn.Do("SREM", "component:"+config.Component, config.ID)
	if err != nil {
		c.AbortWithStatus(500)
		fmt.Println(err)
		return
	}
	_, err = redisConn.Do("SREM", "platform:"+config.Platform, config.ID)
	if err != nil {
		c.AbortWithStatus(500)
		fmt.Println(err)
		return
	}
	_, err = redisConn.Do("SREM", "location:"+config.Location, config.ID)
	if err != nil {
		c.AbortWithStatus(500)
		fmt.Println(err)
		return
	}

	c.BindJSON(&config)
	db.Save(&config)

	_, err = redisConn.Do("SADD", "component:"+config.Component, config.ID)
	if err != nil {
		c.AbortWithStatus(500)
		fmt.Println(err)
		return
	}
	_, err = redisConn.Do("SADD", "platform:"+config.Platform, config.ID)
	if err != nil {
		c.AbortWithStatus(500)
		fmt.Println(err)
		return
	}
	_, err = redisConn.Do("SADD", "location:"+config.Location, config.ID)
	if err != nil {
		c.AbortWithStatus(500)
		fmt.Println(err)
		return
	}

	c.JSON(200, config)
}

func DeleteConfig(c *gin.Context) {
	var config AdConfiguration
	id := c.Param("id")
	if err := db.First(&config, id).Error; err != nil {
		c.AbortWithStatus(400)
		fmt.Println(err)
		return
	}

	db.Delete(&config)

	redisConn := redisPool.Get()
	defer redisConn.Close()

	_, err := redisConn.Do("SREM", "component:"+config.Component, config.ID)
	if err != nil {
		c.AbortWithStatus(500)
		fmt.Println(err)
		return
	}
	_, err = redisConn.Do("SREM", "platform:"+config.Platform, config.ID)
	if err != nil {
		c.AbortWithStatus(500)
		fmt.Println(err)
		return
	}
	_, err = redisConn.Do("SREM", "location:"+config.Location, config.ID)
	if err != nil {
		c.AbortWithStatus(500)
		fmt.Println(err)
		return
	}

	c.JSON(204, nil)
}
