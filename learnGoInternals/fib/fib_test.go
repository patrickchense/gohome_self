package fib

import "testing"

func BenchmarkFib(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Fib(20)
	}
}

func BechmarkComplicated(b *testing.B)  {
	for n := 0; n < b.N; n++ {
		b.StopTimer()
		//complicatedSetup()
		b.StartTimer()
		//function under test
	}
}

func BenchmarkExpensive(b *testing.B)  {
	//boringAndExpensiveSetup()
	b.ResetTimer()
	for n:= 0; n < b.N; n++ {
		// function under test
	}
}

func BenchmarkRead(b *testing.B)  {
	b.ReportAllocs()
	for n:= 0; n < b.N; n++ {
		//function under test
	}
}