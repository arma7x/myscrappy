package publicinfobanjir

import (
  "fmt"
  "net/http"
  "strings"
  "github.com/gin-gonic/gin"
  "github.com/PuerkitoBio/goquery"
  "golang.org/x/net/html"
  "regexp"
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

var RIVER_HEADER = []string{"No","Station ID","Station Name","District","Main Basin","Sub River Basin","Last Updated","Water Level","Threshold"}
var RIVER_LEVEL = []string{"Normal","Alert","Warning","Danger"}

// https://github.com/PuerkitoBio/goquery/issues/17
func RemoveNode(root_node *html.Node, remove_me *html.Node) {
  found_node := false
  check_nodes := make(map[int]*html.Node)
  i := 0
  for n := root_node.FirstChild; n != nil; n = n.NextSibling {
    if n == remove_me {
        found_node = true
        n.Parent.RemoveChild(n)
    }
    check_nodes[i] = n
    i++
  }
  if found_node == false {
    for _, item := range check_nodes {
      RemoveNode(item, remove_me)
    }
  }
}

func GetStateList(c *gin.Context) {
  c.JSON(http.StatusOK, gin.H{
    "state": STATE,
  })
}

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
  }
}

func GetRainLevel(c *gin.Context) {
  // https://publicinfobanjir.water.gov.my/hujan/data-hujan/?state=KEL&lang=en
  state := strings.ToUpper(c.DefaultQuery("state", "SEL"))
  if _, exist := STATE[state]; !exist {
    state = "SEL"
  }
  html := c.DefaultQuery("html", "0")
  url := fmt.Sprintf("http://publicinfobanjir.water.gov.my/wp-content/themes/shapely/agency/searchresultrainfall.php?state=%s&district=ALL&station=ALL&language=1&loginStatus=0", state)
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
      calculatejs := "<script>function calculate(){var e=new URL(document.location.toString());e.searchParams.set('html',0),fetch(e.toString()).then(e=>e.json()).then(e=>{console.clear();let t={};e.data.forEach(e=>{var a=0;if(null!=e['Daily Rainfall']&&Object.keys(e['Daily Rainfall']).length>0){for(let l in e['Daily Rainfall']){let n=parseFloat(e['Daily Rainfall'][l]);n>=0&&(a+=n)}null==t[e.District]&&(t[e.District]=0),t[e.District]+=a}});var a=[];for(var l in t)a.push({name:l,value:t[l]});a.sort((e,t)=>e.value>t.value?-1:1);let n=new Date,i=n.getDate(),r=n.getMonth()+1;n.setTime(n.getTime()-5184e5);let o=n.getDate(),c=n.getMonth()+1;var h=document.createElement('ul');h.setAttribute('id','total_rainfall');var m=`Total rainfall for 7 consecutive days(${i}/${r} - ${o}/${c}):`,s=document.createElement('h3');s.setAttribute('style','margin-left:4px;'),document.body.appendChild(s),s.innerHTML=m,a.forEach(e=>{m+=`${e.name}${'-'.repeat(30-e.name.length)}-> ${e.value.toFixed(2)}mm`;var t=document.createElement('li');h.appendChild(t),t.innerHTML=`${e.name} ${e.value.toFixed(2)}mm`}),document.body.appendChild(h),console.log(m)}).catch(e=>{console.error(e)})}calculate();</script>"
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
        textHtml = strings.Replace(textHtml, "</body>", fmt.Sprintf("%s</body>", calculatejs), 1)
        c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(textHtml))
      }
    } else {
      var results []map[string]any
      var rainHeader []string
      var dailyRailfall [7]string
      table := doc.Find("table")
      table.Children().Each(func(i int, tbody *goquery.Selection) {
        if (i == 0) {
          tbody.Children().Each(func(j int, tr *goquery.Selection) {
            if (j > 0) {
              tr.Children().Each(func(k int, th *goquery.Selection) {
                if (j == 1) {
                  if (k == 6) {
                    re := regexp.MustCompile(`\d{2}/\d{2}/\d{4}`)
                    dailyRailfall[6] = re.FindString(strings.TrimSpace(th.Text()))
                  } else {
                    rainHeader = append(rainHeader, strings.TrimSpace(th.Text()))
                  }
                } else if (j == 2) {
                  dailyRailfall[k] = strings.TrimSpace(th.Text())
                }
              })
            }
          })
        } else if (i == 1) {
          tbody.Children().Each(func(j int, tr *goquery.Selection) {
            t := make(map[string]any)
            daily := make(map[string]any)
            tr.Children().Each(func(k int, td *goquery.Selection) {
              if (k >=0 && k <= 4) {
                t[rainHeader[k]] = strings.TrimSpace(td.Text())
              } else if (k >=5 && k <= 11) {
                daily[dailyRailfall[k - 5]] = strings.TrimSpace(td.Text())
                if (k == 5) {
                  t[rainHeader[5]] = daily
                }
              } else {
                t[rainHeader[len(rainHeader) - 1]] = strings.TrimSpace(td.Text())
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
  }
}
