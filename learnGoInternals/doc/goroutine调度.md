##goroutine调度
Go的调度的实现，涉及到几个重要的数据结构。运行时库用这几个数据结构来实现goroutine的调度，管理goroutine和物理线程的运行。
这些数据结构分别是结构体G，结构体M，结构体P，以及Sched结构体。前三个的定义在文件runtime/runtime.h中，而Sched的定义在runtime/proc.c中。Go语言的调度相关实现也是在文件proc.c中

###G
G是goroutine的缩写，goroutine的控制结构，抽象
其中包括goid是这个goroutine的ID，status是这个goroutine的状态，如Gidle,Grunnable,Grunning,Gsyscall,Gwaiting,Gdead等
```go 这里的字段乱，type居然再var name之前....好奇怪
struct G
{
	uintptr	stackguard;	// 分段栈的可用空间下界
	uintptr	stackbase;	// 分段栈的栈基址
	Gobuf	sched;		//进程切换时，利用sched域来保存上下文
	uintptr	stack0;
	FuncVal*	fnstart;		// goroutine运行的函数
	void*	param;		// 用于传递参数，睡眠时其它goroutine设置param，唤醒时此goroutine可以获取
	int16	status;		// 状态Gidle,Grunnable,Grunning,Gsyscall,Gwaiting,Gdead
	int64	goid;		// goroutine的id号
	G*	schedlink;
	M*	m;		// for debuggers, but offset not hard-coded
	M*	lockedm;	// G被锁定只能在这个m上运行
	uintptr	gopc;	// 创建这个goroutine的go表达式的pc
	...
};
```
栈信息stackbase和stackguard，有运行的函数信息fnstart

PS:在当前机器的1.9中，runtime2.go中没有找到struct G只有g, 然后字段也于上面写的不同，但是上面的可以参考，比如G中含有哪些概念字段，作用，重要性,

