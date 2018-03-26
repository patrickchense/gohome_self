package main

import (
	"net/http"
	"fmt"
	"strings"
	"log"
	"html/template"
	"strconv"
	"regexp"
	"learnGoWeb/tools"
	logging "github.com/op/go-logging"
)

import (
	tt "text/template"
	"time"
	"crypto/md5"
	"io"
)


type MyMux struct {

}

/*
handler 接口 定义方法 ServeHTTP(w, r)
通过HandlerFunc 来实现
自己实现ServeHTTP 方法之后， 可以通过http.ListenAndServe 来注册handler
 */
func (p *MyMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if (r.URL.Path == "/") {
		sayHelloName(w, r)
		return
	}
	//http://localhost:9090/testGetScript
	if (r.URL.Path == "/testGetScript") {
/*		t, err := tt.New("foo").Parse(`{{define "T"}}Hello, {{.}}!{{end}}`)
		err = t.ExecuteTemplate(w, "T", "<script>alert('you have been pwned')</script>")
		if err != nil {
			fmt.Println("output script failed")
		}
		//output: Hello, <script>alert('you have been pwned')</script>!
*/
//		或者:
//		html/template
/*		t, err := template.New("foo").Parse(`{{define "T"}}Hello, {{.}}!{{end}}`)
err = t.ExecuteTemplate(w, "T", template.HTML("<script>alert('you have been pwned')</script>"))
		//output the same as above : Hello, <script>alert('you have been pwned')</script>!
*/

		t, err := template.New("foo").Parse(`{{define "T"}}Hello, {{.}}!{{end}}`)
		err = t.ExecuteTemplate(w, "T", "<script>alert('you have been pwned')</script>")
		//output : Hello, &lt;script&gt;alert(&#39;you have been pwned&#39;)&lt;/script&gt;!
		if err != nil {
			fmt.Println("output script failed")
		}
		return
	}
	http.NotFound(w, r)
	NotFound404(w, r) // 错误处理
	return
}

func sayHelloName(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello myroute!")
	r.ParseForm()       //解析url传递的参数，对于POST则解析响应包的主体（request body）
	//注意:如果没有调用ParseForm方法，下面无法获取表单的数据
	fmt.Println(r.Form) //这些信息是输出到服务器端的打印信息
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hello astaxie!") //这个写入到w的是输出到客户端的
}

func validate(r *http.Request) (res bool, err error) {
	if len(r.Form["username"]) == 0 {
		// 验证字符串空
		fmt.Println("no username passed!")
		res = false
		return
	}
	// if age not exists in Form ["age"] error, have to use Get
	age,err:=strconv.Atoi(r.Form.Get("age"))
	if err!=nil{
		//数字转化出错了，那么可能就不是数字
		fmt.Println(" age not found")
		res = false
		return
	}
	fmt.Printf("age is %d\n", age)
	if m, _ := regexp.MatchString("^[0-9]+$", r.Form.Get("age")); !m {
		res = false
		return
	}
	//中文验证
	if m, _ := regexp.MatchString("^\\p{Han}+$", r.Form.Get("realname")); !m {
		fmt.Println("不是中文")
		res = false
		return
	}
	if m, _ := regexp.MatchString("^[a-zA-Z]+$", r.Form.Get("engname")); !m {
		fmt.Println("not english")
		res = false
		return
	}
	//email validate
	if m, _ := regexp.MatchString(`^([\w\.\_]{2,10})@(\w{1,}).([a-z]{2,4})$`, r.Form.Get("email")); !m {
		fmt.Println("not validated email")
		res = false
		return
	}else{
		fmt.Println("yes")
	}
	// cell phone
	if m, _ := regexp.MatchString(`^(1[3|4|5|8][0-9]\d{4,8})$`, r.Form.Get("mobile")); !m {
		fmt.Println("not validated cell phone")
		res = false
		return
	}
	// in dropdown
	slice:=[]string{"apple","pear","banana"}
	v := r.Form.Get("fruit")
	for _, item := range slice {
		if item == v {
			res = true
			break;
		}
	}
	if !res {
		fmt.Println("not in fruits")
		res = false
		return
	}

	//single select
	genders:=[]string{"1","2"}

	for _, v := range genders {
		if v == r.Form.Get("gender") {
			res = true
			break;
		}
	}
	if !res {
		fmt.Println("not qualified gender")
		res = false
		return
	}

	//multiple select
	sports:=[]string{"football","basketball","tennis"}
	// cannot use (type []string) as type []interface{} !!!!
	// https://golang.org/doc/faq#convert_slice_of_interface
	a:= tools.Slice_diff(tools.Slice_stringToInterface(r.Form["interest"]), tools.Slice_stringToInterface((sports)))
	if a == nil {
		fmt.Println("not qualified interest")
		res = false
		return
	}

	//time check?


	//XSS change html
	fmt.Println("username:", tt.HTMLEscapeString(r.Form.Get("username")))
	fmt.Println("password:", tt.HTMLEscapeString(r.Form.Get("password")))
	//template.HTMLEscape(w, []byte(r.Form.Get("username")))

	return
}

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //获取请求的方法
	if r.Method == "GET" {
		t, _ := template.ParseFiles("login.gtpl")
		log.Println(t.Execute(w, nil))
	} else {
		//请求的是登录数据，那么执行登录的逻辑判断
		fmt.Println("username:", r.Form["username"])
		fmt.Println("password:", r.Form["password"])
	}
	//token 验证
	fmt.Println("method:", r.Method) //获取请求的方法
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		t, _ := template.ParseFiles("login.gtpl")
		t.Execute(w, token)  // template
	} else {
		//请求的是登陆数据，那么执行登陆的逻辑判断
		r.ParseForm()
		token := r.Form.Get("token")
		if token != "" {
			//验证token的合法性
		} else {
			//不存在token报错
		}
		fmt.Println("username length:", len(r.Form["username"][0]))
		fmt.Println("username:", template.HTMLEscapeString(r.Form.Get("username"))) //输出到服务器端
		fmt.Println("password:", template.HTMLEscapeString(r.Form.Get("password")))
		template.HTMLEscape(w, []byte(r.Form.Get("username"))) //输出到客户端
	}
}
/**
<html>
<head>
        <title>上传文件</title>
</head>
<body>
<form enctype="multipart/form-data" action="/upload" method="post">
  <input type="file" name="uploadfile" />
  <input type="hidden" name="token" value="{{.}}"/>
  <input type="submit" value="upload" />
</form>
</body>
</html>
 */
