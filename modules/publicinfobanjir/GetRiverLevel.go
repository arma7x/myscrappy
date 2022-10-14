package publicinfobanjir

import (
  "fmt"
  "net/http"
  "strings"
  "github.com/gin-gonic/gin"
  "github.com/PuerkitoBio/goquery"
)

var riverHeaders = []string{"No","Station ID","Station Name","District","Main Basin","Sub River Basin","Last Updated","Water Level","Threshold"}
var riverLevels = []string{"Normal","Alert","Warning","Danger"}

func GetRiverLevel(c *gin.Context) {
  state := strings.ToUpper(c.DefaultQuery("state", "SEL"))
  if _, exist := STATE[state]; !exist {
    state = "SEL"
  }
  html := c.DefaultQuery("html", "0")
  url := fmt.Sprintf("http://publicinfobanjir.water.gov.my/aras-air/data-paras-air/aras-air-data/?state=%s&district=ALL&station=ALL&lang=en", state)

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

  if doc, err := goquery.NewDocumentFromReader(res.Body); err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{
      "error": err.Error(),
      "url": url,
    })
  } else {
    if (html == "1") {
      doc.Find("script").Each(func(i int, s *goquery.Selection) {
        RemoveNode(doc.Get(0), s.Get(0))
      })
      doc.Find("link").Each(func(i int, s *goquery.Selection) {
        RemoveNode(doc.Get(0), s.Get(0))
      })
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
            t := make(map[string]string)
            tr.Children().Each(func(k int, td *goquery.Selection) {
              if k < 8 {
                t[riverHeaders[k]] = strings.TrimSpace(td.Text())
              } else {
                t[riverLevels[k - 8]] = strings.TrimSpace(td.Text())
              }
            })
            results = append(results, t)
          })
        }
      })
      c.JSON(http.StatusOK, gin.H{
        "data": results,
        "river_header": riverHeaders,
        "level_header": riverLevels,
      })
    }
  }
}
