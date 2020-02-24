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

	redisconn := redisPool.Get()
	defer redisconn.Close()

	configIds, err := redis.Strings(redisconn.Do("SINTER", componentKey, platformKey, locationKey))
	if err != nil {
		c.AbortWithStatus(500)
		fmt.Println(err)
		return
	}

	if len(configIds) > 0 {
		configIdToServe = configIds[0]
	} else {
		configIds, err = redis.Strings(redisconn.Do("SINTER", componentKey, platformKey))
		if err != nil {
			c.AbortWithStatus(500)
			fmt.Println(err)
			return
		}

		if len(configIds) > 0 {
			configIdToServe = configIds[0]
		} else {
			configIds, err = redis.Strings(redisconn.Do("SINTER", componentKey, locationKey))
			if err != nil {
				c.AbortWithStatus(500)
				fmt.Println(err)
				return
			}

			if len(configIds) > 0 {
				configIdToServe = configIds[0]
			} else {
				configIds, err = redis.Strings(redisconn.Do("SMEMBERS", componentKey))
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
		c.JSON(200, nil)
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

	redisconn := redisPool.Get()
	defer redisconn.Close()

	_, err := redisconn.Do("SADD", "component:"+config.Component, config.ID)
	if err != nil {
		c.AbortWithStatus(500)
		fmt.Println(err)
		return
	}
	_, err = redisconn.Do("SADD", "platform:"+config.Platform, config.ID)
	if err != nil {
		c.AbortWithStatus(500)
		fmt.Println(err)
		return
	}
	_, err = redisconn.Do("SADD", "location:"+config.Location, config.ID)
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

	redisconn := redisPool.Get()
	defer redisconn.Close()

	_, err := redisconn.Do("SREM", "component:"+config.Component, config.ID)
	if err != nil {
		c.AbortWithStatus(500)
		fmt.Println(err)
		return
	}
	_, err = redisconn.Do("SREM", "platform:"+config.Platform, config.ID)
	if err != nil {
		c.AbortWithStatus(500)
		fmt.Println(err)
		return
	}
	_, err = redisconn.Do("SREM", "location:"+config.Location, config.ID)
	if err != nil {
		c.AbortWithStatus(500)
		fmt.Println(err)
		return
	}

	c.BindJSON(&config)
	db.Save(&config)

	_, err = redisconn.Do("SADD", "component:"+config.Component, config.ID)
	if err != nil {
		c.AbortWithStatus(500)
		fmt.Println(err)
		return
	}
	_, err = redisconn.Do("SADD", "platform:"+config.Platform, config.ID)
	if err != nil {
		c.AbortWithStatus(500)
		fmt.Println(err)
		return
	}
	_, err = redisconn.Do("SADD", "location:"+config.Location, config.ID)
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

	redisconn := redisPool.Get()
	defer redisconn.Close()

	_, err := redisconn.Do("SREM", "component:"+config.Component, config.ID)
	if err != nil {
		c.AbortWithStatus(500)
		fmt.Println(err)
		return
	}
	_, err = redisconn.Do("SREM", "platform:"+config.Platform, config.ID)
	if err != nil {
		c.AbortWithStatus(500)
		fmt.Println(err)
		return
	}
	_, err = redisconn.Do("SREM", "location:"+config.Location, config.ID)
	if err != nil {
		c.AbortWithStatus(500)
		fmt.Println(err)
		return
	}

	c.JSON(204, nil)
}
