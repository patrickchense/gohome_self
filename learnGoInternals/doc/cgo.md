下面是一个使用cgo的例子：
```go
package rand

/*
#include <stdlib.h>
*/
import "C"

func Random() int {
    return int(C.random())
}

func Seed(i int) {
    C.srandom(C.uint(i))
}
```
参考(../cgo_exp.go)
导入了"C"，但是在Go的标准库中并没有一个"C"包。这是因为"C"是一个伪包，这是一个特殊的名字，cgo通过这个包知道它是引用C命名空间的。Go编译器使用符号"·"来区分命名空间
而C编译器使用不同的约定，因此使用C包中的名字时，Go编译器就知道应该使用C的命名约定


##问题
* Go使用的是分段栈，初始栈大小很小，当发现栈不够时会动态增长。动态增长是通过进入函数时插入检测指令实现的。然而C函数不使用分段栈技术，并且假设栈是足够大的。那么Go是如何处理不让cgo调用发生栈溢出的呢？
* Go中的goroutine都是协作式的，运行到调用runtime库时就有机会进行调度。然而C函数是不会与Go的runtime做这种交互的，所以cgo的函数不是一个协作式的，那么如何避免进入C函数的这个goroutine“失控”？
* cgo不仅仅是从Go调用C，还包括从C中调用Go函数。这里面又有哪些技术难点？举个简单的例子，C中调用Go函数f，而f中是使用了go建立新的goroutine的，但是在C中是不支持Go的runtime的。

##预备
###m的g0栈
问题一的解答:
Go的运行时库中使用了几个重要的结构体，其中M是机器的抽象。每个M可以运行各个goroutine，在结构体M的定义中有一个相对特殊的goroutine叫g0(还有另一个比较特殊的gsignal，与本节内容无关暂且不讲)。那么这个g0特殊在什么地方呢
g0的特殊之处在于它是带有调度栈的goroutine，下文就将其称为“m的g0栈“。Go在执行调度相关代码时，都是使用的m的g0栈。当一个g执行的是调度相关的代码时，它并不是直接在自己的栈中执行，而是先切换到m的g0栈然后再执行代码
m的g0栈是一个特殊的栈，g0的分配和普通goroutine的分配过程不同，g0是在m建立时就生成的，并且给它分配的栈空间比较大，可以假定它的大小是足够大而不必使用分段栈。而普通的goroutine是在runtime.newproc时建立，并且初始栈空间分配得很小(4K)，会在需要时增长。不仅如此，m的g0栈同时也是这个m对应的物理线程的栈。
这样就相当于拥有了一个“无穷”大小的非分段栈，于是回答了前面提的那个问题：Go使用的是分段栈，初始栈大小很小，当发现栈不够时会动态增长。动态增长是通过进入函数时插入检测指令实现的。然而C函数不使用分段栈技术，并且假设栈是足够大的。调用cgo代码时，使用的是m的g0栈，这是一个足够大的不会发生分段的栈。

函数newm是新那一个结构体M，其中调用runtime.allocm分配M的空间。它的g0域是这个分配的
```go
mp->g0 = runtime·malg(8192);
```
__为什么是8k？够用吗？__

###进入系统调用
Go的运行时库对系统调用作了特殊处理，所有涉及到调用系统调用之前，都会先调用runtime.entersyscall，而在出系统调用函数之后，会调用runtime.exitsyscall。这样做原因跟调度器相关，目的是始终维持GOMAXPROCS的数量，
当进入到系统调用时，runtime.entersyscall会将P的M剥离并将它设置为PSyscall状态，告知系统此时其它的P有机会运行，以保证始终是GOMAXPROCS个P在运行

runtime.entersyscall函数会立刻返回，它仅仅是起到一个通知的作用。那么这跟cgo又有什么关系呢？这个关系可大着呢！在执行cgo函数调用之前，其实系统会先调用runtime.entersyscall。这是一个很关键的处理，
Go把cgo的C函数调用像系统调用一样独立出去了，不让它影响运行时库。
这就回答了前面提出的第二个问题：Go中的goroutine都是协作式的，运行到调用runtime库时就有机会进行调度。然而C函数是不会与Go的runtime做这种交互的，所以cgo的函数不是一个协作式的，
那么如何避免进入C函数的这个goroutine“失控”？答案就在这里。将C函数像处理系统调用一样隔离开来，这个goroutine也就不必参与调度了。而其它部分的goroutine正常的运行不受影响

###退出系统调用
退出系统调用跟进入系统调用是一个相反的过程，runtime.exitsyscall函数会查看当前仍然有可用的P，则让它继续运行，否则这个goroutine就要被挂起了。
对于cgo的代码也是同样的作用，出了cgo的C函数调用之后会调用runtime.exitsyscall。


##cgo关键技术
整个cgo的实现依赖于几个部分，依赖于cgo命令生成桩文件，依赖于6c和6g对Go这一端的代码进行编译，依赖gcc对C那一端编译成动态链接库，同时，还依赖于运行时库实现Go和C互操作的一些支持
cgo命令会生成一些桩文件，这些桩文件是给6c和6g命令使用的，它们是Go和C调用之间的桥梁。原始的C文件会使用gcc编译成动态链接库的形式使用

###cgo命令
```C
import "C"
```
字段，就会先调用cgo命令。cgo提取出相应的C函数接口部分，生成桩文件
对它执行cgo命令：
```shell
go tool cgo cgo_exp.go
```
在当前目录下会生成一个_obj的文件夹，文件夹里会包含下列文件：  
[root@localhost _obj]# ls  
cgo_exp.cgo1.go  cgo_exp.cgo2.c  _cgo_export.c  _cgo_export.h  _cgo_flags  _cgo_gotypes.go  _cgo_main.c  _cgo_.o

###桩文件
cgo生成了很多文件，其中大多数作用都是包装现有的函数，或者进行声明








