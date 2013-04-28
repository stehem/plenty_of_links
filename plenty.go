// go get github.com/bmizerany/pq
// go get github.com/dlintw/goconf
package main

import (
  _ "github.com/bmizerany/pq"
  "github.com/dlintw/goconf"
  "database/sql"
  "fmt"
)

type Link struct {
  link string
  source string
  category string
}

func main() {
  results := getlinks()
  fmt.Println(results)
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
  rows, _ := db.Query("SELECT link, source, category FROM links")
  results := []Link{}
  for rows.Next() {
    link := new(Link)
    rows.Scan(&link.link, &link.source, &link.category)
    results = append(results, *link)
  }
  return results
}


