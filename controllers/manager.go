package controllers

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/SRG98/automatas-go/models"
)

func (c *Controller) readJSONFile(filename string) (*models.Automata, error) {
	var automaton models.Automata

	fileBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(fileBytes, &automaton)
	if err != nil {
		return nil, err
	}

	return &automaton, nil
}

func (c *Controller) writeJSONFile(filename string, automata *models.Automata) bool {
	fileBytes, err := json.MarshalIndent(automata, "", "  ")
	if err != nil {
		return false
	}

	err = os.WriteFile(filename, fileBytes, 0644)
	if err != nil {
		return false
	}

	return true
}

func (c *Controller) readInputFile(filepath string) error {
	fileBytes, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(fileBytes), "\n")
	c.inputStrings = lines

	fmt.Print("Texto con ruta predefinida cargado!")
	return nil
}
