package main

import (
	"fmt"
	"github.com/astaxie/go-i18n"
)
/*
这个i18n 有问题， 是不能用的，但是code 的结构可以参考
 */
func main() {
	tr, _ := i18n.NewIL("config/locals", "zh")
	fmt.Println(tr.Translate("submit"))
}