goroutine切换时，上下文信息保存在结构体的sched域中。goroutine是轻量级的线程或者称为协程，切换时并不必陷入到操作系统内核中，所以保存过程很轻量。看一下结构体G中的Gobuf，
其实只保存了当前栈指针，程序计数器，以及goroutine自身。
```go 这里用的1.9的runtime2.go中的
type gobuf struct {
	// The offsets of sp, pc, and g are known to (hard-coded in) libmach.
	//
	// ctxt is unusual with respect to GC: it may be a
	// heap-allocated funcval so write require a write barrier,
	// but gobuf needs to be cleared from assembly. We take
	// advantage of the fact that the only path that uses a
	// non-nil ctxt is morestack. As a result, gogo is the only
	// place where it may not already be nil, so gogo uses an
	// explicit write barrier. Everywhere else that resets the
	// gobuf asserts that ctxt is already nil.
	sp   uintptr
	pc   uintptr
	g    guintptr
	ctxt unsafe.Pointer // this has to be a pointer so that gc scans it
	ret  sys.Uintreg
	lr   uintptr
	bp   uintptr // for GOEXPERIMENT=framepointer
}
```
 记录g是为了恢复当前goroutine的结构体G指针，运行时库中使用了一个常驻的寄存器extern register G* g，这个是当前goroutine的结构体G的指针
 这样做是为了快速地访问goroutine中的信息，比如，Go的栈的实现并没有使用%ebp寄存器，不过这可以通过g->stackbase快速得到
 
 ###M
 M是machine的缩写，是对机器的抽象，每个m都是对应到一条操作系统的物理线程。M必须关联了P才可以执行Go代码，但是当它处理阻塞或者系统调用中时，可以不需要关联P
 ```go from 1.9
type m struct {
	g0      *g     // goroutine with scheduling stack
	morebuf gobuf  // gobuf arg to morestack
	divmod  uint32 // div/mod denominator for arm - known to liblink

	// Fields not known to debuggers.
	procid        uint64       // for debuggers, but offset not hard-coded
	gsignal       *g           // signal-handling g
	goSigStack    gsignalStack // Go-allocated signal handling stack
	sigmask       sigset       // storage for saved signal mask
	tls           [6]uintptr   // thread-local storage (for x86 extern register)
	mstartfn      func()
	curg          *g       // current running goroutine
	caughtsig     guintptr // goroutine running during fatal signal
	p             puintptr // attached p for executing go code (nil if not executing go code)
	nextp         puintptr
	id            int32
	mallocing     int32
	throwing      int32
	preemptoff    string // if != "", keep curg running on this m
	locks         int32
	softfloat     int32
	dying         int32
	profilehz     int32
	helpgc        int32
	spinning      bool // m is out of work and is actively looking for work
	blocked       bool // m is blocked on a note
	inwb          bool // m is executing a write barrier
	newSigstack   bool // minit on C thread called sigaltstack
	printlock     int8
	incgo         bool // m is executing a cgo call
	fastrand      uint32
	ncgocall      uint64      // number of cgo calls in total
	ncgo          int32       // number of cgo calls currently in progress
	cgoCallersUse uint32      // if non-zero, cgoCallers in use temporarily
	cgoCallers    *cgoCallers // cgo traceback if crashing in cgo call
	park          note
	alllink       *m // on allm
	schedlink     muintptr
	mcache        *mcache
	lockedg       *g
	createstack   [32]uintptr // stack that created this thread.
	freglo        [16]uint32  // d[i] lsb and f[i]
	freghi        [16]uint32  // d[i] msb and f[i+16]
	fflag         uint32      // floating point compare flags
	locked        uint32      // tracking for lockosthread
	nextwaitm     uintptr     // next m waiting for lock
	needextram    bool
	traceback     uint8
	waitunlockf   unsafe.Pointer // todo go func(*g, unsafe.pointer) bool
	waitlock      unsafe.Pointer
	waittraceev   byte
	waittraceskip int
	startingtrace bool
	syscalltick   uint32
	thread        uintptr // thread handle

	// these are here because they are too large to be on the stack
	// of low-level NOSPLIT functions.
	libcall   libcall
	libcallpc uintptr // for cpu profiler
	libcallsp uintptr
	libcallg  guintptr
	syscall   libcall // stores syscall parameters on windows

	mOS
}
```
alllink把所有m放在其中, mcache缓存
结构体M中有两个G是需要关注一下的，一个是curg，代表结构体M当前绑定的结构体G。另一个是g0，是带有调度栈的goroutine，这是一个比较特殊的goroutine
普通的goroutine的栈是在堆上分配的可增长的栈，而g0的栈是M对应的线程的栈。所有调度相关的代码，会先切换到该goroutine的栈中再执行

