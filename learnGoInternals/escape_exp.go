package main

func Sum() int {
	const count = 100
	numbers := make([]int, 100)
	for i := range numbers {
		numbers[i] = i + 1
	}

	var sum int
	for _, i := range numbers {
		sum += i
	}
	return sum
}


/*
逃逸分析例子， 这里运行
-bash-3.2$ go build -gcflags=-m escape_exp.go
# command-line-arguments
./escape_exp.go:5:17: Sum make([]int, 100) does not escape
# command-line-arguments
runtime.main_main·f: relocation target main.main not defined
runtime.main_main·f: undefined: "main.main"


 */