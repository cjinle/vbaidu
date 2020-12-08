package conf

import (
	"gopkg.in/ini.v1"
	"strconv"
)

type MainType struct {
	DownloadDir string
}

type CatType struct {
	MaxPageNum int
	PageUrl    string
}

type VbaiduConf struct {
	Main    MainType
	Xiaopin CatType
	Junshi  CatType
}

var VConf *VbaiduConf

func init() {
	ct := CatType{10, ""}
	VConf = &VbaiduConf{
		Main: MainType{
			DownloadDir: "data/",
		},
		Xiaopin: ct,
		Junshi:  ct,
	}
	f, err := ini.Load("conf/vbaidu.ini")
	if err != nil {
		panic(err)
	}

	if val := f.Section("main").Key("download_dir").String(); val != "" {
		VConf.Main.DownloadDir = val
	}
	if val := f.Section("xiaopin").Key("max_page_num").String(); val != "" {
		VConf.Xiaopin.MaxPageNum, _ = strconv.Atoi(val)
	}
	if val := f.Section("xiaopin").Key("page_url").String(); val != "" {
		VConf.Xiaopin.PageUrl = val
	}
	if val := f.Section("junshi").Key("max_page_num").String(); val != "" {
		VConf.Junshi.MaxPageNum, _ = strconv.Atoi(val)
	}
	if val := f.Section("junshi").Key("page_url").String(); val != "" {
		VConf.Junshi.PageUrl = val
	}
}
