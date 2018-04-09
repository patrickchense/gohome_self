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
compares incoming requests against a list of predefined URL paths,and calls the associated handler for the path whenever a match is found
* Handlers
    * responsible for writing response headers and bodies
    * Almost any type ("object") can be a handler, so long as it satisfies the http.Handler interface
    * In lay terms, that simply means it must have a ServeHTTP method with the following signature:ServeHTTP(http.ResponseWriter, *http.Request)

##basic_server
最简单的httpserver
```go
http.HandleFunc("/", someFunc)
http.ListenAndServe(":8080", nil) //nil 默认的router

func someFunc(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world!"))
}
```
##Mux
```go
myMux := http.NewServeMux()
myMux.HandleFunc("/", someFunc)
http.ListenAndServe(":8080", myMux)
```

##serve file
```go
func (this *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    path := r.URL.Path[1:]
    log.Println(path)
//    path := "templates" + r.URL.Path

    data, err := ioutil.ReadFile(string(path))

    if err == nil {
        w.Write(data)
    } else {
        w.WriteHeader(404)
        w.Write([]byte("404 My Friend - " + http.StatusText(404)))
    }
}
func main() {
    http.Handle("/", new(MyHandler))
    http.ListenAndServe(":8080", nil)
}
```
假设启动server,然后运行文件: localhost:8080/templates/home.html,就会发现
控制台打印:
```text
2018/04/09 10:34:29 templates/home.html
2018/04/09 10:34:29 public/css/css_reset.css
2018/04/09 10:34:29 public/css/app.css
2018/04/09 10:34:29 public/css/flyout_button.css
2018/04/09 10:34:29 public/css/flyout_menu.css
2018/04/09 10:34:29 public/scripts/flyout_menu.js
2018/04/09 10:34:30 public/img/favicon.ico
```
读取了html中每个关联的url:
```javascript
    <link rel='stylesheet' href='../public/css/css_reset.css'>
    <link rel='stylesheet' href='http://fonts.googleapis.com/css?family=Droid+Sans:400'>
    <link rel='stylesheet' href='../public/css/app.css'>
    <link rel='stylesheet' href='../public/css/flyout_button.css'>
    <link rel='stylesheet' href='../public/css/flyout_menu.css'>
    <link rel='shortcut icon' href='../public/img/favicon.ico'>
```
但是没有设置需要的context-type

##set context type
替换上面的w.Write(data)为:
```go
if err == nil {
        var contentType string

        if strings.HasSuffix(path, ".css") {
            contentType = "text/css"
        } else if strings.HasSuffix(path, ".html") {
            contentType = "text/html"
        } else if strings.HasSuffix(path, ".js") {
            contentType = "application/javascript"
        } else if strings.HasSuffix(path, ".png") {
            contentType = "image/png"
        } else if strings.HasSuffix(path, ".svg") {
            contentType = "image/svg+xml"
        } else {
            contentType = "text/plain"
        }

        w.Header().Add("Content Type", contentType)
        w.Write(data)
    } else {
        w.WriteHeader(404)
        w.Write([]byte("404 Mi amigo - " + http.StatusText(404)))
    }

```
##buffer
读取文件，可以使用bufio包来buffer处理
```go
 f, err := os.Open(path)

    if err == nil {
        bufferedReader := bufio.NewReader(f)

        var contentType string

        if strings.HasSuffix(path, ".css") {
            contentType = "text/css"
        } else if strings.HasSuffix(path, ".html") {
            contentType = "text/html"
        } else if strings.HasSuffix(path, ".js") {
            contentType = "application/javascript"
        } else if strings.HasSuffix(path, ".png") {
            contentType = "image/png"
        } else if strings.HasSuffix(path, ".svg") {
            contentType = "image/svg+xml"
        } else if strings.HasSuffix(path, ".mp4") {
            contentType = "video/mp4"
        } else {
            contentType = "text/plain"
        }

        w.Header().Add("Content Type", contentType)
        bufferedReader.WriteTo(w)
    }
```

