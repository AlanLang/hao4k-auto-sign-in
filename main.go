package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/simplifiedchinese"
)

func main() {
	gCookiesJar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: gCookiesJar,
	}
	isLogin := LogIn(client)
	if !isLogin {
		os.Exit(3)
	}
	formHash := GetFormHash(client)
	if formHash == "" {
		fmt.Println("formHash获取失败")
		os.Exit(3)
	}
	success := SignIn(client, formHash)
	if success {
		fmt.Println("签到成功")
	} else {
		fmt.Println("签到失败")
		os.Exit(3)
	}
	// cookies := gCookiesJar.Cookies(reqest.URL)
	// fmt.Println(cookies)
}

// LogIn 登录
func LogIn(client *http.Client) bool {
	username := os.Getenv("HAO4K_USERNAME")
	password := os.Getenv("HAO4K_PASSWORD")
	if username == "" {
		fmt.Println("用户名不能为空")
		return false
	}
	if password == "" {
		fmt.Println("密码不能为空")
		return false
	}
	urlString := "https://www.hao4k.cn/member.php?mod=logging&action=login&loginsubmit=yes&loginhash=Lj415&inajax=1"
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	_ = bodyWriter.WriteField("username", username)
	_ = bodyWriter.WriteField("password", password)
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()
	reqest, _ := http.NewRequest("POST", urlString, bodyBuf)
	reqest.Header.Add("content-type", contentType)
	response, err := client.Do(reqest)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	reader := simplifiedchinese.GB18030.NewDecoder().Reader(response.Body)
	buf, err := ioutil.ReadAll(reader)
	isLogin := strings.Contains(string(buf), "欢迎您回来，新手上路")
	if !isLogin {
		fmt.Println("登录失败")
		fmt.Println(string(buf))
	} else {
		fmt.Println("登录成功")
	}

	return isLogin
}

// GetFormHash GetFormHash
func GetFormHash(client *http.Client) string {
	urlString := "https://www.hao4k.cn//k_misign-sign.html"
	reqest, err := http.NewRequest("GET", urlString, nil)
	//处理返回结果
	response, err := client.Do(reqest)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	reader := simplifiedchinese.GB18030.NewDecoder().Reader(response.Body)
	buf, err := ioutil.ReadAll(reader)
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(string(buf)))
	var formHash string
	formHash = ""
	dom.Find("input[name=formhash]").Each(func(i int, selection *goquery.Selection) {
		val, exists := selection.Attr("value")
		if exists {
			formHash = val
		}
	})
	return formHash
}

// SignIn 签到
func SignIn(client *http.Client, formHash string) bool {
	//生成要访问的url
	url := "https://www.hao4k.cn/plugin.php?id=k_misign:sign&operation=qiandao&formhash=" + formHash + "&format=empty&inajax=1&ajaxtarget=JD_sign"
	//提交请求
	reqest, err := http.NewRequest("GET", url, nil)

	//处理返回结果
	response, err := client.Do(reqest)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	reader := simplifiedchinese.GB18030.NewDecoder().Reader(response.Body)
	buf, err := ioutil.ReadAll(reader)
	fmt.Println(string(buf))
	return strings.Contains(string(buf), "![CDATA[")
}
