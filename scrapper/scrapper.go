package scrapper

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type extractedJob struct {
	id       string
	title    string
	company  string
	location string
	salary   string
	summary  string
}

//Scrape Indeed by a term
func Scrape(term string) {
	var baseURL string = "https://kr.indeed.com/jobs?q=" + term
	var jobs []extractedJob
	totalPages := getPages(baseURL)
	c := make(chan []extractedJob)

	for i := 0; i < totalPages; i++ {
		go getPage(baseURL, i, c)
		//extractedJobs := <-c
	}

	for i := 0; i < totalPages; i++ {
		jobs = append(jobs, <-c...)
	}

	writeJobs(jobs)
	fmt.Println("done, extentd: ", len(jobs))

}

func getPage(baseURL string, page int, c chan<- []extractedJob) {
	fmt.Println(page)
	subChan := make(chan extractedJob)
	var jobs []extractedJob
	pageURL := baseURL + "&start=" + strconv.Itoa(page*10)
	fmt.Println("request URL : ", pageURL)
	res, err := http.Get(pageURL)
	checkErr(err)
	checkCode(res)
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)
	sercheCard := doc.Find(".jobsearch-SerpJobCard")
	sercheCard.Each(func(i int, card *goquery.Selection) {
		go extractJob(card, subChan)

	})

	for i := 0; i < sercheCard.Length(); i++ {
		jobs = append(jobs, <-subChan)
		//fmt.Println(jobs)
	}

	c <- jobs

	//return jobs
}

func extractJob(card *goquery.Selection, c chan<- extractedJob) {
	id, _ := card.Attr("data-jk")
	title := card.Find(".title>a").Text()
	location := card.Find(".location accessible-contrast-color-location").Text()
	company := card.Find(".company").Text()
	salary := card.Find(".salaryText").Text()
	summary := card.Find(".summary").Text()

	c <- extractedJob{
		id:       id,
		title:    CleanString(title),
		location: CleanString(location),
		company:  CleanString(company),
		summary:  CleanString(summary),
		salary:   CleanString(salary)}

}

func getPages(baseURL string) int {
	pages := 0
	res, err := http.Get(baseURL)

	checkErr(err)
	checkCode(res)
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)

	checkErr(err)
	doc.Find(".pagination").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		pages = s.Find("a").Length()
		//band := s.Find("a").Text()
		//title, _ := s.Find("a").Attr("href")
		//fmt.Printf("Review %d: %s - %s\n", i, band, title)
	})

	return pages
}

func writeJobs(jobs []extractedJob) {
	file, err := os.Create("jobs.csv")
	checkErr(err)

	w := csv.NewWriter(file)

	defer w.Flush()

	headers := []string{"ID", "Title", "Location", "Salary", "Summary"}

	wErr := w.Write(headers)
	checkErr(wErr)

	for _, job := range jobs {
		jobSlice := []string{"https://kr.indeed.com/%EC%B1%84%EC%9A%A9%EB%B3%B4%EA%B8%B0?jk=" + job.id, job.title, job.location, job.salary, job.summary}
		jwErr := w.Write(jobSlice)
		checkErr(jwErr)
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request failed with Status : ", res.StatusCode)
	}

}

// CleanString is Clen String
func CleanString(text string) string {

	return strings.Join(strings.Fields(strings.TrimSpace(text)), " ")

}
