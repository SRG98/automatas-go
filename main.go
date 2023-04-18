package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/SRG98/automatas-go/models"
)

func main() {
	filename := "data/auto.json"

	// Crea un nuevo autómata
	auto := models.NewAutomaton()
	auto.SetName("Sample Automaton")
	auto.SetAlphabet([]string{"a", "b"})
	// Agrega estados
	auto.NewState("q0", false, true)
	auto.NewState("q1", true, false)

	// Agrega transiciones
	auto.NewTransition("q0", "q1", []string{"a", "b"})
	auto.NewTransition("q1", "q0", []string{"b"})

	fmt.Println(auto.ToString())

	// Writing
	err := writeJSONFile(filename, auto)
	if err != nil {
		fmt.Println("Error writing JSON file:", err)
		os.Exit(1)
	}

	// Reading
	automaton, err := readJSONFile(filename)
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		os.Exit(1)
	}

	fmt.Printf("automata: %+v\n", automaton)

}

func readJSONFile(filename string) (models.Automata, error) {
	var automaton models.Automata

	fileBytes, err := os.ReadFile(filename)
	if err != nil {
		return automaton, err
	}

	err = json.Unmarshal(fileBytes, &automaton)
	if err != nil {
		return automaton, err
	}

	return automaton, nil
}

func writeJSONFile(filename string, automata *models.Automata) error {
	fileBytes, err := json.MarshalIndent(automata, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, fileBytes, 0644)
	if err != nil {
		return err
	}

	return nil

}
