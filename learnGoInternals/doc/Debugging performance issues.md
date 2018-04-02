[原文](https://software.intel.com/en-us/blogs/2014/05/10/debugging-performance-issues-in-go-programs)

##简介
需要debugGo程序来找到问题，定位问题，比如找hotspot(CPU,MEMORY,IO,etc),SQL问题，代码问题(可以更快更简单)

下面会介绍一些工具来达到目的

##CPU
Go runtime内嵌了CPU profiler, 显示CPU占用时间比例，有3个方法
1. 最简单,go test的-cpuprofile flag
```text
$ go test -run=none -bench=ClientServerParallel4 -cpuprofile=cprof net/http
``` 
会把benchmark和产生的CPU profile存放到cprof文件
然后可以通过
```text
$ go tool pprof --text http.test cprof
```
会打印CPU的hottest functions
其他可选的option包括 --web , --list 
唯一的缺点是，这些只支持test

2. net/http/pprof 包，对于网络应用最理想的方式，只需要import包，然后collect就可以了
```text
$ go tool pprof --text mybin http://myserver:6060:/debug/pprof/profile

```
3. 人工profile收集. import runtime/pprof，然后添加代码到main
```go
if *flagCpuprofile != "" {
    f, err := os.Create(*flagCpuprofile)
    if err != nil {
        log.Fatal(err)
    }
    pprof.StartCPUProfile(f)
    defer pprof.StopCPUProfile()
}
```
可以写到特定的文件，用同样的go tool pprof 可以查看
pprof的可视化例子![](images/cpu_profile.png)







