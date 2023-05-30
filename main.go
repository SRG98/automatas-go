package main

import (
	"fmt"

	"github.com/SRG98/automatas-go/controllers"
	//"github.com/SRG98/automatas-go/views"
)

func main() {
	controller := controllers.NewController()
	err := controller.Run()
	if err != nil {
		fmt.Println("Error:", err)
	}
}

/*
func main() {
	controller := controllers.NewController()
	uinterface := views.NewUI(controller)
	err := uinterface.RunUI()
	if err != nil {
		fmt.Println("Error:", err)
	}
}
*/
