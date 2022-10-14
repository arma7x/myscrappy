package financialtimes

import (
  "net/http"
  "github.com/gin-gonic/gin"
  "strings"
  "github.com/PuerkitoBio/goquery"
)

func GetBondsAndRates(c *gin.Context) {
  url := "https://markets.ft.com/data/bonds"
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
    c.JSON(http.StatusOK, gin.H{
      "data": map[string]any{
        "interbankratesovernight": extractBondsAndRates(doc.Find("#interbankratesovernight-panel")),
        "officialinterestrates": extractBondsAndRates(doc.Find("#officialinterestrates-panel")),
        "marketrates": extractBondsAndRates(doc.Find("#marketrates-panel")),
      },
    })
  }
}

func extractBondsAndRates(panel *goquery.Selection) map[string]any {
  var headers []string
  var results []map[string]string
  table := panel.Find("table")
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
  return map[string]any{
    "header": headers,
    "data": results,
  }
}
