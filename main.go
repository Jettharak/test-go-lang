package main

import "fmt"

func main() {
	var num1, num2 int

	fmt.Print("กรอกจำนวนเต็ม 2 จำนวน: ")
	if _, err := fmt.Scan(&num1, &num2); err != nil {
		fmt.Println("กรุณากรอกจำนวนเต็มให้ถูกต้อง:", err)
		return
	}

	result := add(num1, num2)
	fmt.Println("ผลลัพธ์:", result)
}

func add(a, b int) int {
	return a + b
}
