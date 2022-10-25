package openmeteo

import (
  "net/http"
  "strings"
  "github.com/gin-gonic/gin"
  "fmt"
  "io"
  "encoding/json"
)

func GetWeather(c *gin.Context) {
  latitude := c.DefaultQuery("latitude", "37.6")
  longitude := c.DefaultQuery("longitude", "-95.665")
  url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%s&longitude=%s&timezone=auto&&daily=weathercode", latitude, longitude)
  res, err := http.Get(url)
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{
      "error": err.Error(),
      "url": url,
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
      daily := result["daily"].(map[string]interface{})
      weathercode := daily["weathercode"].([]interface{})
      for k, v := range weathercode {
        weathercode[k] = ParseWMOCode(fmt.Sprintf("%.0f", v.(float64)))
      }
      daily["weathercode"] = weathercode
      result["daily"] = daily
      c.JSON(http.StatusOK, gin.H{
        "data": result,
      })
    }
  }
}

func ParseWMOCode(code string) string {
  switch code {
    case "0":
      return "Clear sky"
    case "1", "2", "3":
      kv := map[string]string{"1": "mainly clear", "2": "partly cloudy", "3": "overcast"}
      return kv[code]
    case "45", "48":
      kv := map[string]string{"45": "fog", "48": "depositing rime fog"}
      return kv[code]
    case "51", "53", "55":
      kv := map[string]string{"51": "light", "53": "moderate", "55": "dense intensity"}
      return strings.Join([]string{"Drizzle:", kv[code]}, " ")
    case "56", "57":
      kv := map[string]string{"56": "light", "57": "dense intensity"}
      return strings.Join([]string{"Freezing Drizzle:", kv[code]}, " ")
    case "61", "63", "65":
      kv := map[string]string{"61": "slight", "63": "moderate", "65": "heavy intensity"}
      return strings.Join([]string{"Rain:", kv[code]}, " ")
    case "66", "67":
      kv := map[string]string{"66": "light", "67": "heavy intensity"}
      return strings.Join([]string{"Freezing Rain:", kv[code]}, " ")
    case "71", "73", "75":
      kv := map[string]string{"71": "slight", "73": "moderate", "75": "heavy intensity"}
      return strings.Join([]string{"Snow Fall:", kv[code]}, " ")
    case "77":
      return "Snow grains"
    case "80", "81", "82":
      kv := map[string]string{"80": "slight", "81": "moderate", "82": "violent"}
      return strings.Join([]string{"Rain Showers:", kv[code]}, " ")
    case "85", "86":
      kv := map[string]string{"85": "slight", "86": "heavy"}
      return strings.Join([]string{"Snow Showers:", kv[code]}, " ")
    case "95":
      return "Thunderstorm: Slight or moderate"
    case "96", "99":
      kv := map[string]string{"96": "slight", "99": "heavy hail"}
      return strings.Join([]string{"hunderstorm:", kv[code]}, " ")
    default:
      return "Unknown"
  }
}
