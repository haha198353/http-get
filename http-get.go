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
	LinkOption struct {
		Urlhost     string
		Suburl1     string
		Suburl2     string
		Requestmode string
		Cookieid    string
		ContentType string
		Jsonkey0    string
		Jsonkey1    string
		Jsonkey2    string
		Jsonkey3    string
		Codetype    string
	}
	RquestOption struct {
		PageNum             string
		PageSize            string
		GiveUpCount         string
		LastTouchTimeRange1 string
		LastTouchTimeRange2 string
		LastTouchTimeRange3 string
		LastTouchTimeRange4 string
		StartPageNum        string
		EndPageNum          string
		LastTouchStartDate  string
		LastTouchEndDate    string
		chuancan            string
	}
}

func main() {

	t := time.Now()                                       //取当前时间，精确到秒
	Outputfilename := t.Format("20060102150405") + ".CSV" //定义输出文件名，避免覆盖
	Inifilename := "config.ini"                           //读取配置文件的文件名，写死不能更改

	hl7data := readinifile(Inifilename) //读取配置文件

	hl7data.RquestOption.chuancan = `{` +
		`"pageNum": ` + hl7data.RquestOption.PageNum + `,` +
		`"pageSize": ` + hl7data.RquestOption.PageSize + `,` +
		`"condition": {` +
		`"giveUpCount": ` + hl7data.RquestOption.GiveUpCount + `,` +
		`"lastTouchTimeRange": ["` + hl7data.RquestOption.LastTouchTimeRange1 + `T` + hl7data.RquestOption.LastTouchTimeRange2 + `.342Z", "` + hl7data.RquestOption.LastTouchTimeRange3 + `T` + hl7data.RquestOption.LastTouchTimeRange4 + `.342Z"],` +
		`"startPageNum": ` + hl7data.RquestOption.StartPageNum + `,` +
		`"endPageNum": ` + hl7data.RquestOption.EndPageNum + `,` +
		`"lastTouchStartDate": "` + hl7data.RquestOption.LastTouchStartDate + `00:00:00",` +
		`"lastTouchEndDate": "` + hl7data.RquestOption.LastTouchEndDate + `23:59:59"` +
		`}` +
		`}`
	customerId := strings.Split(Getdata(hl7data, getpage(hl7data)), ",") //获取 customerId 并切成数组

	hl7data.LinkOption.Jsonkey1 = hl7data.LinkOption.Jsonkey3              //重新配置请求参数
	customerName := strings.Split(Getdata(hl7data, getpage(hl7data)), ",") //获取 customerName 并切成数组

	hl7data.LinkOption.Suburl1 = hl7data.LinkOption.Suburl2   //重新配置请求参数1
	hl7data.LinkOption.Jsonkey1 = hl7data.LinkOption.Jsonkey2 //重新配置请求参数2

	f, err := os.Create(Outputfilename) //创建输出文件
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.WriteString("\xEF\xBB\xBF") //写入输出文件格式文件头
	w := csv.NewWriter(f)

	enc := mahonia.NewEncoder(hl7data.LinkOption.Codetype) //GO默认编码方式为UTF-8，但是Windows识别编码默认为ANSI（简中为GBK），故Windows下使用需要转码，enc.ConvertString为具体转码实现

	w.Write([]string{enc.ConvertString("序号"), enc.ConvertString("CID"), enc.ConvertString("学员姓名"), enc.ConvertString("学员性别"), enc.ConvertString("学员手机号"), enc.ConvertString("qq"), enc.ConvertString("微信"), enc.ConvertString("项目(必填)"), enc.ConvertString("学历"), enc.ConvertString("年龄"), enc.ConvertString("证件类型"), enc.ConvertString("证件号码"), enc.ConvertString("客户来源"), enc.ConvertString("创建人"), enc.ConvertString("创建时间(yyyy-MM-dd HH:mm:ss)"), enc.ConvertString("地域(必填)"), enc.ConvertString("归属人"), enc.ConvertString("回访次数"), enc.ConvertString("下次回访时间(yyyy-MM-dd HH:mm:ss)"), enc.ConvertString("备注")}) //标题行

	for count1, cid1 := range customerId { //此处可以声明两个变量，第一个是数组位置的值，第二个代表该数组的值
		cid1 = strings.Trim(strings.Trim(cid1, `[`), `]`)
		countstr := strconv.Itoa(count1 + 1)
		hl7data.RquestOption.chuancan = `{"` + hl7data.LinkOption.Jsonkey0 + `":"` + cid1 + `"}`                        //重新配置请求参数
		fmt.Println("Processed " + countstr + "/" + hl7data.RquestOption.PageSize + " item data.customerId is " + cid1) //屏幕输出，提示进度

		go w.Write([]string{enc.ConvertString(strconv.Itoa(count1 + 1)), enc.ConvertString(cid1), enc.ConvertString(strings.Trim(strings.Trim(strings.Trim(customerName[count1], `[`), `]`), `"`)), "", Getdata(hl7data, getpage(hl7data))}) //获取 phone1 并和 customerName 、customerId 一并写入文件                                                                                                                                                                                                                      //文件写入刷新
	}
	w.Flush()
}

