[原文](https://golang.org/doc/effective_go.html)
##pre
golang官网的effective go阅读，记录一些模糊的，不知道的

###init func
```go
func init() {
    if user == "" {
        log.Fatal("$USER not set")
    }
    if home == "" {
        home = "/home/" + user
    }
    if gopath == "" {
        gopath = home + "/go"
    }
    // gopath may be overridden by --gopath flag on command line.
    flag.StringVar(&gopath, "gopath", gopath, "override default GOPATH")
}
```
init方法在包中所有变量初始化，import包初始化之后
init方法不能像声明这样暴露

###comment
comment是可以用来godoc的，那么写comment的一些规范:
* 提供package comment
* 需要格式, 比如*
* 开头以要comment的内容,比如comment方法Compile, 那么 //Compile 开头， 原因:根据某些内容搜索的时候，能知道需要的comment来自哪里
* 第一句最好是能一句总结

###name
不要设置太长的名字，而通过comment(godoc)来说明  

###string loop
string通常不用range编历, 比如:
```go
for pos, char := range "日本\x80語" { // \x80 is an illegal UTF-8 encoding
    fmt.Printf("character %#U starts at byte position %d\n", char, pos)
}
```
输出:
```text
character U+65E5 '日' starts at byte position 0
character U+672C '本' starts at byte position 3
character U+FFFD '�' starts at byte position 6
character U+8A9E '語' starts at byte position 7
```
用rune是每个UTF-8字符

###type switch
类型断言
```go
var t interface{}
t = functionOfSomeType()
switch t := t.(type) {
default:
    fmt.Printf("unexpected type %T\n", t)     // %T prints whatever type t has
case bool:
    fmt.Printf("boolean %t\n", t)             // t has type bool
case int:
    fmt.Printf("integer %d\n", t)             // t has type int
case *bool:
    fmt.Printf("pointer to boolean %t\n", *t) // t has type *bool
case *int:
    fmt.Printf("pointer to integer %d\n", *t) // t has type *int
}
```

###defer
多个defer类似队列FILO  
组合使用:
```go
func trace(s string) string {
    fmt.Println("entering:", s)
    return s
}

func un(s string) {
    fmt.Println("leaving:", s)
}

func a() {
    defer un(trace("a"))
    fmt.Println("in a")
}

func b() {
    defer un(trace("b"))
    fmt.Println("in b")
    a()
}

func main() {
    b()
}
```
输出
```text
entering: b
in b
entering: a
in a
leaving: a
leaving: b
```

###allocation
####new
build-in func allocation memory, 但是零初始  
new(T), 以零值初始化类型T，然后返回地址（*T),所以有零值的类型可以用new来创建  
new(T) == &T{}  

###make
make(T, args)  
只用来创建slice,channel,map, 返回非零值的T(不是*T), 表明这些类型的对象必须初始化之后才能使用  
比如slice, 有一个pointer(指向数组),一个len, 一个cap,那么必须初始化这几个值之后才可使用，否则为nil, make([]int, 10, 100)  

区别:
```go
var p *[]int = new([]int)       // allocates slice structure; *p == nil; rarely useful
var v  []int = make([]int, 100) // the slice v now refers to a new array of 100 ints

// Unnecessarily complex:
var p *[]int = new([]int)
*p = make([]int, 100, 100)

// Idiomatic:
v := make([]int, 100)
```

###arrays
Go和C中array的区别, Go中:
* array是值，Assigning one array to another copies all the elements
* 传递array到func，是复制array而不是pointer
* The size of an array is part of its type. The types [10]int and [20]int are distinct




