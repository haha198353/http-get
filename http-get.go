package main

import (
	"encoding/csv"
	"fmt"
	"github.com/axgle/mahonia"
	"github.com/tidwall/gjson"
	"gopkg.in/gcfg.v1"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type hl7 struct {
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
		Jsonkey3    string
		Codetype    string
		chuancan    string
	}
}

func main() {

	Qingqiu1t := `{"pageNum":1,"pageSize":`
	Qingqiu1b := `,"condition":{"giveUpCount":5,"lastTouchTimeRange":["2019-08-17T08:25:15.342Z","2019-09-15T08:25:15.342Z"],"startPageNum":100,"endPageNum":100,"lastTouchStartDate":"2019-08-17 00:00:00","lastTouchEndDate":"2019-09-15 23:59:59"}}`
	Qingqiu2t := `{"customerId":"`
	Qingqiu2b := `"}`

	t := time.Now()
	Outputfilename := t.Format("20060102150405") + ".CSV"
	Inifilename := "config.ini"

	hl7data := readinifile(Inifilename)
	hl7data.HL7order.chuancan = Qingqiu1t + hl7data.HL7order.Pagesize + Qingqiu1b
	customerId := strings.Split(Getdata(hl7data, getpage(hl7data)), ",")

	hl7data.HL7order.Jsonkey1 = hl7data.HL7order.Jsonkey3
	customerName := strings.Split(Getdata(hl7data, getpage(hl7data)), ",")

	hl7data.HL7order.Suburl1 = hl7data.HL7order.Suburl2
	hl7data.HL7order.Jsonkey1 = hl7data.HL7order.Jsonkey2

	f, err := os.Create(Outputfilename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.WriteString("\xEF\xBB\xBF")
	w := csv.NewWriter(f)

	enc := mahonia.NewEncoder(hl7data.HL7order.Codetype) //GO默认编码方式为UTF-8，但是Windows识别编码默认为ANSI（简中为GBK），故Windows下使用需要转码，enc.ConvertString为具体转码实现

	w.Write([]string{enc.ConvertString("序号"), enc.ConvertString("CID"), enc.ConvertString("学员姓名"), enc.ConvertString("学员性别"), enc.ConvertString("学员手机号"), enc.ConvertString("qq"), enc.ConvertString("微信"), enc.ConvertString("项目(必填)"), enc.ConvertString("学历"), enc.ConvertString("年龄"), enc.ConvertString("证件类型"), enc.ConvertString("证件号码"), enc.ConvertString("客户来源"), enc.ConvertString("创建人"), enc.ConvertString("创建时间(yyyy-MM-dd HH:mm:ss)"), enc.ConvertString("地域(必填)"), enc.ConvertString("归属人"), enc.ConvertString("回访次数"), enc.ConvertString("下次回访时间(yyyy-MM-dd HH:mm:ss)"), enc.ConvertString("备注")})

	for count1, cid1 := range customerId { //此处可以声明两个变量，第一个是数组位置的值，第二个代表该数组的值
		cid1 = strings.Trim(strings.Trim(cid1, `[`), `]`)
		hl7data.HL7order.chuancan = Qingqiu2t + cid1 + Qingqiu2b
		fmt.Println("Processed " + strconv.Itoa(count1+1) + "/" + hl7data.HL7order.Pagesize + " item data.customerId is " + cid1)

		w.Write([]string{enc.ConvertString(strconv.Itoa(count1 + 1)), enc.ConvertString(cid1), enc.ConvertString(strings.Trim(strings.Trim(strings.Trim(customerName[count1], `[`), `]`), `"`)), "", Getdata(hl7data, getpage(hl7data))})
		w.Flush()
	}
}

func getpage(hl7data hl7) string {

	client1 := &http.Client{}
	req, err := http.NewRequest(hl7data.HL7order.Requestmode, hl7data.HL7order.Urlhost+hl7data.HL7order.Suburl1, strings.NewReader(hl7data.HL7order.chuancan)) //发送请求
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Content-Type", hl7data.HL7order.ContentType) //请求配置连接类型
	req.Header.Set("Cookie", hl7data.HL7order.Cookieid)          //请求配置cookie

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

func Getdata(hl7data hl7, jsonstr string) string {
	resulta := gjson.Get(jsonstr, hl7data.HL7order.Jsonkey1)
	return resulta.String()
}

func readinifile(Inifilename string) hl7 {
	hl7data := hl7{}
	iniexist, err := os.Stat(Inifilename)

	if iniexist != nil {

		//fmt.Println("Found Configure file.\nConfigure file is [" + configurefile + "].")
		config := hl7{}
		err := gcfg.ReadFileInto(&config, Inifilename)
		if err != nil {
			fmt.Println(err, "Failed to parseconfigure file!")
		} else {
			hl7data.HL7order.Urlhost = strings.TrimSpace(config.HL7order.Urlhost)
			hl7data.HL7order.Suburl1 = strings.TrimSpace(config.HL7order.Suburl1)
			hl7data.HL7order.Suburl2 = strings.TrimSpace(config.HL7order.Suburl2)
			hl7data.HL7order.Requestmode = strings.TrimSpace(config.HL7order.Requestmode)
			hl7data.HL7order.Cookieid = "JSESSIONID=" + strings.TrimSpace(config.HL7order.Cookieid)
			hl7data.HL7order.ContentType = strings.TrimSpace(config.HL7order.ContentType)
			hl7data.HL7order.Pagesize = strings.TrimSpace(config.HL7order.Pagesize)
			hl7data.HL7order.Jsonkey1 = strings.TrimSpace(config.HL7order.Jsonkey1)
			hl7data.HL7order.Jsonkey2 = strings.TrimSpace(config.HL7order.Jsonkey2)
			hl7data.HL7order.Jsonkey3 = strings.TrimSpace(config.HL7order.Jsonkey3)
			hl7data.HL7order.Codetype = strings.TrimSpace(config.HL7order.Codetype)
		}
	} else {
		fmt.Println(err)
	}
	return hl7data
}

// func appendToFile(Outputfilename, str string) {
// 	f, err := os.OpenFile(Outputfilename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
// 	if err != nil {
// 		fmt.Printf("Cannot open file %s!\n")
// 		return
// 	}
// 	defer f.Close()
// 	f.WriteString(str)
// }
