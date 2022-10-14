package financialtimes

import (
  "net/http"
  "github.com/gin-gonic/gin"
  "strings"
  "github.com/PuerkitoBio/goquery"
)

func GetCommodities(c *gin.Context) {
  url := "https://markets.ft.com/data/commodities"
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
        "energy": extractCommodities(doc.Find("#energy-panel")),
        "metals": extractCommodities(doc.Find("#metals-panel")),
        "agricultureandlumber": extractCommodities(doc.Find("#agricultureandlumber-panel")),
      },
    })
  }
}

func extractCommodities(panel *goquery.Selection) map[string]any {
  var headers []string
  var results []map[string]any
  table := panel.Find("table")
  table.Children().Each(func(i int, tbody *goquery.Selection) {
    if (i == 0) {
      tbody.Children().Each(func(j int, tr *goquery.Selection) {
        tr.Children().Each(func(k int, td *goquery.Selection) {
          if (k == 4) {
            var spans []string
            td.Children().First().Children().Each(func(j int, span *goquery.Selection) {
              spans = append(spans, strings.TrimSpace(span.Text()))
            })
            headers = append(headers, strings.Join(spans, "@"))
          } else {
            headers = append(headers, strings.TrimSpace(td.Text()))
          }
        })
      })
    } else if (i == 1) {
      tbody.Children().Each(func(j int, tr *goquery.Selection) {
        t := make(map[string]any)
        tr.Children().Each(func(k int, td *goquery.Selection) {
          if (k == 0) {
            t[headers[k]] = []string{strings.TrimSpace(td.Children().First().Text()), strings.TrimSpace(td.Children().Last().Text())}
          } else if (k == 1) {
            full := strings.TrimSpace(td.Text())
            code := strings.TrimSpace(td.Children().Last().Text())
            full =  strings.TrimSpace(strings.Replace(full, code, "", 1))
            t[headers[k]] = []string{full, code}
          } else if (k == 2) {
            full := strings.TrimSpace(td.Text())
            code := strings.TrimSpace(td.Find("span").Find("span").Text())
            full =  strings.TrimSpace(strings.Replace(full, code, "", 1))
            t[headers[k]] = []string{full, code}
          } else if (k == 3) {
            if val, exist := td.Children().Last().Attr("style"); exist {
              t[headers[k]] = strings.TrimSpace(val)
            } else {
              t[headers[k]] = false
            }
          } else if (k == 4) {
            hl := make(map[string][]string)
            keys := map[string]string{
              "low": ".mod-ui-range-bar__container__label--lo",
              "high": ".mod-ui-range-bar__container__label--hi",
            }
            for key := range keys {
              dom := td.Find(keys[key])
              if (strings.TrimSpace(dom.Text()) != "") {
                var s []string
                dom.Children().Each(func(j int, span *goquery.Selection) {
                  s = append(s, strings.TrimSpace(span.Text()));
                })
                hl[key] = s
              } else {
                hl[key] = []string{}
              }
            }
            t[headers[k]] = hl
          }
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
