package controllers

import(
	"fmt"
)

func main() {
	s1 :=  State{
		Data      "A",
		IsInitial true,
		IsFinal   false,
		Adjacent  [],
	}

	fmt.Println(s1)
}