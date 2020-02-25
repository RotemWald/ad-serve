package main

import (
	"fmt"
	unitTest "github.com/Valiben/gin_unit_test"
	"github.com/Valiben/gin_unit_test/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {
	db = NewDatabaseConnection("sqlite3", "spotim-test.db")
	db.AutoMigrate(&AdConfiguration{})

	redisPool = NewRedisPool(80, 12000, "localhost:6379")

	r := gin.Default()
	r.GET("/config", GetConfig)
	r.GET("/configs", GetConfigs)
	r.POST("/configs", CreateConfig)
	r.PUT("/configs/:id", UpdateConfig)
	r.DELETE("/configs/:id", DeleteConfig)

	unitTest.SetRouter(r)
}

func TestAll(t *testing.T) {
	defer db.Close()

	config := AdConfiguration{
		PartnerId: 1,
		AdCode:    "abcd-1234",
		Component: "CONVERSATION",
		Platform:  "NATIVE",
		Location:  "IL",
	}

	configId := testCreateConfig(t, config)
	testGetConfigs(t, configId)
	testGetConfig(t, configId, config)
	testUpdateConfig(t, configId, config)
	testDeleteConfig(t, configId)
}

func testCreateConfig(t *testing.T, config AdConfiguration) uint {
	err := unitTest.TestHandlerUnMarshalResp(utils.POST, "/configs", utils.JSON, config, &config)
	if err != nil {
		t.Errorf("testCreateConfig: %v\n", err)
		return 0
	}
	if config.ID != 0 {
		fmt.Println("testCreateConfig passed")
		return config.ID
	}

	assert.Fail(t, "Config could not be created")
	return 0
}

func testGetConfigs(t *testing.T, id uint) {
	var configs []AdConfiguration
	err := unitTest.TestHandlerUnMarshalResp(utils.GET, "/configs", utils.FORM, nil, &configs)
	if err != nil {
		t.Errorf("testGetConfig: %v\n", err)
		return
	}

	for _, config := range configs {
		if config.ID == id {
			fmt.Println("testGetConfigs passed")
			return
		}
	}

	assert.Fail(t, "Config could not be found")
	return
}

func testGetConfig(t *testing.T, id uint, config AdConfiguration) {
	param := make(map[string]interface{})
	param["component"] = config.Component
	param["platform"] = config.Platform
	param["location"] = config.Location

	err := unitTest.TestHandlerUnMarshalResp(utils.GET, "/config", utils.FORM, param, &config)
	if err != nil {
		t.Errorf("testGetConfig: %v\n", err)
		return
	}

	if config.ID == id {
		fmt.Println("testGetConfig passed")
		return
	}

	assert.Fail(t, "Config could not be served")
	return
}

func testUpdateConfig(t *testing.T, id uint, config AdConfiguration) {
	config.ID = id
	config.Location = "USA"

	err := unitTest.TestHandlerUnMarshalResp(utils.PUT, "/configs/"+fmt.Sprint(id), utils.JSON, config, &config)
	if err != nil {
		t.Errorf("testUpdateConfig: %v\n", err)
		return
	}
	if config.ID == id && config.Location == "USA" {
		fmt.Println("testUpdateConfig passed")
		return
	}

	assert.Fail(t, "Config could not be updated")
	return
}

func testDeleteConfig(t *testing.T, id uint) {
	body, err := unitTest.TestOrdinaryHandler(utils.DELETE, "/configs/"+fmt.Sprint(id), utils.FORM, nil)
	if err != nil {
		t.Errorf("testDeleteConfig: %v\n", err)
		return
	}

	if len(body) == 0 {
		fmt.Println("testDeleteConfig passed")
		return
	}

	assert.Fail(t, "Config could not be deleted")
	return
}
