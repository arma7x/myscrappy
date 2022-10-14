package publicinfobanjir

import (
  "golang.org/x/net/html"
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
