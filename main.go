package main

import (
  "fmt"
  "net/http"
  "runtime"
  "strings"
  "github.com/gin-gonic/gin"
  "github.com/PuerkitoBio/goquery"
)

var STATE = map[string]string{
  "KDH": "Kedah",
  "PNG": "Pulau Pinang",
  "PRK": "Perak",
  "SEL": "Selangor",
  "WLH": "Wilayah Persekutuan Kuala Lumpur",
  "PTJ": "Wilayah Persekutuan Putrajaya",
  "NSN": "Negeri Sembilan",
  "MLK": "Melaka",
  "JHR": "Johor",
  "PHG": "Pahang",
  "TRG": "Terengganu",
  "KEL": "Kelantan",
  "SRK": "Sarawak",
  "SAB": "Sabah",
  "WLP": "Wilayah Persekutuan Labuan",
}

var RIVER_HEADER = []string{"No","StationID","StationName","District","MainBasin","SubRiverBasin","LastUpdated","WaterLevel","Threshold"}
var RIVER_LEVEL = []string{"Normal","Alert","Warning","Danger"}

func main() {
  r := gin.Default()
  r.GET("/", func(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
      "go_version": runtime.Version(),
    })
  })
  publicinfobanjir := r.Group("/publicinfobanjir/v1")
  {
    publicinfobanjir.GET("/state", getStateList)
    publicinfobanjir.GET("/river", getRiverLevel)
    publicinfobanjir.GET("/rain", getRainLevel)
  }
  r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func getStateList(c *gin.Context) {
  c.JSON(http.StatusOK, gin.H{
    "state": STATE,
  })
}

func getRiverLevel(c *gin.Context) {
  state := strings.ToUpper(c.DefaultQuery("state", "SEL"))
  if _, exist := STATE[state]; !exist {
    state = "SEL"
  }
  html := c.DefaultQuery("html", "0")
  url := fmt.Sprintf("http://publicinfobanjir.water.gov.my/aras-air/data-paras-air/aras-air-data/?state=%s&district=ALL&station=ALL&lang=en", state)
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
    c.JSON(http.StatusInternalServerError, gin.H{
      "error": res.StatusCode,
      "url": res.Status,
    })
    return
  }

  if doc, err := goquery.NewDocumentFromReader(res.Body); err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{
      "error": err.Error(),
      "url": url,
    })
    return
  } else {
    if (html == "1") {
      if textHtml, err := doc.Html(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
          "error": err.Error(),
          "url": url,
        })
      } else {
        c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(textHtml))
      }
    } else {
      var results []map[string]string
      table := doc.Find("table")
      table.Children().Each(func(i int, tbody *goquery.Selection) {
        if (i == 1) {
          tbody.Children().Each(func(j int, tr *goquery.Selection) {
            t := make(map[string]string);
            tr.Children().Each(func(k int, td *goquery.Selection) {
              if k < 8 {
                t[RIVER_HEADER[k]] = strings.TrimSpace(td.Text())
              } else {
                t[RIVER_LEVEL[k - 8]] = strings.TrimSpace(td.Text())
              }
            })
            results = append(results, t)
          })
        }
      })
      c.JSON(http.StatusOK, gin.H{
        "data": results,
      })
    }
    return
  }
}

func getRainLevel(c *gin.Context) {
  // https://publicinfobanjir.water.gov.my/hujan/data-hujan/?state=KEL&lang=en
  state := strings.ToUpper(c.DefaultQuery("state", "SEL"))
  if _, exist := STATE[state]; !exist {
    state = "SEL"
  }
  html := c.DefaultQuery("html", "0")
  url := fmt.Sprintf("http://publicinfobanjir.water.gov.my/hujan/data-hujan/?state=%s&lang=en", state)
  c.JSON(http.StatusOK, gin.H{
    "html": html,
    "state": state,
    "url": url,
  })
}
