package main

import "github.com/gin-gonic/gin"

import "github.com/OkanUysal/go-response"

func GetUsers(c *gin.Context) {
	response.Success(c, []string{})
}
