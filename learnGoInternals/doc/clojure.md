##闭包
###GO中的闭包
```go
func f(i int) func() int {
	return func() int {
		i++
		return i
	}
}
```
这里f的返回，就是一个函数，就是闭包,函数本身没有定义变量i,引用了环境（函数f）中的变量i
所以可以说
__闭包=函数+引用环境__
闭包的环境中引用的变量不能够在栈上分配

###escape analyze
Go的一个语言特性
```go
func f() *Cursor {
	var c Cursor
	c.X = 500
	noinline()
	return &c
}
```
Cursor是一个结构体，这种写法在C语言中是不允许的，因为变量c是在栈上分配的，当函数f返回后c的空间就失效了。但是，在Go语言规范中有说明，这种写法在Go语言中合法的。
语言会自动地识别出这种情况并在堆上分配c的内存，而不是函数f的栈上

查看汇编:
```compile
MOVQ	$type."".Cursor+0(SB),(SP)	// 取变量c的类型，也就是Cursor
PCDATA	$0,$16
PCDATA	$1,$0
CALL	,runtime.new(SB)	// 调用new函数，相当于new(Cursor)
PCDATA	$0,$-1
MOVQ	8(SP),AX	// 取c.X的地址放到AX寄存器
MOVQ	$500,(AX)	// 将AX存放的内存地址的值赋为500
MOVQ	AX,"".~r0+24(FP)
ADDQ	$16,SP
```
识别出变量需要在堆上分配，是由编译器的一种叫escape analyze的技术实现的。如果输入命令
go build --gcflags=-m main.go
