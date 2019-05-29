package main

import "fmt"

func main () {
	a := map[string]string{}
	for i := 0; i < 5; i ++ {
		a[string(i)] = string(i +1)
	}
	fmt.Println(a)
}