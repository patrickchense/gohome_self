[原文](https://dave.cheney.net/2013/10/12/how-to-use-conditional-compilation-with-the-go-build-tool)
##总结
conditional compilation  
大概讲述了在build的时候如何选择compile的平台等，通过flag或者file suffix等方式

##golist
go list gives you access to the internal data structures which power the build process  
go list的arguments和大部分go build，test和install 相似，但是不产生编译
使用 -f 可以提供一小部分text/template代码（执行在go/build.Package架构下）
```text
% go list -f '{{.GoFiles}}' os/exec
[exec.go lp_unix.go]
```

##Build tags
第一种声明conditional compilation的方式是通过在source code中annotation，被称为build tag  
build tag的实现是在尽量靠近文件顶端的地方添加comment
go build 命令在编译时会去寻找build tag  
build tag遵循：
* 空格分割==OR
* 逗号分割==AND
* 每个term是文字和数字组成，!取反

例子:
```go
// +build darwin freebsd netbsd openbsd
```
表示文件只在BSD系统编译，支持kqueue
可以使用多个build tag， OR
```go
// +build linux darwin
// +build 386
```
可以在linux/386或darwin/386

###注意
build tag 要跟package声明至少隔一行
```go
// +build !linux
package mypkg // wrong
```
正确的:
```go
// +build !linux

package mypkg // correct
```
使用go vet:
```text
% go vet mypkg
mypkg.go:1: +build comment appears too late in file
exit status 1
```

一个常用的格式:
```text
% head headspin.go 
// Copyright 2013 Way out enterprises. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build someos someotheros thirdos,!amd64

// Package headspin implements calculates numbers so large
// they will make your head spin.
package headspin
```

##File suffix
第二种conditional compilation的方式是文件后缀  
例子:
```text
mypkg_freebsd_arm.go // only builds on freebsd/arm systems
mypkg_plan9.go       // only builds on plan9
```
不能只有file suffix  

##两者选择
通常使用file suffix来实现目的,但是如果多个平台或者要排除某个平台用build tag  

##使用+build来debug test
[原文](https://dave.cheney.net/2014/09/28/using-build-to-switch-between-debug-and-release)
在go test -intergration的时候，很多test很难找到出错的  
```text
% go test -integration -v
=== RUN TestUnmarshalAttrs
--- PASS: TestUnmarshalAttrs (0.00s)
=== RUN TestNewClient
--- PASS: TestNewClient (0.00s)
=== RUN TestClientLstat
Unknown message 0
```
用debug的方式:
```go
// +build debug

package sftp

import "log"

func debug(fmt string, args ...interface{}) {
	log.Printf(fmt, args...)
}
```
发布版：
```go
// +build !debug

package sftp

func debug(fmt string, args ...interface{}) {}
```
再测试:
```text
% go test -tags debug -integration -v -run=Lstat
2014/09/28 11:18:31 send packet sftp.sshFxInitPacket, len: 38
=== RUN TestClientLstat
2014/09/28 11:18:31 send packet sftp.sshFxInitPacket, len: 5
2014/09/28 11:18:31 send packet sftp.sshFxpLstatPacket, len: 62
Unknown message 0
```







