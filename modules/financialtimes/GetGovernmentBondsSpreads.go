package financialtimes

import (
  "net/http"
  "github.com/gin-gonic/gin"
  "strings"
  "github.com/PuerkitoBio/goquery"
)

func GetGovernmentBondsSpreads(c *gin.Context) {
  url := "https://markets.ft.com/data/bonds/government-bonds-spreads"
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
    var headers []string
    var results []map[string]string
    table := doc.Find("table")
    table.Children().Each(func(i int, tbody *goquery.Selection) {
      if (i == 0) {
        tbody.Children().Each(func(j int, tr *goquery.Selection) {
          tr.Children().Each(func(k int, td *goquery.Selection) {
            headers = append(headers, strings.TrimSpace(td.Text()))
          })
        })
      } else if (i == 1) {
        tbody.Children().Each(func(j int, tr *goquery.Selection) {
          t := make(map[string]string)
          tr.Children().Each(func(k int, td *goquery.Selection) {
            t[headers[k]] = strings.TrimSpace(td.Text())
          })
          results = append(results, t)
        })
      }
    })
    c.JSON(http.StatusOK, gin.H{
      "header": headers,
      "data": results,
    })
  }
}
