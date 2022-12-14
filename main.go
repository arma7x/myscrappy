package main

import (
  "net/http"
  "runtime"
  "github.com/gin-contrib/cors"
  "github.com/gin-gonic/gin"
  "myscrappy/modules/publicinfobanjir"
  "myscrappy/modules/financialtimes"
  "myscrappy/modules/openmeteo"
)

func main() {
  r := gin.Default()
  r.Use(cors.Default())

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

  routeWeather := r.Group("/open-meteo/api/v1")
  {
    routeWeather.GET("/weather", openmeteo.GetWeather)
  }

  r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
