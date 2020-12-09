package vbaidu

import (
	"encoding/json"
	"fmt"
	. "github.com/cjinle/vbaidu/conf"
	"net/http"
	"net/url"
	"log"
	"io/ioutil"
	// "os"
	"os/exec"
	"sync"
	"time"
	"regexp"
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
	VideoUrl   string `json:"video_url`
}

type ResultData struct {
	Videos []*Video `json:videos`
}

type Result struct {
	ErrorNo int        `json:"errno"`
	Msg     string     `json:"msg"`
	Data    ResultData `json:"data"`
}


// var videoChan chan *Video
var wg sync.WaitGroup

func CrawlUrls(cat *CatType) ([]*Video) {
	ret := []*Video{}
	maxPageNum := cat.MaxPageNum
	page := 1
	now := time.Now().Unix()
	for page <= maxPageNum {
		url := fmt.Sprintf(cat.PageUrl, page, now)
		page++
		log.Println(url)
		res, err := http.Get(url)
		if err != nil {
			log.Println(err, url)
			continue
		}
		bytes, _ := ioutil.ReadAll(res.Body)
		res.Body.Close()
		v := &Result{}
		err = json.Unmarshal(bytes, v)
		if err != nil {
			log.Println(err, string(bytes), page)
			continue
		}
		ret = append(ret, v.Data.Videos...)

	}
	return ret
}


func StartCrawl() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	videoList := CrawlUrls(&VConf.Xiaopin)
	// log.Println(videoList)
	if len(videoList) == 0 {
		return
	}

	// videoChan = make(chan *Video)
	for _, video := range videoList {
		log.Println(video.Url)
		go ParseVideoUrl(video)
		wg.Add(1)
		// ParseVideoUrl(video)
	}
	// for v := range videoChan {
	// 	log.Println(v.VideoUrl)
	// }
	wg.Wait()
	log.Println("done")
}

func ParseVideoUrl(video *Video) {
	defer wg.Done()
	res, err := http.Get(video.Url)
	if err != nil {
		log.Println(err)
		return
	}
	defer res.Body.Close()
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return
	}
	reg := `videoFlashPlayUrl\s=\s\'(.*)\'`
	submatch := regexp.MustCompile(reg).FindAllSubmatch(bytes, -1)
	if len(submatch) > 0 && len(submatch[0]) > 0 {
		vv := submatch[0][1]
		urlInfo, err := url.Parse(string(vv))
		if err != nil {
			log.Println(err, string(vv))
			return
		}
		log.Println(urlInfo.Query()["video"][0])
		video.VideoUrl = urlInfo.Query()["video"][0]
		go DownloadVideo(video)
		wg.Add(1)
		// videoChan <- video
	}
}

func DownloadVideo(video *Video) {
	defer wg.Done()
	if video.VideoUrl == "" {
		return
	}
	// $cmd = sprintf("ffmpeg -i \"%s\" -c copy \"%s/%s.mp4\"", $param['video'], $dir, $val['title']);
	// cmd := fmt.Sprintf("ffmpeg -i \"%s\" -c copy \"%s/%s.mp4\"", 
	// 	video.VideoUrl, VConf.Main.DownloadDir, video.Title)
	// log.Println(cmd)
	cmd := exec.Command("ffmpeg", "-i", `"`+video.VideoUrl+`"`, 
		"-c copy", fmt.Sprintf(`"%s/%s.mp4"`, VConf.Main.DownloadDir, video.Title))
	log.Println(cmd.Path, cmd.Args)
	err := cmd.Start()
	if err != nil {
		log.Println(err)
		return
	}
	
}