##template
3步走:
* [template.NeW](http://golang.org/pkg/text/template/#New)
    create the template "object"
* [template.Parse](http://golang.org/pkg/text/template/#Template.Parse)
    put your template into the template "object"
* [template.Execute](http://golang.org/pkg/text/template/#Template.Execute)
    merge your template with data
    
```go
func main() {
    http.HandleFunc("/", myHandlerFunc)
    http.ListenAndServe(":8080", nil)
    // nil means use default ServeMux
}

func myHandlerFunc(w http.ResponseWriter, req *http.Request) {
    w.Header().Add("Content Type", "text/html")
    tmpl, err := template.New("anyNameForTemplate").Parse(doc)
    if err == nil {
        tmpl.Execute(w, nil)
        // nil means no data to pass in
    }
}
```

传递data:
```go
func main() {
    http.HandleFunc("/", myHandlerFunc)
    http.ListenAndServe(":8080", nil)
}

func myHandlerFunc(w http.ResponseWriter, req *http.Request) {
    w.Header().Add("Content Type", "text/html")
    tmpl, err := template.New("anyNameForTemplate").Parse(doc)
    if err == nil {
        tmpl.Execute(w, req.URL.Path)
//        tmpl.Execute(w, req.URL.Path[1:])
    }
}

const doc = `
<!DOCTYPE html>
<html>
<head lang="en">
    <meta charset="UTF-8">
    <title>First Template</title>
</head>
<body>
    <h1>Hello {{.}}</h1>
</body>
</html>
`
```
例子2：
```go
const doc = `
<!DOCTYPE html>
<html>
<head lang="en">
    <meta charset="UTF-8">
    <title>First Template</title>
</head>
<body>
    <h1>My name is {{.FirstName}}</h1>
    <p>{{.Message}}</p>
</body>
</html>
`

func toddFunc(w http.ResponseWriter, req *http.Request) {
    w.Header().Add("Content Type", "text/html")
    tmpl, err := template.New("anyNameForTemplate").Parse(doc)
    if err == nil {
        context := Context{"Todd", "more Go, please"}
        tmpl.Execute(w, context)
    }
}

func main() {
    http.HandleFunc("/todd", toddFunc)
    http.HandleFunc("/ming", mingFunc)
    http.HandleFunc("/rio", rioFunc)
    http.HandleFunc("/", jamesFunc)
    http.ListenAndServe(":8080", nil)
}
func mingFunc(w http.ResponseWriter, req *http.Request) {
    w.Header().Add("Content Type", "text/html")
    tmpl, err := template.New("anyNameForTemplate").Parse(doc)
    if err == nil {
        context := Context{"Ming", "I am a problem solver!"}
        tmpl.Execute(w, context)
    }
}

func rioFunc(w http.ResponseWriter, req *http.Request) {
    w.Header().Add("Content Type", "text/html")
    tmpl, err := template.New("anyNameForTemplate").Parse(doc)
    if err == nil {
        context := Context{"Rio", "I drank the google-aid"}
        tmpl.Execute(w, context)
    }
}

func jamesFunc(w http.ResponseWriter, req *http.Request) {
    w.Header().Add("Content Type", "text/html")
    tmpl, err := template.New("anyNameForTemplate").Parse(doc)
    if err == nil {
        context := Context{"James", "Another beer, please"}
        tmpl.Execute(w, context)
    }
}

```
例子3,template conditional :
```go
type Context struct {
    FirstName string
    Message string
    URL string
}

func main() {
    http.HandleFunc("/", myHandlerFunc)
    http.ListenAndServe(":8080", nil)
}

func myHandlerFunc(w http.ResponseWriter, req *http.Request) {
    w.Header().Add("Content Type", "text/html")
    tmpl, err := template.New("anyNameForTemplate").Parse(doc)
    if err == nil {
        context := Context{"Todd", "more beer, please", req.URL.Path}
        tmpl.Execute(w, context)
    }
}

const doc = `
<!DOCTYPE html>
<html>
<head lang="en">
    <meta charset="UTF-8">
    <title>First Template</title>
</head>
<body>
    {{if eq .URL "/nobeer"}}
        <h1>We're out of beer. Sorry!</h1>
    {{else}}
        <h1>Yes, grab another beer, {{.FirstName}}</h1>
    {{end}}

    <hr>

    <h2>Here's all the data:</h2>
    <p>{{.}}</p>
</body>
</html>
`
/*
conditionals
if
if / else
if / else if

testing
eq - equal
---- an unlimited number of conditions can be tested against the first term
------ eq 1 (0+1) (2-1)
------ if they all evaluate to be the same, the test evaluates to TRUE
------ operator is listed first, followed by operands
ne - not equal
lt - less than
---- first arg compared to second
---- first condition less than second - evals to true
gt - greater than
le - less than or equal to
ge - greater than or equal to
*/
```
第四个例子,loop:
```go
type Context struct {
    FirstName string
    Message string
    URL string
    Beers []string
    Title string

}

func main() {
    http.HandleFunc("/", myHandlerFunc)
    http.ListenAndServe(":8080", nil)
}

func myHandlerFunc(w http.ResponseWriter, req *http.Request) {
    w.Header().Add("Content Type", "text/html")
    tmpl, err := template.New("anyNameForTemplate").Parse(doc)
    if err == nil {
        context := Context{
            "Todd",
            "more beer, please",
            req.URL.Path,
            []string{"New Belgium", "La Fin Du Monde", "The Alchemist"},
            "Favorite Beers",
        }
        tmpl.Execute(w, context)
    }
}

const doc = `
<!DOCTYPE html>
<html>
<head lang="en">
    <meta charset="UTF-8">
    <title>{{.Title}}</title>
</head>
<body>

    <h1>{{.FirstName}} says, "{{.Message}}"</h1>

    {{if eq .URL "/nobeer"}}
        <h2>We're out of beer, {{.FirstName}}. Sorry!</h2>
    {{else}}
        <h2>Yes, grab another beer, {{.FirstName}}</h2>
        <ul>
            {{range .Beers}}
            <li>{{.}}</li>
            {{end}}
        </ul>
    {{end}}

    <hr>

    <h2>Here's all the data:</h2>
    <p>{{.}}</p>
</body>
</html>
`

/*
range
allows you to loop over data with many items
array, slice, map, channel
when you "range" loop over data
the pipeline {{.}} gets set to the current item in the data
another way to say this: "the range operator resets the pipeline
to be the individual item in the collection"
range / else
same as range
however, if the data is len == 0, then
the else block gets executed
(eg, empty shopping cart)

sub-templates
include templates in templates
a view can include many different templates
-- call this, call that, call another thing

*/
```

第五个例子, multiple template:
```go
type Context struct {
    FirstName string
    Message string
    URL string
    Beers []string
    Title string

}

func main() {
    http.HandleFunc("/", myHandlerFunc)
    http.ListenAndServe(":8080", nil)
}

func myHandlerFunc(w http.ResponseWriter, req *http.Request) {
    w.Header().Add("Content Type", "text/html")
    templates := template.New("template")
    templates.New("test").Parse(doc)
    templates.New("header").Parse(head)
    templates.New("footer").Parse(foot)
    context := Context{
        "Todd",
        "more beer, please",
        req.URL.Path,
        []string{"New Belgium", "La Fin Du Monde", "The Alchemist"},
        "Favorite Beers",
    }
    templates.Lookup("test").Execute(w, context)
}

const doc = `
{{template "header" .Title}}
<body>

    <h1>{{.FirstName}} says, "{{.Message}}"</h1>

    {{if eq .URL "/nobeer"}}
        <h2>We're out of beer, {{.FirstName}}. Sorry!</h2>
    {{else}}
        <h2>Yes, grab another beer, {{.FirstName}}</h2>
        <ul>
            {{range .Beers}}
            <li>{{.}}</li>
            {{end}}
        </ul>
    {{end}}

    <hr>

    <h2>Here's all the data:</h2>
    <p>{{.}}</p>
</body>
{{template "footer"}}
`

const head = `
<!DOCTYPE html>
<html>
<head lang="en">
    <meta charset="UTF-8">
    <title>{{.}}</title>
</head>
`

const foot = `
</html>
`

/*
create a template that contains all of your templates
any sub-templates invoked have to be either
-- siblings
-- descendents
of the parent template

{{template "header"}}
"header" is the name we gave the template with template.New

func (*Template) Lookup
func (t *Template) Lookup(name string) *Template
Lookup returns the template with the given name that is associated with t,
or nil if there is no such template.

*/
```
##防止注入
使用html/template
```go
import (
    "net/http"
//    "text/template"
    "html/template"
)

//var Message string = "more beer, please sir"
//var Message string = "alert('you have been pwned')"
var Message string = "<script>alert('you have been pwned, BIATCH')</script>"

func main() {
    http.HandleFunc("/", myHandlerFunc)
    http.ListenAndServe(":8080", nil)
}

func myHandlerFunc(w http.ResponseWriter, req *http.Request) {
    w.Header().Add("Content Type", "text/html")
    tmpl, err := template.New("anyNameForTemplate").Parse(doc)
    if err == nil {
        tmpl.Execute(w, Message)
    }
}

const doc = `
<!DOCTYPE html>
<html>
<head lang="en">
    <meta charset="UTF-8">
    <title>Injection Safe</title>
</head>
<body>

    <p>{{.}}</p>

    <script>{{.}}</script>

</body>
</html>
`

/*
html/template
Package template (html/template) implements data-driven templates for generating
HTML output safe against code injection. It provides the same interface as package
text/template and should be used instead of text/template whenever the output is HTML.

HTML templates treat data values as plain text which should be encoded so they can be
safely embedded in an HTML document. The escaping is contextual, so actions can appear
within JavaScript, CSS, and URI contexts.

http://golang.org/pkg/html/template/

to run the above code ...
try the different Message variables with text/template import
... then ...
try the different Message variables with html/template import


*/
```

##mvc
* MODEL
business logic & rules
data storage

* VIEW
what the client sees

* CONTROLLER
the glue between model & view
coordinates the model & view layers
determines how the model needs to be interacted with to meet a user's request
passes the results of the model layers work to the view layer
responsibilities:
    * generate output and send it back to client
        * templates
        * bind data
    * receive user actions
        * ajax
        * forms

steps:
(1) create a template "cache"
* one template to hold other templates
* all of the "held" templates will be siblings of each other
    * this means the "held" templates can call/include each other
    * templates can call/include sibling & descendent templates

例子1，最简单的找到对应的template:
```go
func main() {
    templates := populateTemplates()

    http.HandleFunc("/",
    func(w http.ResponseWriter, req *http.Request) {
        requestedFile := req.URL.Path[1:]
        template := templates.Lookup(requestedFile + ".html")

        if template != nil {
            template.Execute(w, nil)
        } else {
            w.WriteHeader(404)
        }
    })


    http.ListenAndServe(":8080", nil)
}

func populateTemplates() *template.Template {
    result := template.New("templates")

    basePath := "templates"
    templateFolder, _ := os.Open(basePath)
    defer templateFolder.Close()

    templatePathsRaw, _ := templateFolder.Readdir(-1)
    // -1 means all of the contents
    templatePaths := new([]string)
    for _, pathInfo := range templatePathsRaw {
        if !pathInfo.IsDir() {
            *templatePaths = append(*templatePaths,
            basePath + "/" + pathInfo.Name())
        }
    }

    result.ParseFiles(*templatePaths...)

    return result
}
```
例子2，根据不同的content-type,不同的handler:
```go
func main() {
    templates := populateTemplates()

    http.HandleFunc("/",
    func(w http.ResponseWriter, req *http.Request) {
        requestedFile := req.URL.Path[1:]
        template := templates.Lookup(requestedFile + ".html")

        if template != nil {
            template.Execute(w, nil)
        } else {
            w.WriteHeader(404)
        }
    })
    http.HandleFunc("/img/", serveResource)
    http.HandleFunc("/css/", serveResource)
    http.HandleFunc("/scripts/", serveResource)
    http.ListenAndServe(":8080", nil)
}

func serveResource(w http.ResponseWriter, req *http.Request) {
    path := "public" + req.URL.Path
    var contentType string
    if strings.HasSuffix(path, ".css") {
        contentType = "text/css"
    } else if strings.HasSuffix(path, ".png") {
        contentType = "image/png"
    } else if strings.HasSuffix(path, ".jpg") {
        contentType = "image/jpg"
    } else if strings.HasSuffix(path, ".svg") {
        contentType = "image/svg+xml"
    } else if strings.HasSuffix(path, ".js") {
        contentType = "application/javascript"
    } else {
        contentType = "text/plain"
    }

    log.Println(path)
    log.Println(contentType)

    f, err := os.Open(path)

    if err == nil {
        defer f.Close()
        w.Header().Add("Content Type", contentType)
        br := bufio.NewReader(f)
        br.WriteTo(w)
    } else {
        w.WriteHeader(404)
    }
}

func populateTemplates() *template.Template {
    result := template.New("templates")

    basePath := "templates"
    templateFolder, _ := os.Open(basePath)
    defer templateFolder.Close()

    templatePathsRaw, _ := templateFolder.Readdir(-1)
    // -1 means all of the contents
    templatePaths := new([]string)
    for _, pathInfo := range templatePathsRaw {
        log.Println(pathInfo.Name())
        if !pathInfo.IsDir() {
            *templatePaths = append(*templatePaths,
            basePath + "/" + pathInfo.Name())
        }
    }

    result.ParseFiles(*templatePaths...)

    return result
}
```
例子3,subtemplate:
```go
func main() {
    templates := populateTemplates()

    http.HandleFunc("/",
    func(w http.ResponseWriter, req *http.Request) {
        requestedFile := req.URL.Path[1:]
        template :=
        templates.Lookup(requestedFile + ".html")

        var context interface{} = nil
        switch requestedFile {
            case "home":
            context = viewmodels.GetHome()
        }
        if template != nil {
            template.Execute(w, context)
        } else {
            w.WriteHeader(404)
        }
    })

    http.HandleFunc("/img/", serveResource)
    http.HandleFunc("/css/", serveResource)
    http.HandleFunc("/scripts/", serveResource)
    http.ListenAndServe(":8080", nil)
}

func serveResource(w http.ResponseWriter, req *http.Request) {
    path := "public" + req.URL.Path
    var contentType string
    if strings.HasSuffix(path, ".css") {
        contentType = "text/css"
    } else if strings.HasSuffix(path, ".png") {
        contentType = "image/png"
    } else if strings.HasSuffix(path, ".jpg") {
        contentType = "image/jpg"
    } else if strings.HasSuffix(path, ".svg") {
        contentType = "image/svg+xml"
    } else if strings.HasSuffix(path, ".js") {
        contentType = "application/javascript"
    } else {
        contentType = "text/plain"
    }

    log.Println(path)
    log.Println(contentType)

    f, err := os.Open(path)

    if err == nil {
        defer f.Close()
        w.Header().Add("Content Type", contentType)
        br := bufio.NewReader(f)
        br.WriteTo(w)
    } else {
        w.WriteHeader(404)
    }
}

func populateTemplates() *template.Template {
    result := template.New("templates")

    basePath := "templates"
    templateFolder, _ := os.Open(basePath)
    defer templateFolder.Close()

    templatePathsRaw, _ := templateFolder.Readdir(-1)
    templatePaths := new([]string)
    for _, pathInfo := range templatePathsRaw {
        log.Println(pathInfo.Name())
        if !pathInfo.IsDir() {
            *templatePaths = append(*templatePaths,
            basePath + "/" + pathInfo.Name())
        }
    }

    result.ParseFiles(*templatePaths...)

    return result
}
/*
we're going to break our main html page down into different parts
--- header
--- content1
--- content2
--- footer

This will help with
-- code reusability
-- organizing our data and keeping it clean

we are going to separate the data that is used in the VIEW layer
from the rest of the data that the application uses
-- good practice as the needs of the VIEW and MODEL layer differ over time

create:
viewmodels / home.go

add this to main.go imports:
"viewmodels"
*/
```

例子4： 传递data:
```go
func main() {
    templates := populateTemplates()

    http.HandleFunc("/",
    func(w http.ResponseWriter, req *http.Request) {
        requestedFile := req.URL.Path[1:]
        template :=
        templates.Lookup(requestedFile + ".html")

        var context interface{} = nil
        switch requestedFile {
            case "home":
            context = viewmodels.GetHome()
            case "search":
            context = viewmodels.GetSearch()
        }
        if template != nil {
            template.Execute(w, context)
        } else {
            w.WriteHeader(404)
        }
    })

    http.HandleFunc("/img/", serveResource)
    http.HandleFunc("/css/", serveResource)
    http.HandleFunc("/scripts/", serveResource)
    http.ListenAndServe(":8080", nil)
}

func serveResource(w http.ResponseWriter, req *http.Request) {
    path := "public" + req.URL.Path
    var contentType string
    if strings.HasSuffix(path, ".css") {
        contentType = "text/css"
    } else if strings.HasSuffix(path, ".png") {
        contentType = "image/png"
    } else if strings.HasSuffix(path, ".jpg") {
        contentType = "image/jpg"
    } else if strings.HasSuffix(path, ".svg") {
        contentType = "image/svg+xml"
    } else if strings.HasSuffix(path, ".js") {
        contentType = "application/javascript"
    } else {
        contentType = "text/plain"
    }

    log.Println(path)
    log.Println(contentType)

    f, err := os.Open(path)

    if err == nil {
        defer f.Close()
        w.Header().Add("Content Type", contentType)
        br := bufio.NewReader(f)
        br.WriteTo(w)
    } else {
        w.WriteHeader(404)
    }
}

func populateTemplates() *template.Template {
    result := template.New("templates")

    basePath := "templates"
    templateFolder, _ := os.Open(basePath)
    defer templateFolder.Close()

    templatePathsRaw, _ := templateFolder.Readdir(-1)
    templatePaths := new([]string)
    for _, pathInfo := range templatePathsRaw {
        log.Println(pathInfo.Name())
        if !pathInfo.IsDir() {
            *templatePaths = append(*templatePaths,
            basePath + "/" + pathInfo.Name())
        }
    }

    result.ParseFiles(*templatePaths...)

    return result
}

```

