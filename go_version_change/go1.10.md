[原文](https://golang.org/doc/go1.10)
1.9之后6个月发布，主要的改变在toolchain,runtime的实现和lib.
具体的包括:
* improve [caching of built packages](https://golang.org/doc/go1.10#build)
* add [caching of successful test results](https://golang.org/doc/go1.10#test)
* runs [vet automatically during tests](https://golang.org/doc/go1.10#test-vet)
* permits [passing string values directly between Go and C using cgo](https://golang.org/doc/go1.10#cgo)
* add new compiler option whitelist may caused [invalid falg error](https://golang.org/s/invalidflag)

##语法
没有太多change
method expression允许any type as a receiver,这个编译器早就实现，现在在语法上放开
比如： 