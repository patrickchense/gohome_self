package main

import (
	"crypto/md5"
	"io"
	"strconv"
	"fmt"
	"time"
)

func main()  {
	//生成
	h := md5.New()
	crutime := time.Now()
	io.WriteString(h, strconv.FormatInt(crutime.Unix(), 10))
	//io.WriteString(h, "ganraomaxxxxxxxxx")
	token := fmt.Sprintf("%x", h.Sum(nil))
	fmt.Println("token:" + token)
/*	t, _ := template.ParseFiles("login.gtpl")
	t.Execute(w, token)*/
	//验证
/*	r.ParseForm()
	token := r.Form.Get("token")*/
	if token != "" {
		//验证token的合法性
	} else {
		//不存在token报错
	}
}
