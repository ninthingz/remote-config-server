package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron/v2"
	"io"
	"log"
	"net/http"
	"time"
)

type CBSCommonResult[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message" json:"message,omitempty"`
	Data    T      `json:"data"`
}

type CBSUserLoginDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Source   string `json:"source"`
}

type CBSUserToken struct {
	Token     string `json:"token"`
	TokenHead string `json:"tokenHead"`
}

type CBSUserInfo struct {
	Username    string    `json:"username"`
	Nickname    string    `json:"nickname"`
	Email       string    `json:"email"`
	ExpiredTime time.Time `json:"expiredTime" json:"expiredTime,omitempty"`
}

var cbsServiceHost = "http://192.168.12.245:8032"

var cbsTokenUserInfoMap = make(map[string]*CBSUserInfo)

var s gocron.Scheduler

func initTask() error {

	var err error
	s, err = gocron.NewScheduler()
	if err != nil {
		return err
	}

	j, err := s.NewJob(
		gocron.DurationJob(
			1*time.Hour,
		),
		gocron.NewTask(
			func() {
				nowTime := time.Now()
				deleteKeys := make([]string, 0)
				for key, value := range cbsTokenUserInfoMap {
					if nowTime.After(value.ExpiredTime) {
						deleteKeys = append(deleteKeys, key)
					}
				}

				if len(deleteKeys) > 0 {
					for _, key := range deleteKeys {
						delete(cbsTokenUserInfoMap, key)
					}
				}
			},
		),
	)
	if err != nil {
		// handle error
		return err
	}
	// each job has a unique id
	fmt.Println(j.ID())

	// start the scheduler
	s.Start()

	return nil
}

func cbsUserLogin(context *gin.Context) {

	var cbsUserLoginDTO CBSUserLoginDTO
	err := context.ShouldBindJSON(&cbsUserLoginDTO)
	if err != nil {
		context.JSON(200, gin.H{"code": 500, "message": err.Error()})
		return
	}

	var cbsUserToken CBSCommonResult[CBSUserToken]
	var resp *http.Response
	var body []byte
	var errResp error

	body, err = json.Marshal(cbsUserLoginDTO)

	if err != nil {
		context.JSON(200, gin.H{"code": 500, "message": err.Error()})
		return
	}

	resp, err = http.Post(cbsServiceHost+"/api/v1/user/login", "application/json", bytes.NewBuffer(body))

	if err != nil {
		context.JSON(200, gin.H{"code": 500, "message": err.Error()})
		return
	}

	if resp.StatusCode != 200 {
		errResp = errors.New(fmt.Sprintf("request cbs user login failed, status code: %d", resp.StatusCode))
		context.JSON(200, gin.H{"code": 500, "message": errResp.Error()})
		return
	} else {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Println(err.Error())
			}
		}(resp.Body)
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			context.JSON(200, gin.H{"code": 500, "message": err.Error()})
			return
		}
		err = json.Unmarshal(body, &cbsUserToken)
		if err != nil {
			context.JSON(200, gin.H{"code": 500, "message": err.Error()})
			return
		}

		if cbsUserToken.Code != 200 {
			context.JSON(200, gin.H{"code": 500, "message": cbsUserToken.Message})
			return
		}

		context.JSON(200, gin.H{"code": 200, "data": cbsUserToken.Data})
	}

}

func cbsUserLogout(context *gin.Context) {

	token := context.GetHeader("Authorization")
	if token == "" {
		context.JSON(200, gin.H{"code": 500, "message": "no token"})
		return
	}

	if cbsTokenUserInfoMap[token] != nil {
		delete(cbsTokenUserInfoMap, token)
	}

	context.JSON(200, gin.H{"code": 200, "message": "success"})
}

func cbsUserInfo(context *gin.Context) {
	token := context.GetHeader("Authorization")
	if token == "" {
		context.JSON(200, gin.H{"code": 500, "message": "no token"})
		return
	}

	if cbsTokenUserInfoMap[token] != nil {
		userInfo := cbsTokenUserInfoMap[token]
		if time.Now().Before(userInfo.ExpiredTime) {
			userInfo.ExpiredTime = time.Now().Add(24 * 7 * time.Hour)
			cbsTokenUserInfoMap[token] = userInfo

			context.JSON(200, gin.H{"code": 200, "data": *userInfo})
			return
		}
	}

	req, err := http.NewRequest("GET", cbsServiceHost+"/api/v1/user/info", nil)

	if err != nil {
		context.JSON(200, gin.H{"code": 500, "message": err.Error()})
		return
	}
	req.Header.Set("Authorization", token)
	resp, err := http.DefaultClient.Do(req)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)

	if err != nil {
		context.JSON(200, gin.H{"code": 500, "message": err.Error()})
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		context.JSON(200, gin.H{"code": 500, "message": err.Error()})
		return
	}
	var cbsUserInfo CBSCommonResult[CBSUserInfo]
	err = json.Unmarshal(body, &cbsUserInfo)
	if err != nil {
		context.JSON(200, gin.H{"code": 500, "message": err.Error()})
		return
	}

	if cbsUserInfo.Code != 200 {
		context.JSON(200, gin.H{"code": 500, "message": cbsUserInfo.Message})
		return
	}

	cbsUserInfo.Data.ExpiredTime = time.Now().Add(24 * 7 * time.Hour)
	cbsTokenUserInfoMap[token] = &cbsUserInfo.Data

	context.JSON(200, gin.H{"code": 200, "data": cbsUserInfo.Data})

}

func Authorization() gin.HandlerFunc {
	return func(context *gin.Context) {
		requestSign := context.GetHeader("SIGN")
		if requestSign == sign {
			context.Next()
			return
		}

		if context.Request.URL.Path == "/api/v1/cbs_user/login" ||
			context.Request.URL.Path == "/api/v1/cbs_user/info" ||
			context.Request.URL.Path == "/api/v1/config/element" ||
			(context.Request.URL.Path == "/api/v1/config" && context.Request.Method == "GET") {
			context.Next()
			return
		}
		token := context.GetHeader("Authorization")
		if cbsTokenUserInfoMap[token] == nil {
			context.AbortWithStatusJSON(200, gin.H{"code": 500, "message": "no token"})
			return
		} else {
			nowTime := time.Now()
			if nowTime.After(cbsTokenUserInfoMap[token].ExpiredTime) {
				context.AbortWithStatusJSON(200, gin.H{"code": 500, "message": "token expired"})
				return
			}
		}
		context.Next()
	}
}
