package publicinfobanjir

import (
  "net/http"
  "github.com/gin-gonic/gin"
)

func GetStateList(c *gin.Context) {
  c.JSON(http.StatusOK, gin.H{
    "state": STATE,
  })
}
