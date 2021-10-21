package main

import "fmt"

func main() {
	arr := []string{"gaze", "at", "the", "stars"}

	for range arr {
		arr = append(arr, "tonight")
	}
	fmt.Println(arr)
}
