package main

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

var configMap = make(map[string]*Config)

var sign = "0D000721"

func main() {
	var err error

	err = initTask()
	if err != nil {
		println(err.Error())
		return
	}

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
	r.GET("/config/:id/history", getConfigHistory)
	r.POST("/cbs_user/login", cbsUserLogin)
	r.POST("/cbs_user/logout", cbsUserLogout)
	r.GET("/cbs_user/info", cbsUserInfo)
	err = r.Run(":19200")
	if err != nil {
		println(err.Error())
		return
	}
}

func getConfigHistory(context *gin.Context) {
	id, err := strconv.Atoi(context.Param("id"))
	if err != nil {
		context.JSON(200, gin.H{"code": 500, "message": "id is required"})
		return
	}
	pageSize, err := strconv.Atoi(context.Query("pageSize"))
	if err != nil {
		context.JSON(200, gin.H{"code": 500, "message": "pageSize is required"})
		return
	}

	pageIndex, err := strconv.Atoi(context.Query("pageIndex"))
	if err != nil {
		context.JSON(200, gin.H{"code": 500, "message": "pageIndex is required"})
		return
	}

	list, err := configHistoryDao.list(id, pageSize, pageIndex)
	if err != nil {
		context.JSON(200, gin.H{"code": 500, "message": err.Error()})
		return
	}

	count, err := configHistoryDao.getCount(id)
	if err != nil {
		context.JSON(200, gin.H{"code": 500, "message": err.Error()})
		return
	}

	context.JSON(200, gin.H{"code": 200, "data": gin.H{
		"list":  list,
		"count": count,
	}})
	return
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
	pageSize, err := strconv.Atoi(context.Query("pageSize"))

	if err != nil {
		context.JSON(200, gin.H{"code": 500, "message": "pageSize is required"})
		return
	}

	pageIndex, err := strconv.Atoi(context.Query("pageIndex"))

	if err != nil {
		context.JSON(200, gin.H{"code": 500, "message": "pageIndex is required"})
		return
	}

	list, err := configDao.list(keyword, pageSize, pageIndex)
	if err != nil {
		context.JSON(200, gin.H{"code": 500, "message": err.Error()})
		return
	}

	count, err := configDao.getCount(keyword)
	if err != nil {
		context.JSON(200, gin.H{"code": 500, "message": err.Error()})
		return
	}

	context.JSON(200, gin.H{"code": 200, "data": gin.H{
		"list":  list,
		"count": count,
	}})
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

	oldConfig := configMap[config.Name]
	if oldConfig == nil {
		context.JSON(200, gin.H{"code": 500, "message": "config not found"})
		return
	}

	userInfo := cbsTokenUserInfoMap[context.GetHeader("Authorization")]
	err = configDao.update(config, userInfo.Nickname)
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

	userInfo := cbsTokenUserInfoMap[context.GetHeader("Authorization")]
	err = configDao.create(config, userInfo.Nickname)
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
	configList, err := configDao.list("", 1000000, 1)
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
