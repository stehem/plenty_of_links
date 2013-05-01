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
  "strings"
  "time"
)

type Link struct {
  url string
  subreddit string
  title string
}

type RedditJSON struct {
  Data struct {
    Children []struct {
      Data struct {
        Title string
        Url string
      }
    }
  }
}



func main() {
  //results := getlinks()
  //fmt.Println(results)
  //SaveLinks(GetTwitterLinks())
  subs := []string{"ruby", "golang", "python", "javascript", "clojure", "scala"}
  for _, sub := range subs {
    links := GetRedditLinks(sub)
    SaveLinks(links)
  }
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

func GetLinks() []Link {
  db := getdb()
  defer db.Close()
  rows, _ := db.Query("SELECT url, subreddit, title FROM links ORDER BY created_at LIMIT 50")
  results := []Link{}
  for rows.Next() {
    link := new(Link)
    rows.Scan(&link.url, &link.subreddit, &link.title)
    results = append(results, *link)
  }
  return results
}


func SaveLinks(links []Link) {
  db := getdb()
  defer db.Close()
  for _, link := range links {
    db.Exec("INSERT INTO links (url, subreddit, title, created_at) VALUES ($1, $2, $3, $4)", link.url, link.subreddit, link.title, time.Now())
  }
}


func GetRedditLinks(sub string) []Link {
  response, _ := http.Get("http://www.reddit.com/r/" + sub + ".json?limit=100")
  defer response.Body.Close()
  contents, _ := ioutil.ReadAll(response.Body)
  var stories RedditJSON
  json.Unmarshal(contents, &stories)
  children := stories.Data.Children
  links := []Link{}
  for _, child := range children {
    if IsNotReddit(child.Data.Url) {
      link := Link{url: child.Data.Url, subreddit: sub, title: child.Data.Title}
      links = append(links, link)
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

func IsNotReddit(url string) bool {
  if strings.Contains(url, "reddit") {
    return false
  }
  return true
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
