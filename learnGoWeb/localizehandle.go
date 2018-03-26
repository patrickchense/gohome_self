package main

import (
	"fmt"
	"time"
	//"os"
	//"html/template"
)

var locales map[string]map[string]string

func main() {
	locales = make(map[string]map[string]string, 2)
	en := make(map[string]string, 10)
	en["pea"] = "pea"
	en["bean"] = "bean"
	locales["en"] = en
	cn := make(map[string]string, 10)
	cn["pea"] = "豌豆"
	cn["bean"] = "毛豆"
	locales["zh-CN"] = cn
	lang := "zh-CN"
	fmt.Println(msg(lang, "pea"))
	fmt.Println(msg(lang, "bean"))
	//替换
	en["how old"] ="I am %d years old"
	cn["how old"] ="我今年%d岁了"

	fmt.Printf(msg(lang, "how old"), 30)
	fmt.Println()
	//time
	en["time_zone"]="America/Chicago"
	cn["time_zone"]="Asia/Shanghai"

	loc,_:=time.LoadLocation(msg(lang,"time_zone"))
	t:=time.Now()
	t = t.In(loc)
	fmt.Println(t.Format(time.RFC3339))

	//time formate
	en["date_format"]="%Y-%m-%d %H:%M:%S"
	cn["date_format"]="%Y年%m月%d日 %H时%M分%S秒"

	//fmt.Println(date(msg(lang,"date_format"),t))

	//money
	en["money"] ="USD %d"
	cn["money"] ="￥%d元"

	fmt.Println(money_format(msg(lang,"money"),100))

	//本地化国际化资源 views/lang/file  格式
	//s1, _ := template.ParseFiles("views/"+lang+"/index.tpl")
	//VV.Lang=lang
	//s1.Execute(os.Stdout, VV)
}

func msg(locale, key string) string {
	if v, ok := locales[locale]; ok {
		if v2, ok := v[key]; ok {
			return v2
		}
	}
	return ""
}


func date(fomate string,t time.Time) string{
	//year, month, day = t.Date()
	//hour, min, sec = t.Clock()
	//解析相应的%Y %m %d %H %M %S然后返回信息
	//%Y 替换成2012
	//%m 替换成10
	//%d 替换成24
	return ""
}

func money_format(fomate string,money int64) string{
	return fmt.Sprintf(fomate,money)
}