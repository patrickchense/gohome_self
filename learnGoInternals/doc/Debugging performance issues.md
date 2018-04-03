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
选项 --list=functionName 可以查看某个func的耗时
解决cpu func问题的一些hits:
* runtime.mallocgc方法耗时长，证明程序分配了很多小的空间，看看哪里产生的小对象
* channel操作，sync.Mutext,其他同步操作耗时长，系统资源竞争激烈，考虑重构，处理共享资源，通常sharding/partitioning,local buffering
/batching, copy-on-write等手段
* syscall.Read/Write耗时长,程序可能有很多小的文件读写，尝试用bufio包的缓存
* GC时间长，可能的问题是分配对象太多，heap设置过小等

##Memory
显示方法对于heap的占用情况
使用: go test --memprofile , net/http/pprof http://myserver:6060/debug/pprof/heap 查看，使用 runtime/pprof.WriteHeapProfile
可选option:
* --memprofilerate, 设置收集速率，默认1 每512k, 
* --inuse_space, 默认,live的allocate
* --alloc_space, 程序启动时的分配
* --inuse/alloc_space and --inuse/alloc_objects, 显示分配的bytes或者obj数量
PS： 对于persistent和transient的对象，如果是大对象persistent，程序开始就分配了，那么肯定可以再profile中发现，但是也不会影响程序执行，而transient的大对象，
可能很难在profile中发现(默认的 inuse_space)，但是恰恰会影响系统（GC），所以想关注程序执行的内存占用情况，需要用inuse_space,而如果想关注执行速度，需要--alloc_objects
* --functions 粒度调整为函数级别(granularity) 这也是默认的granularity, --lines, --addresses, --files
优化建议:
* combine objects into larger objects, 比如 使用bytes.Buffer而不是*bytes.Buffer（使用预申请大小的buffer而不是bytes.Buffer自动增长）
* 注意escape analysis，注意local var的作用域，注意函数调用时候的参数复制
* 预定义slice大小
* 使用更小的类型，比如int8而不是int
* 对象尽量不要包含指针，指针对象不会被GC
* 使用缓存，重复分配临时对象

##Blocking
查看goroutine的block点
使用方式:go test --blockprofile, net/http/pprof http://myserver:6060:/debug/pprof/block, 调用runtime/pprof.Lookup("block").WriteTo
可选option:
* 启用，必须设置runtime.SetBlockProfileRate, =1 就是每个blocking event
* --lines, 监控到行
不是所有的Blocking都是有问题的，比如sync.WaitGroup的block通常是没有问题的，sync.Conn就不一定了，Consumer如果block在channel，意味着server慢
或者缺乏worker，Producer如果block在consumer的channel，就一般ok。
Block在sync.Mutex和sync.RWMutex通常是不好的
* --ignore, pprof可以排除不感兴趣的block
常用的优化建议：
* 使用buffered的channel来实现producer-consumer模型
* 使用sync.RWMutex而不是sync.Mutex，RWMutex不会阻塞其他reader
* 可以使用copy-on-write来替换Mutex
copy-on-write例子:
```go
type Config struct {
    Routes   map[string]net.Addr
    Backends []net.Addr
}
 
var config unsafe.Pointer  // actual type is *Config
 
// Worker goroutines use this function to obtain the current config.
func CurrentConfig() *Config {
    return (*Config)(atomic.LoadPointer(&config))
}
 
// Background goroutine periodically creates a new Config object
// as sets it as current using this function.
func UpdateConfig(cfg *Config) {
    atomic.StorePointer(&config, unsafe.Pointer(cfg))
}
```
* 分区，减少共享数据blocking的手段，例子:
```go
type Partition struct {
    sync.RWMutex
    m map[string]string
}
 
const partCount = 64
var m [partCount]Partition
 
func Find(k string) string {
    idx := hash(k) % partCount
    part := &m[idx]
    part.RLock()
    v := part.m[k]
    part.RUnlock()
    return v
}
```
* 本地缓存或批处理来减少不可分区数据的blocking,例子:
```go
const CacheSize = 16
 
type Cache struct {
    buf [CacheSize]int
    pos int
}
 
func Send(c chan [CacheSize]int, cache *Cache, value int) {
    cache.buf[cache.pos] = value
    cache.pos++
    if cache.pos == CacheSize {
        c <- cache.buf
        cache.pos = 0
    }
}
```
* 使用sync.Pool

##goroutine
显示所有的运行时goroutine，用于debug load balance问题或者死锁
使用方式，不支持go test，net/http/pprof 然后 http://myserver:6060:/debug/pprof/goroutine 查看，svg/pdf包调用runtime/pprof.Lookup("goroutine").WriteTo

