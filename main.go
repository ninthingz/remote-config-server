package main

import (
	"github.com/gin-gonic/gin"
	"strings"
)

var configMap = make(map[string]*Config)

var sign = "0D000721"

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
	r.Use(Cors())
	r.Use(Authorization())
	r.POST("/config", createConfig)
	r.PATCH("/config", updateConfig)
	r.DELETE("/config", deleteConfig)
	r.GET("/config", listConfig)
	r.GET("/config/element", getConfig)
	err = r.Run(":19200")
	if err != nil {
		println(err.Error())
		return
	}
}

func getConfig(context *gin.Context) {
	name := context.Query("name")
	if name == "" {
		context.JSON(200, gin.H{"code": 500, "message": "name is required"})
		return
	}
	if configMap[name] == nil {
		context.JSON(200, gin.H{"code": 500, "message": "config not found"})
		return
	}
	context.JSON(200, gin.H{"code": 200, "data": configMap[name]})
}

func listConfig(context *gin.Context) {
	keyword := context.Query("keyword")
	if keyword == "" {
		var configList = make([]*Config, 0)
		for _, config := range configMap {
			configList = append(configList, config)
		}
		context.JSON(200, gin.H{"code": 200, "data": configList})
		return
	}

	var configList = make([]*Config, 0)
	for _, config := range configMap {
		if strings.Contains(config.Name, keyword) {
			configList = append(configList, config)
		}
	}

	context.JSON(200, gin.H{"code": 200, "data": configList})
	return

}

func deleteConfig(context *gin.Context) {
	name := context.Query("name")
	if name == "" {
		context.JSON(200, gin.H{"code": 500, "message": "name is required"})
		return
	}
	if configMap[name] == nil {
		context.JSON(200, gin.H{"code": 500, "message": "config not found"})
		return
	}
	err := configDao.delete(configMap[name].Id)
	if err != nil {
		context.JSON(200, gin.H{"code": 500, "message": err.Error()})
		return
	}
	delete(configMap, name)
	context.JSON(200, gin.H{"code": 200, "message": "success"})
}

func updateConfig(context *gin.Context) {
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

	if configMap[config.Name] == nil {
		context.JSON(200, gin.H{"code": 500, "message": "config not found"})
		return
	}

	err = configDao.update(config)
	if err != nil {
		context.JSON(200, gin.H{"code": 500, "message": err.Error()})
		return
	}

	configMap[config.Name] = &config

	context.JSON(200, gin.H{"code": 200, "message": "success"})

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

	selectConfig, err := configDao.getByName(config.Name)
	if err != nil {
		context.JSON(200, gin.H{"code": 500, "message": err.Error()})
		return
	}

	configMap[config.Name] = selectConfig

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

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		}

		c.Next()
	}
}

func Authorization() gin.HandlerFunc {
	return func(context *gin.Context) {
		requestSign := context.GetHeader("SIGN")
		if requestSign != sign {
			context.AbortWithStatusJSON(200, gin.H{"code": 500, "message": "sign error"})
			return
		}
		context.Next()
	}
}
