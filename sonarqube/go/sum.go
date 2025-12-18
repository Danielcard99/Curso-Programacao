package main

import "fmt"

func main() {
	fmt.Println(sum(2, 2))
	fmt.Println(sub(2, 2))
	fmt.Println(div(2, 2))
	fmt.Println(mod(2, 2))
	fmt.Println(pow(2, 2))
	fmt.Println(sqrt(2))
	fmt.Println(cube(2))
	fmt.Println(square(2))
}

func sum(a int, b int) int {
	return a + b
}

func sub(a int, b int) int {
	return a - b
}

func div(a int, b int) int {
	return a / b
}

func mod(a int, b int) int {
	return a % b
}

func pow(a int, b int) int {
	return a ^ b
}

func sqrt(a int) int {
	return a * a
}

func cube(a int) int {
	return a * a * a
}

func square(a int) int {
	return a * a
}
