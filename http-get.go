package main

import (
	"encoding/csv"
	"fmt"
	"github.com/tidwall/gjson"
	"gopkg.in/gcfg.v1"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
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
	Jsonkey3    string
	chuancan    string
}

func main() {

	Qingqiu1t := `{"pageNum":1,"pageSize":`
	Qingqiu1b := `,"condition":{"giveUpCount":5,"lastTouchTimeRange":["2019-08-17T08:25:15.342Z","2019-09-15T08:25:15.342Z"],"startPageNum":100,"endPageNum":100,"lastTouchStartDate":"2019-08-17 00:00:00","lastTouchEndDate":"2019-09-15 23:59:59"}}`
	Qingqiu2t := `{"customerId":"`
	Qingqiu2b := `"}`

	Outputfilename := "PN.CSV"
	Inifilename := "config.ini"

	hl7data := readinifile(Inifilename)
	hl7data.chuancan = Qingqiu1t + hl7data.Pagesize + Qingqiu1b
	customerId := strings.Split(Getdata(hl7data, getpage(hl7data)), ",")

	hl7data.Jsonkey1 = hl7data.Jsonkey3
	customerName := strings.Split(Getdata(hl7data, getpage(hl7data)), ",")

	hl7data.Suburl1 = hl7data.Suburl2
	hl7data.Jsonkey1 = hl7data.Jsonkey2
	//PNumberList := ""
	count1 := 0

	f, err := os.Create(Outputfilename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.WriteString("\xEF\xBB\xBF")
	w := csv.NewWriter(f)
	w.Write([]string{"学员姓名", "学员性别", "学员手机号", "qq", "微信", "项目(必填)", "学历", "年龄", "证件类型", "证件号码", "客户来源", "创建人", "创建时间(yyyy-MM-dd HH:mm:ss)", "地域(必填)", "归属人", "回访次数", "下次回访时间(yyyy-MM-dd HH:mm:ss)", "备注"})

	for _, cid1 := range customerId {
		cid1 = strings.Trim(cid1, "[")
		cid1 = strings.Trim(cid1, "]")
		hl7data.chuancan = Qingqiu2t + cid1 + Qingqiu2b
		//PNumberList = customerName[count1] + "," + Getdata(hl7data, getpage(hl7data)) + "\n"
		fmt.Println("Processed " + strconv.Itoa(count1) + "/" + hl7data.Pagesize + " item data.customerId is " + cid1)
		//appendToFile(Outputfilename, PNumberList)
		//fmt.Println(PNumberList)
		//fmt.Println(strings.Trim(strings.Trim(strings.Trim(customerName[count1],`[`),`]`),`"`),Getdata(hl7data, getpage(hl7data)))
		w.Write([]string{strings.Trim(strings.Trim(strings.Trim(customerName[count1], `[`), `]`), `"`), "", Getdata(hl7data, getpage(hl7data))})
		w.Flush()
		count1 = count1 + 1
	}
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

func Getdata(hl7data hl7, jsonstr string) string {
	resulta := gjson.Get(jsonstr, hl7data.Jsonkey1)
	return resulta.String()
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
				Jsonkey3    string
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
			hl7data.Jsonkey3 = strings.TrimSpace(config.HL7order.Jsonkey3)
		}
	} else {
		fmt.Println(err)
	}
	return hl7data

}