##tracer
###GC tracer
使用方式：运行时配置环境变量 GODEBUG=gctrace=1
显示:
```text
gc9(2): 12+1+744+8 us, 2 -> 10 MB, 108615 (593983-485368) objects, 4825/3620/0 sweeps, 0(0) handoff, 6(91) steal, 16/1/0 yields
gc10(2): 12+6769+767+3 us, 1 -> 1 MB, 4222 (593983-589761) objects, 4825/0/1898 sweeps, 0(0) handoff, 6(93) steal, 16/10/2 yields
gc11(2): 799+3+2050+3 us, 1 -> 69 MB, 831819 (1484009-652190) objects, 4825/691/0 sweeps, 0(0) handoff, 5(105) steal, 16/1/0 yields
```
gc9：GC的id，第几次
gc9(2): 2是参与的worker thread
12+1+744+8 us： STW的时间
2->10MB, heap从2M到10M，GC之后
108615 (593983-485368) objects， 所有的对象数量，free的数量，之后的数量
4825/3620/0 sweeps： 描述sweep阶段，一共4825spans，3620swept了，0swept在STW时期
0(0) handoff, 6(91) : 并行标记阶段，LB的情况，0对象handoff,6次steal操作（91个对象），
16/1/0 yields： 并行标记的效率，17次yeild操作等待其他线程
GC: 是mark-and-sweep， 耗时公式：
```text
Tgc = Tseq + Tmark + Tsweep
```
Tseq：是停止goroutine和其他preparation操作的时间
Tmark：是STW标记时间，是影响latency的重要部分
Tsweep：是并行清扫时间
预估mark时间公式：
```text
Tmark = C1*Nlive + C2*MEMlive_ptr + C3*Nlive_ptr
```
Nlive:GC时期live的对象数量
MEMlive_ptr：live对象的指针占用内存
Nlive_ptr：live对象指针数量
预估sweep时间公式：
```text
Tsweep = C4*MEMtotal + C5*MEMgarbage
```
MEMtotal：heap总大小
MEMgarbage：heap中garbage大小

GOGC 环境变量，触发GC之后，下次触发需要heap的大小调整，默认100
比如，heap当前4M，GC之后，下次GC发生在heap为8M的时候，线性的增长
sweep是依赖heap大小的，设置GOGC大的话，可以减少GC次数而latency不变
GOMAXPROCS 设置的总核数，可以影响GC的worker，而默认GC8个线程

###memory tracer
使用方式：运行使用GODEBUG=allocfreetrace=1环境变量
显示：
```text
tracealloc(0xc208062500, 0x100, array of parse.Node)
goroutine 16 [running]:
runtime.mallocgc(0x100, 0x3eb7c1, 0x0)
    runtime/malloc.goc:190 +0x145 fp=0xc2080b39f8
runtime.growslice(0x31f840, 0xc208060700, 0x8, 0x8, 0x1, 0x0, 0x0, 0x0)
    runtime/slice.goc:76 +0xbb fp=0xc2080b3a90
text/template/parse.(*Tree).parse(0xc2080820e0, 0xc208023620, 0x0, 0x0)
    text/template/parse/parse.go:289 +0x549 fp=0xc2080b3c50
...
 
tracefree(0xc208002d80, 0x120)
goroutine 16 [running]:
runtime.MSpan_Sweep(0x73b080)
       runtime/mgc0.c:1880 +0x514 fp=0xc20804b8f0
runtime.MCentral_CacheSpan(0x69c858)
       runtime/mcentral.c:48 +0x2b5 fp=0xc20804b920
runtime.MCache_Refill(0x737000, 0xc200000012)
       runtime/mcache.c:78 +0x119 fp=0xc20804b950
...
```
显示memroy的address，block，size，type，goroutine id和stack trace

###scheduler trace
可以深入goroutine的调度细节，调试LB和伸缩性的问题
GODEBUG=schedtrace=1000 环境变量，1000ms显示一次
显示：
```text
SCHED 1004ms: gomaxprocs=4 idleprocs=0 threads=11 idlethreads=4 runqueue=8 [0 1 0 3]
SCHED 2005ms: gomaxprocs=4 idleprocs=0 threads=11 idlethreads=5 runqueue=6 [1 5 4 0]
SCHED 3008ms: gomaxprocs=4 idleprocs=0 threads=11 idlethreads=4 runqueue=10 [2 2 2 1]
```
"1004ms"：从程序开始到现在的时间
Gomaxprocs： 当前的GOMAXPROCS
Idleprocs：idling处理器
Threads：调度器创建的worker数量，thread有3个state，execute Go code (gomaxprocs-idleprocs), execute syscalls/cgocalls or idle
Runqueue：全局运行goroutine队列的长度
[0 1 0 3]：每个processor的运行goroutine数量

###tracer 总结
3个tracer集成使用
GODEBUG=gctrace=1,allocfreetrace=1,schedtrace=1000

##memory statistics
通过runtime.ReadMemStats调用，或者net/http/pprof, http://myserver:6060/debug/pprof/heap?debug=1
HeapAlloc: 当前heap大小
HeapSys：heap总大小
HeapObjects：heap中对象数量
HeapReleased：heap释放的内存，5分钟没有使用
Sys：OS分配的内存大小
Sys-HeapReleased：程序的有效内存
StackSys：goroutine的stack大小
MSpanSys/MCacheSys/BuckHashSys/GCSys/OtherSys： 运行时内存数据
PauseNs：上次GC STW时间

##heap dumper
使用runtime/debug.WriteHeapDump
```go
f, err := os.Create("heapdump")
if err != nil { ... }
debug.WriteHeapDump(f.Fd())
```
然后可以把文件转成hprof格式或者dot文件
go get github.com/randall77/hprof/dumptodot
dumptodot heapdump mybinary > heap.dot
使用Graphviz打开dot

hprof：
go get github.com/randall77/hprof/dumptohprof
dumptohprof heapdump heap.hprof
jhat heap.hprof
使用浏览器http://localhost:7000查看






