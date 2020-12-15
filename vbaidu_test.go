package vbaidu

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestCrawlUrls(t *testing.T) {
	log.Println("TestCrawlUrls start ... ")
	StartCrawl()

}

func TestResultData(t *testing.T) {
	f, err := os.OpenFile("data/example.json", os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}

	bytes, _ := ioutil.ReadAll(f)
	v := &Result{}
	err = json.Unmarshal(bytes, v)
	if err != nil {
		panic(err)
	}
	log.Println(v)
}
