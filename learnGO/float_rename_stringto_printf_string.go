package main

import (
	"fmt"
	"strconv"
)

type Celsius float64

//定义float64, 可以定制String, 直接println输出
func (c Celsius) String() string {
	return "The temperature is: " + strconv.FormatFloat(float64(c), 'f', 1, 32) + " °C"

}
func main() {
	var c Celsius = 18.36
	fmt.Println(c)
}
