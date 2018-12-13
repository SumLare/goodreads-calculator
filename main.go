package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const DaysInYear = 365

type Reviews struct {
	XMLName xml.Name `xml:"GoodreadsResponse"`
	Reviews []Review `xml:"reviews>review"`
}

type Review struct {
	XMLName  xml.Name `xml:"review"`
	NumPages int      `xml:"book>num_pages"`
}

func main() {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.goodreads.com/review/list", nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	q := req.URL.Query()
	key := flag.String("key", "", "")
	id := flag.String("id", "", "")
	flag.Parse()

	q.Add("key", *key)
	q.Add("v", "2")
	q.Add("id", *id)
	q.Add("shelf", "to-read")
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Errorf("Read body: %v", err)
	}
	calculate(data)
}

func calculate(data []byte) {
	var reviews Reviews
	err := xml.Unmarshal(data, &reviews)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	pages := 0
	for _, review := range reviews.Reviews {
		pages += review.NumPages
	}
	avg := pages / DaysInYear
	fmt.Println(avg)
}
