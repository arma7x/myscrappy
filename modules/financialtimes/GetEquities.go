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

func GetEquities(c *gin.Context) {
  rankingType := c.DefaultQuery("rankingType", "percentgainers")
  if _, exist := RANKING_TYPES[rankingType]; !exist {
    rankingType = "percentgainers"
  }
  rankingSet := c.DefaultQuery("rankingSet", "SP500")
  if _, exist := RANKING_SETS[rankingSet]; !exist {
    rankingSet = "SP500"
  }
  url := fmt.Sprintf("https://markets.ft.com/data/equities/ajax/updatemarketmovers?rankingType=%s&rankingSet=%s", rankingType, rankingSet)
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
      if _, exist := result["data"]; !exist {
        c.JSON(http.StatusInternalServerError, gin.H{
          "error": err.Error(),
          "url": url,
        })
      } else {
        data := result["data"].(map[string]interface{})
        if _, exist := data["content"]; !exist {
          c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
            "url": url,
          })
        } else {
          reader := strings.NewReader(data["content"].(string))
          if doc, err := goquery.NewDocumentFromReader(reader); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
              "error": err.Error(),
              "url": url,
            })
          } else {
            var headers []string
            var result []map[string]any
            table := doc.Find("table")
            table.Children().Each(func(i int, tbody *goquery.Selection) {
              if (i == 0) {
                tbody.Children().Each(func(j int, tr *goquery.Selection) {
                  tr.Children().Each(func(k int, th *goquery.Selection) {
                    if (k < 4) {
                      headers = append(headers, strings.TrimSpace(th.Text()))
                    }
                  })
                })
              } else if (i == 1) {
                tbody.Children().Each(func(j int, tr *goquery.Selection) {
                  temp := make(map[string]any)
                  tr.Children().Each(func(k int, th *goquery.Selection) {
                    if (k <= 1) {
                      full := strings.TrimSpace(th.Text())
                      code := strings.TrimSpace(th.Children().Last().Text())
                      full =  strings.TrimSpace(strings.Replace(full, code, "", 1))
                      temp[headers[k]] = []string{full, code}
                    } else if (k == 2) {
                      full := strings.TrimSpace(th.Text())
                      code := strings.TrimSpace(th.Find("span").Find("span").Text())
                      full =  strings.TrimSpace(strings.Replace(full, code, "", 1))
                      temp[headers[k]] = []string{full, code}
                    } else if (k == 3) {
                      temp[headers[k]] = strings.TrimSpace(th.Text())
                    }
                  })
                  result = append(result, temp)
                })
              }
            });
            c.JSON(http.StatusOK, gin.H{
              "header": headers,
              "data": result,
            })
          }
        }
      }
    }
  }
}
