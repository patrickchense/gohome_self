[视频地址](https://www.youtube.com/watch?v=Vlie-srOU8c)
一个web相关的教程
提供了一组资源:
[GO语法例子&book](https://github.com/GoesToEleven/csuf)
[julieschmidt httprouter](https://github.com/julienschmidt/httprouter)
[golang-book](http://www.golang-book.com/)

##简介
Go web   
* configuration: APP Engine & golang & Google Cloud
* RESOURCE SERVER
    + listen on a tcp port
    + handle requests: route a URL to a file
* ServeMux = HTTP request router = multiplexor

##basic_server
最简单的httpserver
```go
http.HandleFunc("/", someFunc)
http.ListenAndServe(":8080", nil) //nil 默认的router

func someFunc(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world!"))
}
```

