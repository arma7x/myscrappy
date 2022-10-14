package main

import (
  "net/http"
  "runtime"
  "github.com/gin-gonic/gin"
  "myscrappy/modules/publicinfobanjir"
  "myscrappy/modules/financialtimes"
)

func main() {
  r := gin.Default()
  r.GET("/", func(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
      "go_version": runtime.Version(),
    })
  })
  routePublicInfoBanjir := r.Group("/publicinfobanjir/api/v1")
  {
    routePublicInfoBanjir.GET("/state", publicinfobanjir.GetStateList)
    routePublicInfoBanjir.GET("/river", publicinfobanjir.GetRiverLevel)
    routePublicInfoBanjir.GET("/rain", publicinfobanjir.GetRainLevel)
  }
  routeFinancialTimes := r.Group("/ft/api/v1")
  {
    routeFinancialTimes.GET("/currencies", financialtimes.GetCurrencies)
    routeFinancialTimes.GET("/commodities", financialtimes.GetCommodities)
    routeFinancialTimes.GET("/bondsandrates", financialtimes.GetBondsAndRates)
    routeFinancialTimes.GET("/governmentbondsspreads", financialtimes.GetGovernmentBondsSpreads)
    routeFinancialTimes.GET("/equities", financialtimes.GetEquities)
  }

  r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
