package main

import (
  "net/http"
  "runtime"
  "github.com/gin-gonic/gin"
  "myscrappy/modules/publicinfobanjir"
)

func main() {
  r := gin.Default()
  r.GET("/", func(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
      "go_version": runtime.Version(),
    })
  })
  routepublicinfobanjir := r.Group("/publicinfobanjir/v1")
  {
    routepublicinfobanjir.GET("/state", publicinfobanjir.GetStateList)
    routepublicinfobanjir.GET("/river", publicinfobanjir.GetRiverLevel)
    routepublicinfobanjir.GET("/rain", publicinfobanjir.GetRainLevel)
  }
  r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