###P
Go1.1中新加入的一个数据结构，它是Processor的缩写。结构体P的加入是为了提高Go程序的并发度，实现更好的调度。M代表OS线程。P代表Go代码执行时需要的资源。当M执行Go代码时，它需要关联一个P，
当M为idle或者在系统调用中时，它也需要P。有刚好GOMAXPROCS个P。所有的P被组织为一个数组，在P上实现了工作流窃取的调度器
```go
type p struct {
	lock mutex

	id          int32
	status      uint32 // one of pidle/prunning/...
	link        puintptr
	schedtick   uint32     // incremented on every scheduler call
	syscalltick uint32     // incremented on every system call
	sysmontick  sysmontick // last tick observed by sysmon
	m           muintptr   // back-link to associated m (nil if idle)
	mcache      *mcache
	racectx     uintptr

	deferpool    [5][]*_defer // pool of available defer structs of different sizes (see panic.go)
	deferpoolbuf [5][32]*_defer

	// Cache of goroutine ids, amortizes accesses to runtime·sched.goidgen.
	goidcache    uint64
	goidcacheend uint64

	// Queue of runnable goroutines. Accessed without lock.
	runqhead uint32
	runqtail uint32
	runq     [256]guintptr
	// runnext, if non-nil, is a runnable G that was ready'd by
	// the current G and should be run next instead of what's in
	// runq if there's time remaining in the running G's time
	// slice. It will inherit the time left in the current time
	// slice. If a set of goroutines is locked in a
	// communicate-and-wait pattern, this schedules that set as a
	// unit and eliminates the (potentially large) scheduling
	// latency that otherwise arises from adding the ready'd
	// goroutines to the end of the run queue.
	runnext guintptr

	// Available G's (status == Gdead)
	gfree    *g
	gfreecnt int32

	sudogcache []*sudog
	sudogbuf   [128]*sudog

	tracebuf traceBufPtr

	// traceSweep indicates the sweep events should be traced.
	// This is used to defer the sweep start event until a span
	// has actually been swept.
	traceSweep bool
	// traceSwept and traceReclaimed track the number of bytes
	// swept and reclaimed by sweeping in the current sweep loop.
	traceSwept, traceReclaimed uintptr

	palloc persistentAlloc // per-P to avoid mutex

	// Per-P GC state
	gcAssistTime     int64 // Nanoseconds in assistAlloc
	gcBgMarkWorker   guintptr
	gcMarkWorkerMode gcMarkWorkerMode

	// gcw is this P's GC work buffer cache. The work buffer is
	// filled by write barriers, drained by mutator assists, and
	// disposed on certain GC state transitions.
	gcw gcWork

	runSafePointFn uint32 // if 1, run sched.safePointFn at next safe point

	pad [sys.CacheLineSize]byte
}
```
结构体P中也有相应的状态：
Pidle,
Prunning,
Psyscall,
Pgcstop,
Pdead,

注意，跟G不同的是，P不存在waiting状态。MCache被移到了P中，但是在结构体M中也还保留着。在P中有一个Grunnable的goroutine队列，这是一个P的局部队列。
当P执行Go代码时，它会优先从自己的这个局部队列中取，这时可以不用加锁，提高了并发度。如果发现这个队列空了，则去其它P的队列中拿一半过来，这样实现工作流窃取的调度。这种情况下是需要给调用器加锁的

###Sched
Sched是调度实现中使用的数据结构，该结构体的定义在文件proc.c中, 1.9,应该改为schedt, 而且再runtime2中
```go
type schedt struct {
	// accessed atomically. keep at top to ensure alignment on 32-bit systems.
	goidgen  uint64
	lastpoll uint64

	lock mutex

	midle        muintptr // idle m's waiting for work
	nmidle       int32    // number of idle m's waiting for work
	nmidlelocked int32    // number of locked m's waiting for work
	mcount       int32    // number of m's that have been created
	maxmcount    int32    // maximum number of m's allowed (or die)

	ngsys uint32 // number of system goroutines; updated atomically

	pidle      puintptr // idle p's
	npidle     uint32
	nmspinning uint32 // See "Worker thread parking/unparking" comment in proc.go.

	// Global runnable queue.
	runqhead guintptr
	runqtail guintptr
	runqsize int32

	// Global cache of dead G's.
	gflock       mutex
	gfreeStack   *g
	gfreeNoStack *g
	ngfree       int32

	// Central cache of sudog structs.
	sudoglock  mutex
	sudogcache *sudog

	// Central pool of available defer structs of different sizes.
	deferlock mutex
	deferpool [5]*_defer

	gcwaiting  uint32 // gc is waiting to run
	stopwait   int32
	stopnote   note
	sysmonwait uint32
	sysmonnote note

	// safepointFn should be called on each P at the next GC
	// safepoint if p.runSafePointFn is set.
	safePointFn   func(*p)
	safePointWait int32
	safePointNote note

	profilehz int32 // cpu profiling rate

	procresizetime int64 // nanotime() of last change to gomaxprocs
	totaltime      int64 // ∫gomaxprocs dt up to procresizetime
}
```
Sched结构体只是一个壳。可以看到，其中有M的idle队列，P的idle队列，以及一个全局的就绪的G队列。Sched结构体中的Lock是非常必须的，如果M或P等做一些非局部的操作，它们一般需要先锁住调度器

