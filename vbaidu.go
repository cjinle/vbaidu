package vbaidu

import (
	"encoding/json"
	"fmt"
	. "github.com/cjinle/vbaidu/conf"
	"net/http"
	"log"
	"io/ioutil"
	// "os"
	"time"
)

type Video struct {
	BlockSign  string `json:"block_sign"`
	Wid        string `json:"wid"`
	Title      string `json:"title"`
	Url        string `json:"url"`
	Vid        string `json:"vid"`
	Hao123Url  string `json:"hao123_url"`
	DubaUrl    string `json:"duba_url"`
	ImgvUrl    string `json:"imgv_url"`
	Duration   string `json:"duration"`
	Episode    string `json:"episode"`
	BeginTime  string `json:"begin_time"`
	EndTime    string `json:"end_time"`
	UpdateTime string `json:"update_time"`
	Type       string `json:"type"`
	IsBaishi   string `json:"is_baishi"`
	PlayLink   string `json:"play_link"`
	PlayNum    string `json:"play_num"`
}

type ResultData struct {
	Videos []Video `json:videos`
}

type Result struct {
	ErrorNo int        `json:"errno"`
	Msg     string     `json:"msg"`
	Data    ResultData `json:"data"`
}

var finishChan chan int

func CrawlUrls(cat *CatType, videoChan chan Video) {
	maxPageNum := cat.MaxPageNum
	page := 1
	now := time.Now().Unix()
	for page <= maxPageNum {
		url := fmt.Sprintf(cat.PageUrl, page, now)
		page++
		fmt.Println(url)
		res, err := http.Get(url)
		if err != nil {
			log.Println(err)
			continue
		}
		bytes, _ := ioutil.ReadAll(res.Body)
		res.Body.Close()
		v := &Result{}
		err = json.Unmarshal(bytes, v)
		if err != nil {
			log.Println(err, string(bytes))
			continue
		}
		for _, video := range v.Data.Videos {
			videoChan <- video
		}
	}
	log.Println("finish")
	finishChan <- 1
}

func StartCrawl() {
	videoChan := make(chan Video)
	finishChan := make(chan int)
	go CrawlUrls(&VConf.Xiaopin, videoChan)

	for {
		select {
		case v := <- videoChan:
			log.Println(v)
		case <- finishChan:
			break
		}
	}
}
