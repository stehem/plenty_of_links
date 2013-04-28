// go get github.com/bmizerany/pq
// go get github.com/dlintw/goconf
package main

import (
  _ "github.com/bmizerany/pq"
  "github.com/dlintw/goconf"
  "database/sql"
  "fmt"
  "net/http"
  "io/ioutil"
  "encoding/json"
  "unicode/utf8"
)

type Link struct {
  url string
  source string
  category string
  text string
}

type TwitterJSON struct {
  Results []Result
}

type Result struct {
  Text string
  Entities Entity
}

type Entity struct {
  Urls []Url
}

type Url struct {
  Expanded_url string
}



func main() {
  //results := getlinks()
  //fmt.Println(results)
  SaveLinks(GetTwitterLinks())
}

func getdb() *sql.DB {
  c, _ := goconf.ReadConfigFile("db.config")
  user, _ := c.GetString("dev", "user") 
  pass, _ := c.GetString("dev", "password") 
  base := "user=%s password=%s dbname=postgres"
  settings := fmt.Sprintf(base, user, pass)
  db, _ := sql.Open("postgres", settings)
  return db
}

func getlinks() []Link {
  db := getdb()
  defer db.Close()
  rows, _ := db.Query("SELECT url, source, category, text FROM links")
  results := []Link{}
  for rows.Next() {
    link := new(Link)
    rows.Scan(&link.url, &link.source, &link.category, &link.text)
    results = append(results, *link)
  }
  return results
}


func SaveLinks(links []Link) {
  db := getdb()
  defer db.Close()
  for _, link := range links {
    db.Exec("INSERT INTO links (url, source, category, text) VALUES ($1, $2, $3, $4)", link.url, link.source, link.category, link.text)
  }
}

func GetTwitterLinks() []Link {
  response, _ := http.Get("http://search.twitter.com/search.json?q=golang&include_entities=true&rpp=300&lang=en")
  defer response.Body.Close()
  contents, _ := ioutil.ReadAll(response.Body)
  var resp TwitterJSON
  json.Unmarshal(contents, &resp)
  var links []Link
  for _, result := range resp.Results {
    if len(result.Entities.Urls) > 0 {
      link := Link{url: result.Entities.Urls[0].Expanded_url, source: "Twitter", category: "Go", text: result.Text}
      if GoodUrl(links, link.url) {
        links = append(links, link)
      }
    }
  }
  return links
}

func Contains(seq[]Link, url string) bool{
  for _, link := range seq {
    if link.url == url {
      return true
    }
  }
  return false
}

func NoShortUrl(url string) bool {
  count := utf8.RuneCountInString(url)
  if count > 25 {
    return true
  }
  return false
}

func GoodUrl(links []Link, url string) bool {
  res := !Contains(links, url) && NoShortUrl(url)
  return res
}
