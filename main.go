package main

import (
  "net/http"
  "runtime"
  "github.com/gin-gonic/gin"
  "myscrappy/modules/publicinfobanjir"


  "io"
  "fmt"
  "strings"
  "github.com/PuerkitoBio/goquery"
  "encoding/json"
  // "golang.org/x/net/html"
  // "regexp"
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
    routeFinancialTimes.GET("/currencies", GetCurrencies)
  }

  r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

var CURRENCY = map[string]string{
  "Majors": "Majors Currency",
  "Europe": "Europe Currency",
  "Americas": "Americas Currency",
  "Africa": "Africa Currency",
  "AsiaPacific": "Asia Pacific Currency",
}

var RANKING_SETS = map[string]string{
  "AustralianStockExchange": "Australia",
  "BrusselsStockExchange": "Belgium",
  "SaoPauloStockExchange": "Brazil",
  "TorrontoStockExchange": "Canada",
  "PragueStockExchange": "Czech Republic",
  "CopenhagenStockExchange": "Denmark",
  "HelsinkiStockExchange": "Finland",
  "ParisStockExchange": "France",
  "FrankfurtStockExchange": "Germany",
  "AthensStockExchange": "Greece",
  "StockExchangeOfHongKongLimited": "Hong Kong",
  "BudapestStockExchange": "Hungary",
  "BombayStockExchange": "India",
  "JakartaStockExchange": "Indonesia",
  "TelAvivStockExchange": "Israel",
  "TokyoStockExchange": "Japan",
  "MexicanStockExchange": "Mexico",
  "EuronextAmsterdam": "Netherlands",
  "NewZealandStockExchange": "New Zealand",
  "OsloStockExchange": "Norway",
  "LisbonStockExchange": "Portugal",
  "SingaporeExchangeSecuritiesTrading": "Singapore",
  "JohannesburgStockExchange": "South Africa",
  "MercadoContinuoEspanol": "Spain",
  "StockholmStockExchange": "Sweden",
  "SwissExchange": "Switzerland",
  "LondonStockExchange": "United Kingdom",
  "SP500": "United States",
}

var RANKING_TYPES = map[string]string{
  "percentgainers": "Gainers",
  "percentlosers": "Losers",
  "highestvolume": "Movers",
}

func GetCurrencies(c *gin.Context) {
  // https://markets.ft.com/data/currencies/ajax/crossratesforselectedregion?group=Majors
  group := c.DefaultQuery("group", "Majors")
  if _, exist := CURRENCY[group]; !exist {
    group = "Majors"
  }
  url := fmt.Sprintf("https://markets.ft.com/data/currencies/ajax/crossratesforselectedregion?group=%s", group)
  fmt.Println(url)
  res, err := http.Get(url)
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{
      "error": err.Error(),
      "url": url,
    })
    return
  }

  defer res.Body.Close()
  if res.StatusCode != 200 {
    c.JSON(res.StatusCode, gin.H{
      "error": res.StatusCode,
      "url": res.Status,
    })
    return
  }

  b, err := io.ReadAll(res.Body)
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{
      "error": err.Error(),
      "url": url,
    })
  } else {
    var result map[string]interface{}
    if err := json.Unmarshal(b, &result); err != nil {
      fmt.Println(err)
    } else {
      if _, exist := result["html"]; !exist {
        c.JSON(http.StatusInternalServerError, gin.H{
          "error": err.Error(),
          "url": url,
        })
      } else {
        reader := strings.NewReader(result["html"].(string))
        if doc, err := goquery.NewDocumentFromReader(reader); err != nil {
          c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
            "url": url,
          })
          return
        } else {
          var headers []string
          var result [][][]string
          table := doc.Find("table")
          table.Children().Each(func(i int, tbody *goquery.Selection) {
            if (i == 0) {
              tbody.Children().Each(func(j int, tr *goquery.Selection) {
                tr.Children().Each(func(k int, th *goquery.Selection) {
                  if (k == 0) {
                    headers = append(headers, strings.TrimSpace(th.Children().Last().Text()))
                  } else {
                    t := fmt.Sprintf("%s/%s", strings.TrimSpace(th.Children().First().Text()), strings.TrimSpace(th.Children().Last().Text()))
                    headers = append(headers, t)
                  }
                })
              })
            } else if (i == 1) {
              tbody.Children().Each(func(j int, tr *goquery.Selection) {
                var trd [][]string
                tr.Children().Each(func(k int, th *goquery.Selection) {
                  trd = append(trd, []string{headers[k], strings.TrimSpace(th.Text())})
                })
                result = append(result, trd)
              })
            }
          });
          c.JSON(http.StatusOK, gin.H{
            "data": result,
          })
        }
      }
    }
  }
}

func GetCommodities(c *gin.Context) {
  // https://markets.ft.com/data/commodities
}

func GetBondsAndRates(c *gin.Context) {
  // https://markets.ft.com/data/bonds
}

func GetGovernmentBondsSpreads(c *gin.Context) {
  // https://markets.ft.com/data/bonds/government-bonds-spreads
}

func GetEquities(c *gin.Context) {
  // https://markets.ft.com/data/equities/ajax/updatemarketmovers
}
