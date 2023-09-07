package alertmode

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

func RequestDingding(dingdingUrl, jsonData string) {
	//jsonData := `{"msgtype": "markdown","markdown": {"title":"====侦测到故障====","text":"====侦测到故障==== \n\n 即将删除资源:+警告"}}`
	client := &http.Client{}
	req, err := http.NewRequest("POST", dingdingUrl, bytes.NewBuffer([]byte(jsonData)))
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/json")
	// 发送请求并获取响应
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("发送HTTP请求失败:", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应失败:", err)
		return
	}

	// 打印响应内容
	fmt.Println("响应状态码:", resp.Status)
	fmt.Println("响应内容:", string(body))
}
