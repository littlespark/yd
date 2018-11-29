package main

import (
	"encoding/json"
	"fmt"
	"os"

	"io/ioutil"
	"net/http"
	"strings"

	"net/url"

	"crypto/md5"

	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
)

var (
	//此处抄了https://github.com/TimothyYe/ydict的创意
	Version = "0.1"
	logo    = `
 __   __  ____    
  \ \ / / |  _"\   
   \ V / /| | | |  
  U_|"|_uU| |_| |\ 
    |_|   |____/ u 
.-,//|(_   |||_    
 \_) (__) (__)_)     

YD V%s
https://github.com/littlespark/go-youdao

`

	word, response, reqBody = "", "", ""

	openApiUrl = "http://openapi.youdao.com/api"
	appKey     = "7ba6a5d54e032df3"
	secertKey  = "s5jxz5MWms4d10DHdQOoLp7xc5Y6MA3d"
)

func main() {

	//.获取终端写入的字符
	checkInput()

	//. 显示执行进度条
	s := spinner.New(spinner.CharSets[37], 30*time.Millisecond)
	s.Color("green", "bold")
	s.Start()

	//.组装调用数据
	build()

	//.调用有道翻译http api并解析返回结果
	httpPost()

	//. 关闭执行进度条
	s.Stop()

	//.格式化返回内容，并美化输出 + 无翻译结果时输出错误
	output()

}

func checkInput() string {
	args := os.Args
	if len(args) != 2 {
		color.Cyan(logo, Version)
		color.Cyan("Usage:")
		color.Cyan("yd <word(s) to query>        Query the word(s)")
		os.Exit(0)
	} else {
		word = args[1]
	}

	return word
}

func build() {
	//参数较少直接用string拼装
	sign := appKey + word + "salt" + secertKey //md5(appKey+q+salt+应用密钥)
	reqBody = "appKey=" + appKey + "&from=auto&to=auto" + "&q=" + url.QueryEscape(word) + "&salt=salt" + "&sign=" + encrypt(sign)
}

func encrypt(str string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}

func httpPost() {
	resp, err := http.Post(openApiUrl, "application/x-www-form-urlencoded", strings.NewReader(reqBody))
	if err != nil {
		color.Blue("\r\n    word '%s' not found", word)
		os.Exit(0)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		color.Blue("\r\n    word '%s' not found", word)
		os.Exit(0)
	}

	response = string(body)

}

func output() {
	var t TransResult

	if err := json.Unmarshal([]byte(response), &t); err != nil {
		color.Blue("\r\n    word '%s' not found", word)
		os.Exit(0)
	}

	if len(t.Basic.UsPhonetic) > 0 && len(t.Basic.UkPhonetic) > 0 {
		color.Green("\r   英:[%v]      美:[%v]", t.Basic.UkPhonetic, t.Basic.UsPhonetic)
	}

	if len(t.Basic.Explains) > 0 {
		fmt.Println("")
		for _, v := range t.Basic.Explains {
			color.Green("\r    %v", v)
		}
	} else {
		color.Red("\r\n    word '%s' not found", word)
	}

	os.Exit(0)

}

//https://mholt.github.io/json-to-go/
type TransResult struct {
	Basic struct {
		UsPhonetic string   `json:"us-phonetic"`
		UkPhonetic string   `json:"uk-phonetic"`
		Explains   []string `json:"explains"`
	}
}

