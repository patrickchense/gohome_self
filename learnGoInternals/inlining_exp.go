package main

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func F() {
	const a,b = 100, 20
	if Max(a,b) == b {
		panic(b)
	}
}

/*
-bash-3.2$ go build -gcflags=-m inlining_exp.go
# command-line-arguments
./inlining_exp.go:3:6: can inline Max
./inlining_exp.go:12:8: inlining call to Max
# command-line-arguments
runtime.main_main·f: relocation target main.main not defined
runtime.main_main·f: undefined: "main.main"
 */

 /*
 看优化后的汇编： go build -gcflags="-m -S" inlining_exp.go  2>&1 | less
 > inlining_exp_huibian.txt
  */