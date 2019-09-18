package main

import (
	"fmt"
	"github.com/tidwall/gjson"
	"gopkg.in/gcfg.v1"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type hl7 struct {
	Urlhost     string
	Suburl1     string
	Suburl2     string
	Requestmode string
	Cookieid    string
	ContentType string
	Pagesize    string
	Jsonkey1    string
	Jsonkey2    string
	chuancan    string
}

func main() {

	Qingqiu1t := `{"pageNum":1,"pageSize":`
	Qingqiu1b := `,"condition":{"giveUpCount":5,"lastTouchTimeRange":["2019-08-17T08:25:15.342Z","2019-09-15T08:25:15.342Z"],"startPageNum":100,"endPageNum":100,"lastTouchStartDate":"2019-08-17 00:00:00","lastTouchEndDate":"2019-09-15 23:59:59"}}`
	Qingqiu2t := `{"customerId":"`
	Qingqiu2b := `"}`

	Outputfilename := "PN.txt"
	Inifilename := "config.ini"

	hl7data := readinifile(Inifilename)
	hl7data.chuancan = Qingqiu1t + hl7data.Pagesize + Qingqiu1b
	customerId := strings.Split(GetcustomerId(hl7data, getpage(hl7data)), ",")

	hl7data.Suburl1 = hl7data.Suburl2
	PNumberList := ""

	for _, cid1 := range customerId {
		cid1 = strings.Trim(cid1, "[")
		cid1 = strings.Trim(cid1, "]")
		hl7data.chuancan = Qingqiu2t + cid1 + Qingqiu2b
		fmt.Printf("Processed 1 item data.customerId is ")
		fmt.Println(cid1)
		PNumberList = PNumberList + GetPhoneNumber(hl7data, getpage(hl7data)) + "\n"
	}

	//fmt.Println(name2)
	appendToFile(Outputfilename, PNumberList)
}

func getpage(hl7data hl7) string {

	client1 := &http.Client{}
	req, err := http.NewRequest(hl7data.Requestmode, hl7data.Urlhost+hl7data.Suburl1, strings.NewReader(hl7data.chuancan)) //发送请求
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Content-Type", hl7data.ContentType) //请求配置连接类型
	req.Header.Set("Cookie", hl7data.Cookieid)          //请求配置cookie

	resp, err := client1.Do(req) //接收返回信息
	if err != nil {
		fmt.Println(err)
		resp.Body.Close()
	}
	//defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body) //接收返回信息
	if err != nil {
		fmt.Println(err)
	}
	return string(body)
}

func GetcustomerId(hl7data hl7, jsonstr string) string {
	resulta := gjson.Get(jsonstr, hl7data.Jsonkey1)
	return resulta.String()
}

func GetPhoneNumber(hl7data hl7, jsonstr string) string {
	resulta := gjson.Get(jsonstr, hl7data.Jsonkey2)
	return resulta.String()
}

func appendToFile(Outputfilename, str string) {
	f, err := os.OpenFile(Outputfilename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
	if err != nil {
		fmt.Printf("Cannot open file %s!\n")
		return
	}
	defer f.Close()
	f.WriteString(str)
}

func readinifile(Inifilename string) hl7 {
	hl7data := hl7{}
	iniexist, err := os.Stat(Inifilename)

	if iniexist != nil {

		//fmt.Println("Found Configure file.\nConfigure file is [" + configurefile + "].")
		config := struct {
			HL7order struct {
				Urlhost     string
				Suburl1     string
				Suburl2     string
				Requestmode string
				Cookieid    string
				ContentType string
				Pagesize    string
				Jsonkey1    string
				Jsonkey2    string
			}
		}{}

		err := gcfg.ReadFileInto(&config, Inifilename)
		if err != nil {
			fmt.Println(err, "Failed to parseconfigure file!")
		} else {
			hl7data.Urlhost = strings.TrimSpace(config.HL7order.Urlhost)
			hl7data.Suburl1 = strings.TrimSpace(config.HL7order.Suburl1)
			hl7data.Suburl2 = strings.TrimSpace(config.HL7order.Suburl2)
			hl7data.Requestmode = strings.TrimSpace(config.HL7order.Requestmode)
			hl7data.Cookieid = "JSESSIONID=" + strings.TrimSpace(config.HL7order.Cookieid)
			hl7data.ContentType = strings.TrimSpace(config.HL7order.ContentType)
			hl7data.Pagesize = strings.TrimSpace(config.HL7order.Pagesize)
			hl7data.Jsonkey1 = strings.TrimSpace(config.HL7order.Jsonkey1)
			hl7data.Jsonkey2 = strings.TrimSpace(config.HL7order.Jsonkey2)
		}
	} else {
		fmt.Println(err)
	}
	return hl7data

}
