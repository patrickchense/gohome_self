package main

import (
	"errors"
	"fmt"
	"log"
	"reflect"
)

// functional programming with go

func main() {
	println(sqrPlusOne(1))
	res, err := sqrPlusOnePipe(1)
	if err != nil {
		log.Fatal(err.Error())
	}
	println(res)
}

//simple pipe, because func only return one result
func sqrPlusOne(x int) int {
	sqr := func(x int) int { return x * x }
	inc := func(x int) int { return x + 1 }
	return inc(sqr(x))
}

// with err, can't work
func sqrPlusOneErr(x int) (int, error) {
	sqr := func(x int) (int, error) {
		if x < 0 {
			return 0, errors.New("x should not be negative")
		}
		return x * x, nil
	}
	inc := func(x int) int { return x + 1 }
	y, err := sqr(x)
	if err != nil {
		return 0, err
	}
	return inc(y), nil
}

func sqrPlusOnePipe(x int) (int, error) {
	var result int
	err := Pipe(
		func(x int) (int, error) {
			if x < 0 {
				return 0, errors.New("x should not be negative")
			}
			return x * x, nil
		},
		func(x int) int { return x + 1 },
		func(x int) { result = x }, // the sink
	)(x) // the execution of pipeline
	if err != nil {
		return 0, err
	}
	return result, nil
}

// errType is the type of error interface.
var errType = reflect.TypeOf((*error)(nil)).Elem()

// Pipeline is the func type for the pipeline result.
type Pipeline func(...interface{}) error

func empty(...interface{}) error { return nil }

func Pipe(fs ...interface{}) Pipeline {
	if len(fs) == 0 {
		return empty
	}

	return func(args ...interface{}) (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("pipeline panicked: %v", r)
			}
		}()

		var inputs []reflect.Value
		for _, arg := range args {
			inputs = append(inputs, reflect.ValueOf(arg))
		}

		for fIndex, f := range fs {
			outputs := reflect.ValueOf(f).Call(inputs)
			inputs = inputs[:0]

			funcType := reflect.TypeOf(f)

			for oIndex, output := range outputs {
				if funcType.Out(oIndex).Implements(errType) {
					if !output.IsNil() {
						err = fmt.Errorf("%s func failed: %w", ord(fIndex), output.Interface().(error))
						return
					}
				} else {
					inputs = append(inputs, output)
				}
			}
		}

		return
	}
}

func ord(index int) string {
	order := index + 1
	switch {
	case order > 10 && order < 20:
		return fmt.Sprintf("%dth", order)
	case order%10 == 1:
		return fmt.Sprintf("%dst", order)
	case order%10 == 2:
		return fmt.Sprintf("%dnd", order)
	case order%10 == 3:
		return fmt.Sprintf("%drd", order)
	default:
		return fmt.Sprintf("%dth", order)
	}

}