func getpage(hl7data hl7) string {

	client1 := &http.Client{}
	req, err := http.NewRequest(hl7data.LinkOption.Requestmode, hl7data.LinkOption.Urlhost+hl7data.LinkOption.Suburl1, strings.NewReader(hl7data.RquestOption.chuancan)) //发送请求
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Content-Type", hl7data.LinkOption.ContentType) //请求配置连接类型
	req.Header.Set("Cookie", hl7data.LinkOption.Cookieid)          //请求配置cookie

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
	resulta := gjson.Get(jsonstr, hl7data.LinkOption.Jsonkey1)
	return resulta.String()
}

func readinifile(Inifilename string) hl7 {
	hl7data := hl7{}
	iniexist, err := os.Stat(Inifilename)

	if iniexist != nil {

		config := hl7{}
		err := gcfg.ReadFileInto(&config, Inifilename) //读取配置文件内容，项目必须一一严格对应
		if err != nil {
			fmt.Println(err, "Failed to parseconfigure file!")
		} else {
			hl7data.LinkOption.Urlhost = strings.TrimSpace(config.LinkOption.Urlhost)
			hl7data.LinkOption.Suburl1 = strings.TrimSpace(config.LinkOption.Suburl1)
			hl7data.LinkOption.Suburl2 = strings.TrimSpace(config.LinkOption.Suburl2)
			hl7data.LinkOption.Requestmode = strings.TrimSpace(config.LinkOption.Requestmode)
			hl7data.LinkOption.Cookieid = "JSESSIONID=" + strings.TrimSpace(config.LinkOption.Cookieid)
			hl7data.LinkOption.ContentType = strings.TrimSpace(config.LinkOption.ContentType)
			hl7data.LinkOption.Jsonkey0 = strings.TrimSpace(config.LinkOption.Jsonkey0)
			hl7data.LinkOption.Jsonkey1 = strings.TrimSpace(config.LinkOption.Jsonkey1)
			hl7data.LinkOption.Jsonkey2 = strings.TrimSpace(config.LinkOption.Jsonkey2)
			hl7data.LinkOption.Jsonkey3 = strings.TrimSpace(config.LinkOption.Jsonkey3)
			hl7data.LinkOption.Codetype = strings.TrimSpace(config.LinkOption.Codetype)

			hl7data.RquestOption.PageNum = strings.TrimSpace(config.RquestOption.PageNum)
			hl7data.RquestOption.PageSize = strings.TrimSpace(config.RquestOption.PageSize)
			hl7data.RquestOption.GiveUpCount = strings.TrimSpace(config.RquestOption.GiveUpCount)
			hl7data.RquestOption.LastTouchTimeRange1 = strings.TrimSpace(config.RquestOption.LastTouchTimeRange1)
			hl7data.RquestOption.LastTouchTimeRange2 = strings.TrimSpace(config.RquestOption.LastTouchTimeRange2)
			hl7data.RquestOption.LastTouchTimeRange3 = strings.TrimSpace(config.RquestOption.LastTouchTimeRange3)
			hl7data.RquestOption.LastTouchTimeRange4 = strings.TrimSpace(config.RquestOption.LastTouchTimeRange4)
			hl7data.RquestOption.StartPageNum = strings.TrimSpace(config.RquestOption.StartPageNum)
			hl7data.RquestOption.EndPageNum = strings.TrimSpace(config.RquestOption.EndPageNum)
			hl7data.RquestOption.LastTouchStartDate = strings.TrimSpace(config.RquestOption.LastTouchStartDate)
			hl7data.RquestOption.LastTouchEndDate = strings.TrimSpace(config.RquestOption.LastTouchEndDate)
		}
	} else {
		fmt.Println(err)
	}
	return hl7data
}

// func appendToFile(Outputfilename, str string) {  //追加方式写入文件
// 	f, err := os.OpenFile(Outputfilename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
// 	if err != nil {
// 		fmt.Printf("Cannot open file %s!\n")
// 		return
// 	}
// 	defer f.Close()
// 	f.WriteString(str)
// }
