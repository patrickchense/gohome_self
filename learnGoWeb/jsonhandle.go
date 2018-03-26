package main

import (
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"os"
)

type Server struct {
	ServerName string `json:"serverName"`
	ServerIP  string `json:"serverIP"`
}

type Serverslice struct {
	Servers []Server `json:"servers"`
}

type Server1 struct {
	// ID 不会导出到JSON中
	ID int `json:"-"`

	// ServerName2 的值会进行二次JSON编码
	ServerName  string `json:"serverName"`
	ServerName2 string `json:"serverName2,string"`

	// 如果 ServerIP 为空，则不输出到JSON串中
	ServerIP   string `json:"serverIP,omitempty"`
}


func main() {
	var s Serverslice
	str := `{"servers":[{"serverName":"Shanghai_VPN","serverIP":"127.0.0.1"},{"serverName":"Beijing_VPN","serverIP":"127.0.0.2"}]}`
	json.Unmarshal([]byte(str), &s)
	fmt.Println(s) //{[{Shanghai_VPN 127.0.0.1} {Beijing_VPN 127.0.0.2}]}

	b := []byte(`{"Name":"Wednesday","Age":6,"Parents":["Gomez","Morticia"]}`)
	var f interface{}
	json.Unmarshal(b, &f)
	m := f.(map[string]interface{})

	for k, v := range m {
		switch vv := v.(type) {
		case string:
			fmt.Println(k, "is string", vv)
		case int:
			fmt.Println(k, "is int", vv)
		case float64:
			fmt.Println(k,"is float64",vv)
		case []interface{}:
			fmt.Println(k, "is an array:")
			for i, u := range vv {
				fmt.Println(i, u)
			}
		default:
			fmt.Println(k, "is of a type I don't know how to handle")
		}
	}
	/*
	Age is float64 6
Parents is an array:
0 Gomez
1 Morticia
Name is string Wednesday
	 */
	println("--------------------------------------------------")
	println("simplejson use")
	js, _ := simplejson.NewJson([]byte(`{
	"test": {
		"array": [1, "2", 3],
		"int": 10,
		"float": 5.150,
		"bignum": 9223372036854775807,
		"string": "simplejson",
		"bool": true
	}
}`))

	arr, _ := js.Get("test").Get("array").Array()
	i, _ := js.Get("test").Get("int").Int()
	ms := js.Get("test").Get("string").MustString()
	println(arr)
	println(i)
	println(ms)

	println("---------------------")
	var s1 Serverslice
	s1.Servers = append(s1.Servers, Server{ServerName: "Shanghai_VPN", ServerIP: "127.0.0.1"})
	s1.Servers = append(s1.Servers, Server{ServerName: "Beijing_VPN", ServerIP: "127.0.0.2"})
	b, err := json.Marshal(s1)
	if err != nil {
		fmt.Println("json err:", err)
	}
	fmt.Printf("generate json: %s\n", string(b))
	// 没用 tag : {"Servers":[{"ServerName":"Shanghai_VPN","ServerIP":"127.0.0.1"},{"ServerName":"Beijing_VPN","ServerIP":"127.0.0.2"}]}
	// 用了tag: {"servers":[{"serverName":"Shanghai_VPN","serverIP":"127.0.0.1"},{"serverName":"Beijing_VPN","serverIP":"127.0.0.2"}]}

	println("-------------------------------------")
	//测试tag 的使用 string 转换, omitempty,
	s2 := Server1 {
		ID:         3,
		ServerName:  `Go "1.0" `,
		ServerName2: `Go "1.0" `,
		ServerIP:   ``,
	}
	b1, _ := json.Marshal(s2)
	os.Stdout.Write(b1) //{"serverName":"Go \"1.0\" ","serverName2":"\"Go \\\"1.0\\\" \""}

}