###状态G,M,P
G的状态

空闲中(_Gidle): 表示G刚刚新建, 仍未初始化
待运行(_Grunnable): 表示G在运行队列中, 等待M取出并运行
运行中(_Grunning): 表示M正在运行这个G, 这时候M会拥有一个P
系统调用中(_Gsyscall): 表示M正在运行这个G发起的系统调用, 这时候M并不拥有P
等待中(_Gwaiting): 表示G在等待某些条件完成, 这时候G不在运行也不在运行队列中(可能在channel的等待队列中)
已中止(_Gdead): 表示G未被使用, 可能已执行完毕(并在freelist中等待下次复用)
栈复制中(_Gcopystack): 表示G正在获取一个新的栈空间并把原来的内容复制过去(用于防止GC扫描)
M的状态

M并没有像G和P一样的状态标记, 但可以认为一个M有以下的状态:

自旋中(spinning): M正在从运行队列获取G, 这时候M会拥有一个P
执行go代码中: M正在执行go代码, 这时候M会拥有一个P
执行原生代码中: M正在执行原生代码或者阻塞的syscall, 这时M并不拥有P
休眠中: M发现无待运行的G时会进入休眠, 并添加到空闲M链表中, 这时M并不拥有P
自旋中(spinning)这个状态非常重要, 是否需要唤醒或者创建新的M取决于当前自旋中的M的数量.

P的状态

空闲中(_Pidle): 当M发现无待运行的G时会进入休眠, 这时M拥有的P会变为空闲并加到空闲P链表中
运行中(_Prunning): 当M拥有了一个P后, 这个P的状态就会变为运行中, M运行G会使用这个P中的资源
系统调用中(_Psyscall): 当go调用原生代码, 原生代码又反过来调用go代码时, 使用的P会变为此状态
GC停止中(_Pgcstop): 当gc停止了整个世界(STW)时, P会变为此状态
已中止(_Pdead): 当P的数量在运行时改变, 且数量减少时多余的P会变为此状态


