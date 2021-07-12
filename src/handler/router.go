package handler

import "github.com/gin-gonic/gin"

func AddRouter(e *gin.Engine) {
	e.GET("/getblocks", GetBlocks)
}
