package utils

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ParamUInt(c *gin.Context, name string) (uint, bool) {
	raw := c.Param(name)
	val, err := strconv.ParseUint(raw, 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": fmt.Sprintf("Invalid param %s", name)})
		c.Abort()
		return 0, false
	}
	return uint(val), true
}

func ParamInt(c *gin.Context, name string) (int, bool) {
	raw := c.Param(name)
	val, err := strconv.Atoi(raw)
	if err != nil {
		c.JSON(400, gin.H{"error": fmt.Sprintf("Invalid param %s", name)})
		c.Abort()
		return 0, false
	}
	return val, true
}
