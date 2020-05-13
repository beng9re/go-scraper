package main

import (
	"os"
	"strings"

	"github.com/byoungheekim/goScraper/scrapper"

	"github.com/labstack/echo"
)

const filename string = "jobs.csv"

func handleHome(c echo.Context) error {
	return c.File("home.html")
}

func handleScrapper(c echo.Context) error {
	defer os.Remove(filename)
	//c.FormParams("want")
	term := strings.ToLower(scrapper.CleanString(c.FormValue("want")))
	scrapper.Scrape(term)

	return c.Attachment("jobs.csv", "jobs.csv")
}

func main() {
	e := echo.New()
	e.GET("/", handleHome)
	e.POST("/scrape", handleScrapper)
	e.Logger.Fatal(e.Start(":1323"))
}
