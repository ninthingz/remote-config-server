package main

import (
	"github.com/gin-gonic/gin"
)

var configMap = make(map[string]*Config)

func main() {
	var err error

	err = openDatabase()
	if err != nil {
		println(err.Error())
		return
	}

	err = loadConfig()

	if err != nil {
		println(err.Error())
		return
	}

	r := gin.Default()
	r.POST("/config", createConfig)
	err = r.Run()
	if err != nil {
		println(err.Error())
		return
	}
}

func createConfig(context *gin.Context) {
	var config Config
	err := context.BindJSON(&config)
	if err != nil {
		context.JSON(200, gin.H{"code": 500, "message": err.Error()})
		return
	}

	if config.Name == "" {
		context.JSON(200, gin.H{"code": 500, "message": "name is required"})
		return
	}

	if configMap[config.Name] != nil {
		context.JSON(200, gin.H{"code": 500, "message": "name already exists"})
		return
	}

	err = configDao.create(config)
	if err != nil {
		context.JSON(200, gin.H{"code": 500, "message": err.Error()})
		return
	}

	configMap[config.Name] = &config

	context.JSON(200, gin.H{"code": 200, "message": "success"})
}

func loadConfig() error {
	configList, err := configDao.list()
	if err != nil {
		return err
	}
	for _, config := range configList {
		configMap[config.Name] = config
	}
	return nil
}
