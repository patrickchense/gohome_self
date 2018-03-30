interface是Go语言中最成功的设计之一，空的interface可以被当作“鸭子”类型使用，它使得Go这样的静态语言拥有了一定的动态性，但却又不损失静态语言在类型安全方面拥有的编译时检查的优势
Go中的interface在底层是如何实现的呢？
###Eface和Iface
1.9 runtime2.go
```go
type iface struct {
	tab  *itab
	data unsafe.Pointer
}

type eface struct {
	_type *_type
	data  unsafe.Pointer
}
```
interface实际上就是一个结构体，包含两个成员。其中一个成员是指向具体数据的指针，另一个成员中包含了类型信息。空接口和带方法的接口略有不同，下面分别是空接口和带方法的接口是使用的数据结构
先看Eface，它是interface{}底层使用的数据结构。数据域中包含了一个void*指针，和一个类型结构体的指针。interface{}扮演的角色跟C语言中的void*是差不多的，Go中的任何对象都可以表示为interface{}。不同之处在于，interface{}中有类型信息，于是可以实现反射
类型信息的结构体定义如下：
```go
type _type struct {
	size       uintptr
	ptrdata    uintptr // size of memory prefix holding all pointers
	hash       uint32
	tflag      tflag
	align      uint8
	fieldalign uint8
	kind       uint8
	alg        *typeAlg
	// gcdata stores the GC type data for the garbage collector.
	// If the KindGCProg bit is set in kind, gcdata is a GC program.
	// Otherwise it is a ptrmask bitmap. See mbitmap.go for details.
	gcdata    *byte
	str       nameOff
	ptrToThis typeOff
}
```
其实在前面我们已经见过它了。精确的垃圾回收中，就是依赖Type结构体中的gc域的。不同类型数据的类型信息结构体并不完全一致，
Type是类型信息结构体中公共的部分，其中size描述类型的大小，hash数据的hash值，align是对齐，fieldAlgin是这个数据嵌入结构体时的对齐，kind是一个枚举值，每种类型对应了一个编号
alg是一个函数指针的数组，存储了hash/equal/print/copy四个函数操作。UncommonType是指向一个函数指针的数组，收集了这个类型的实现的所有方法
在reflect包中有个KindOf函数，返回一个interface{}的Type，其实该函数就是简单的取Eface中的Type域
Iface和Eface略有不同，它是带方法的interface底层使用的数据结构。data域同样是指向原始数据的，而Itab的结构如下：
```go
type itab struct {
	inter  *interfacetype
	_type  *_type
	link   *itab
	hash   uint32 // copy of _type.hash. Used for type switches.
	bad    bool   // type does not implement interface
	inhash bool   // has this itab been added to hash?
	unused [2]byte
	fun    [1]uintptr // variable sized
}
```

###具体类型向接口类型赋值
将具体类型数据赋值给interface{}这样的抽象类型，中间会涉及到类型转换操作。从接口类型转换为具体类型(也就是反射)，也涉及到了类型转换。这个转换过程中做了哪些操作呢？
1. 先看将具体类型转换为接口类型。如果是转换成空接口，这个过程比较简单，就是返回一个Eface，将Eface中的data指针指向原型数据，type指针会指向数据的Type结构体
将某个类型数据转换为带方法的接口时，会复杂一些。中间涉及了一道检测，该类型必须要实现了接口中声明的所有方法才可以进行转换
```go
type I interface {
	String()
}
var a int = 5
var b I = a
```
编译会报错：
```text
cannot use a (type int) as type I in assignment:
	int does not implement I (missing String method)
```
说明具体类型转换为带方法的接口类型是在编译过程中进行检测的
那么这个检测是如何实现的呢？在runtime下找到了iface.c文件，应该是早期版本是在运行时检测留下的，其中有一个itab函数就是判断某个类型是否实现了某个接口，如果是则返回一个Itab结构体
类型转换时的检测就是比较具体类型的方法表和接口类型的方法表，看具体类型是实现了接口类型所声明的所有的方法。还记得Type结构体中是有个UncommonType字段的，里面有张方法表，类型所实现的方法都在里面。
而在Itab中有个InterfaceType字段，这个字段中也有一张方法表，就是这个接口所要求的方法。这两处方法表都是排序过的，只需要一遍顺序扫描进行比较，应该可以知道Type中否实现了接口中声明的所有方法。
最后还会将Type方法表中的函数指针，拷贝到Itab的fun字段中

Type的UncommonType中有一个方法表，某个具体类型实现的所有方法都会被收集到这张表中。reflect包中的Method和MethodByName方法都是通过查询这张表实现的。表中的每一项是一个Method，其数据结构如下
```go
type method struct {
	name nameOff
	mtyp typeOff
	ifn  textOff
	tfn  textOff
}
```
Iface的Itab的InterfaceType中也有一张方法表，这张方法表中是接口所声明的方法。其中每一项是一个IMethod，数据结构如下：
```go
type imethod struct {
	name nameOff
	ityp typeOff
}
```
跟上面的Method结构体对比可以发现，这里是只有声明没有实现的
Iface中的Itab的func域也是一张方法表，这张表中的每一项就是一个函数指针，也就是只有实现没有声明
类型转换时的检测就是看Type中的方法表是否包含了InterfaceType的方法表中的所有方法，并把Type方法表中的实现部分拷到Itab的func那张表中

###reflect
reflect就是给定一个接口类型的数据，得到它的具体类型的类型信息，它的Value等。reflect包中的TypeOf和ValueOf函数分别做这个事情
还有像
```go
v, ok = i.(T)
```
这样的语法，也是判断一个接口i的具体类型是否为类型T，如果是则将其值返回给v。这跟上面的类型转换一样，也会检测转换是否合法。不过这里的检测是在运行时执行的。在runtime下的iface.c文件中，
有一系统的assetX2X函数，比如runtime.assetE2T，runtime.assetI2T等等。这个实现起来比较简单，只需要比较Iface中的Itab的type是否与给定Type为同一个




