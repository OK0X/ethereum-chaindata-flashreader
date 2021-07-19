package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetBlocks(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")
	if from == "" || to == "" {
		c.Abort()
		c.JSON(http.StatusNotFound, "lack parameter from or to")
		return
	}
	ifrom, _ := strconv.ParseInt(from, 10, 64)
	ito, _ := strconv.ParseInt(to, 10, 64)
	blocks, _ := FlashRead.IndexTransactions(uint64(ifrom), uint64(ito), false, nil)
	c.JSON(http.StatusOK, blocks)
}
