package main

import (
	"fmt"

	"github.com/SRG98/automatas-go/controllers"
)

func main() {
	controller := controllers.NewController()
	err := controller.Run()
	if err != nil {
		fmt.Println("Error:", err)
	}
}
