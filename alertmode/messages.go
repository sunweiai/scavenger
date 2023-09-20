package alertmode

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"scavenger/utils"
	"strings"
	"time"
)

type MarkdownMes struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type MarkdownType struct {
	TextMes *MarkdownMes `json:"markdown"`
	MesType string       `json:"msgtype"`
}

type DingTalk struct {
	Url  string        `json:"url"`
	Body *MarkdownType `json:"body"`
}

func UnmarshalJson(jsondir string) DingTalk {
	//file, err := os.Open("E:\\gitee_code\\go_project\\scavenger\\dingtalk-mes.json")
	file, err := os.Open(jsondir)
	if err != nil {
		fmt.Println("dingtalk json file do not open.", err)
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Error read json file.", err)
	}

	// 将json中的占位符替换
	configJson := string(data)
	// 解析JSON数据到结构体
	var dingtalk DingTalk
	err = json.Unmarshal([]byte(configJson), &dingtalk)
	if err != nil {
		fmt.Println("解析JSON失败:", err)
	}
	return dingtalk
}

func MarkdownBody(dingtalk DingTalk, metricsInfo []utils.MetricsInfo) (string, string) {
	sourceMesList := make([]string, 0)
	var result string
	if len(metricsInfo) > 0 {
		for _, metrics := range metricsInfo {
			sourceMes := strings.ReplaceAll(dingtalk.Body.TextMes.Text, "{{namespace}}", metrics.Namespace)
			sourceMes = strings.ReplaceAll(sourceMes, "{{sourcename}}", metrics.Pod)
			sourceMes = strings.ReplaceAll(sourceMes, "{{timestamp}}", time.DateTime)
			sourceMes = strings.ReplaceAll(sourceMes, ";", "\n\n")
			sourceMesList = append(sourceMesList, sourceMes)
			sourceMesList = append(sourceMesList, "\n\n----------\n\n")
		}
		for _, sourceMes := range sourceMesList {
			result += sourceMes
		}
		dingtalk.Body.TextMes.Text = result
	}
	//fmt.Printf("url:%v,body:%v\n", string(dingtalk.Url), dingtalk.Body)
	markdown, err := json.Marshal(dingtalk.Body)
	if err != nil {
		panic(err)
	}
	var url string
	if dingtalk.Url != "" {
		url = dingtalk.Url
	} else {
		fmt.Println("Error get dingtalk url.")
	}

	return string(markdown), url
}
