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
