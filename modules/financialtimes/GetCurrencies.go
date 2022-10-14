package financialtimes

import (
  "net/http"
  "github.com/gin-gonic/gin"
  "io"
  "fmt"
  "strings"
  "github.com/PuerkitoBio/goquery"
  "encoding/json"
)

func GetCurrencies(c *gin.Context) {
  group := c.DefaultQuery("group", "Majors")
  if _, exist := CURRENCY[group]; !exist {
    group = "Majors"
  }
  url := fmt.Sprintf("https://markets.ft.com/data/currencies/ajax/crossratesforselectedregion?group=%s", group)
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
      c.JSON(http.StatusInternalServerError, gin.H{
        "error": err.Error(),
        "url": url,
      })
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