##goroutine life cycle
通过goroutine的创建，消亡，阻塞和恢复等过程，来观察Go语言的调度策略
为了避免结构体M和结构体P引入的其它干扰，这里主要将注意力集中到结构体G中，以goroutine为主线
###创建
runtime.newproc。这就是一个goroutine的出生，所有新的goroutine都是通过这个函数创建的
runtime.newproc(size, f, args)功能就是创建一个新的g，这个函数不能用分段栈，因为它假设参数的放置顺序是紧接着函数f的（见前面函数调用协议一章，有关go关键字调用时的内存布局）。
分段栈会破坏这个布局，所以在代码中加入了标记#pragma textflag 7表示不使用分段栈。它会调用函数newproc1，在newproc1中可以使用分段栈。真正的工作是调用newproc1完成的。newproc1进行下面这些动作
首先，它会检查当前结构体M中的P中，是否有可用的结构体G。如果有，则直接从中取一个，否则，需要分配一个新的结构体G。如果分配了新的G，需要将它挂到runtime的相关队列中。
获取了结构体G之后，将调用参数保存到g的栈，将sp，pc等上下文环境保存在g的sched域，这样整个goroutine就准备好了，整个状态和一个运行中的goroutine被中断时一样，只要等分配到CPU，它就可以继续运行
然后将这个“准备好”的结构体G挂到当前M的P的队列中。这里会给予新的goroutine一次运行的机会，即：如果当前的P的数目没有到上限，也没有正在自旋抢CPU的M，则调用wakep将P立即投入运行
wakep函数唤醒P时，调度器会试着寻找一个可用的M来绑定P，必要的时候会新建M。让我们看看新建M的函数newm
```go
// 新建一个m，它将以调用fn开始，或者是从调度器开始
static void
newm(void(*fn)(void), P *p)
{
	M *mp;
	mp = runtime·allocm(p);
	mp->nextp = p;
	mp->mstartfn = fn;
	runtime·newosproc(mp, (byte*)mp->g0->stackbase);
}
```
runtime.newm功能跟newproc相似,前者分配一个goroutine,而后者分配一个M。其实一个M就是一个操作系统线程的抽象，可以看到它会调用runtime.newosproc。
看到了从Go的运行时库到操作系统的接口，runtime.newosproc(平台相关的)会调用系统的runtime.clone(平台相关的)来新建一个线程，新的线程将以runtime.mstart为入口函数。runtime.newosproc是个很有意思的函数，
还有一些信号处理方面的细节，但是对鉴于我们是专注于调度方面，就不对它进行更细致的分析了，感兴趣的读者可以自行去runtime/os_linux.c看看源代码。runtime.clone是用汇编实现的,代码在sys_linux_amd64.s
既然线程是以runtime.mstart为入口的，那么接下来看mstart函数
mstart是runtime.newosproc新建的系统线程的入口地址，新线程执行时会从这里开始运行。新线程的执行和goroutine的执行是两个概念，由于有m这一层对机器的抽象，是m在执行g而不是线程在执行g。
所以线程的入口是mstart，g的执行要到schedule才算入口。函数mstart最后调用了schedule
如果是从mstart进入到schedule的，那么schedule中逻辑非常简单，大概就这几步
```text
找到一个等待运行的g
如果g是锁定到某个M的，则让那个M运行
否则，调用execute函数让g在当前的M中运行
```
execute会恢复newproc1中设置的上下文，这样就跳转到新的goroutine去执行了。从newproc出生一直到运行的过程分析，到此结束!
newproc -> newproc1 -> (如果P数目没到上限)wakep -> startm -> (可能引发)newm -> newosproc -> (线程入口)mstart -> schedule -> execute -> goroutine运行

__PS: 再1.9的代码里首先runtime没有proc函数，而是再proc.go中看到（应该也是创建g的，比如获取空闲g,然后配置newg,sp处理等, 还会配置p,貌似racegostart方法启动g,然后放到对列中 runqput__
__后面也有wakep,再wakeup中触发startm->newm->newosproc->clone等__