//处理上传文件
func upload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //获取请求的方法
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		t, _ := template.ParseFiles("login.gtpl")
		t.Execute(w, token)
	} else {
		//请求的是登陆数据，那么执行登陆的逻辑判断
		r.ParseForm()
		token := r.Form.Get("token")
		if token != "" {
			//验证token的合法性
		} else {
			//不存在token报错
		}
		fmt.Println("username length:", len(r.Form["username"][0]))
		fmt.Println("username:", template.HTMLEscapeString(r.Form.Get("username"))) //输出到服务器端
		fmt.Println("password:", template.HTMLEscapeString(r.Form.Get("password")))
		template.HTMLEscape(w, []byte(r.Form.Get("username"))) //输出到客户端
	}
}

func main()  {
	mux := &MyMux{}
/*
http.ListenAndServe 方法， 这个底层其实这样处理的：初始化一个server对象，然后调用了net.Listen("tcp", addr)，也就是底层用TCP协议搭建了一个服务，然后监控我们设置的端口
 */
	http.ListenAndServe(":9090", mux)
}

/*
1 实例化Server
2 调用Server的ListenAndServe()
3 调用net.Listen("tcp", addr)监听端口
4 启动一个for循环，在循环体中Accept请求
5 对每个请求实例化一个Conn，并且开启一个goroutine为这个请求进行服务go c.serve()
6 读取每个请求的内容w, err := c.readRequest()
7 判断handler是否为空，如果没有设置handler（这个例子就没有设置handler），handler就设置为DefaultServeMux
8 调用handler的ServeHttp
9 在这个例子中，下面就进入到DefaultServeMux.ServeHttp
10 根据request选择handler，并且进入到这个handler的ServeHTTP
11 选择handler：
A 判断是否有路由能满足这个request（循环遍历ServeMux的muxEntry）
B 如果有路由满足，调用这个路由handler的ServeHTTP
C 如果没有路由满足，调用NotFoundHandler的ServeHTTP
 */
var golog = logging.MustGetLogger("example")

func NotFound404(w http.ResponseWriter, r *http.Request) {
	golog.Error("页面找不到")   //记录错误日志
	t, _ := template.ParseFiles("tmpl/404.html", nil)  //解析模板文件
	ErrorInfo := "文件找不到" //获取当前用户信息
	t.Execute(w, ErrorInfo)  //执行模板的merger操作
}

func SystemError(w http.ResponseWriter, r *http.Request) {
	golog.Critical("系统错误")   //系统错误触发了Critical，那么不仅会记录日志还会发送邮件
	t, _ := template.ParseFiles("tmpl/error.html", nil)  //解析模板文件
	ErrorInfo := "系统暂时不可用" //获取当前用户信息
	t.Execute(w, ErrorInfo)  //执行模板的merger操作
}