###进出系统调用
假设goroutine"生病"了，它要进入系统调用了，暂时无法继续执行。进入系统调用时，如果系统调用是阻塞的，goroutine会被剥夺CPU，将状态设置成Gsyscall后放到就绪队列
Go的syscall库中提供了对系统调用的封装，它会在真正执行系统调用之前先调用函数.entersyscall，并在系统调用函数返回后调用.exitsyscall函数。这两个函数就是通知Go的运行时库这个goroutine进入了系统调用或者完成了系统调用，调度器会做相应的调度
比如syscall包中的Open函数，它会调用Syscall(SYS_OPEN, uintptr(unsafe.Pointer(_p0)), uintptr(mode), uintptr(perm))实现。这个函数是用汇编写的，在syscall/asm_linux_amd64.s中可以看到它的定义
这个函数是用汇编写的，在syscall/asm_linux_amd64.s中可以看到它的定义
```text
TEXT	·Syscall(SB),7,$0
	CALL	runtime·entersyscall(SB)
	MOVQ	16(SP), DI
	MOVQ	24(SP), SI
	MOVQ	32(SP), DX
	MOVQ	$0, R10
	MOVQ	$0, R8
	MOVQ	$0, R9
	MOVQ	8(SP), AX	// syscall entry
	SYSCALL
	CMPQ	AX, $0xfffffffffffff001
	JLS	ok
	MOVQ	$-1, 40(SP)	// r1
	MOVQ	$0, 48(SP)	// r2
	NEGQ	AX
	MOVQ	AX, 56(SP)  // errno
	CALL	runtime·exitsyscall(SB)
	RET
ok:
	MOVQ	AX, 40(SP)	// r1
	MOVQ	DX, 48(SP)	// r2
	MOVQ	$0, 56(SP)	// errno
	CALL	runtime·exitsyscall(SB)
	RET
```
可以看到它进系统调用和出系统调用时分别调用了runtime.entersyscall和runtime.exitsyscall函数, 其中处理了什么？
1. 将函数的调用者的SP,PC等保存到结构体G的sched域中。同时，也保存到g->gcsp和g->gcpc等，这个是跟垃圾回收相关的。
2. 检查结构体Sched中的sysmonwait域，如果不为0，则将它置为0，并调用runtime·notewakeup(&runtime·sched.sysmonnote)。做这这一步的原因是，目前这个goroutine要进入Gsyscall状态了，它将要让出CPU。如果有人在等待CPU的话，会通知并唤醒等待者，马上就有CPU可用了
3. 将m的MCache置为空，并将m->p->m置为空，表示进入系统调用后结构体M是不需要MCache的，并且P也被剥离了，将P的状态设置为PSyscall
有一个与entersyscall函数稍微不同的函数叫entersyscallblock，它会告诉提示这个系统调用是会阻塞的，因此会有一点点区别。它调用的releasep和handoffp
releasep将P和M完全分离，使p->m为空，m->p也为空，剥离m->mcache，并将P的状态设置为Pidle。注意这里的区别，在非阻塞的系统调用entersyscall中只是设置成Psyscall，并且也没有将m->p置为空
handoffp切换P。将P从处于syscall或者locked的M中，切换出来交给其它M。每个P中是挂了一个可执行的G的队列的，如果这个队列不为空，即如果P中还有G需要执行，则调用startm让P与某个M绑定后立刻去执行，否则将P挂到idlep队列中
出系统调用时会调用到runtime·exitsyscall，这个函数跟进系统调用做相反的操作。它会先检查当前m的P和它状态，如果P不空且状态为Psyscall，则说明是从一个非阻塞的系统调用中返回的，这时是仍然有CPU可用的。
因此将p->m设置为当前m，将p的mcache放回到m，恢复g的状态为Grunning。否则，它是从一个阻塞的系统调用中返回的，因此之前m的P已经完全被剥离了。这时会查看调用中是否还有idle的P，如果有，则将它与当前的M绑定
4. 如果从一个阻塞的系统调用中出来，并且出来的这一时刻又没有idle的P了，要怎么办呢？这种情况代码当前的goroutine无法继续运行了，调度器会将它的状态设置为Grunnable，将它挂到全局的就绪G队列中，然后停止当前m并调用schedule函数

###goroutine消亡
goroutine的消亡比较简单，注意在函数newproc1，设置了fnstart为goroutine执行的函数，而将新建的goroutine的sched域的pc设置为了函数runtime.exit
当fnstart函数执行完返回时，它会返回到runtime.exit中。这时Go就知道这个goroutine要结束了，runtime.exit中会做一些回收工作，会将g的状态设置为Gdead等，并将g挂到P的free队列中

从以上的分析中，其实已经基本上经历了goroutine的各种状态变化。在newproc1中新建的goroutine被设置为Grunnable状态，投入运行时设置成Grunning。
在entersyscall的时候goroutine的状态被设置为Gsyscall，到出系统调用时根据它是从阻塞系统调用中出来还是非阻塞系统调用中出来，又会被设置成Grunning或者Grunnable的状态。
在goroutine最终退出的runtime.exit函数中，goroutine被设置为Gdead状态  
goroutine的状态变迁图：
![](images/5.2.goroutine_state.jpg?raw=true)


