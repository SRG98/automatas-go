package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/OverCV/go-automats/models"
	"github.com/awalterschulze/gographviz"
)

func main() {
	projectPath, e := filepath.Abs(".")
	if e != nil {
		fmt.Println("Error getting project path:", e)
		os.Exit(1)
	}

	filename := "data/auto.json"

	// Crea un nuevo autómata
	auto := models.NewAutomaton()
	auto.SetName("Sample Automaton")
	auto.SetAlphabet([]string{"a", "b"})

	// Agrega estados
	auto.NewState("q0", false, true)
	auto.NewState("q1", true, false)

	// Agrega transiciones
	auto.NewTransition("q1", "q0", []string{"b"})
	auto.NewTransition("q0", "q1", []string{"a", "b"})

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

	graph, err := automataToGraphviz(&automaton)
	if err != nil {
		fmt.Println("Error generating Graphviz graph:", err)
		os.Exit(1)
	}

	imagePath := filepath.Join(projectPath, "data", "graph.png")
	err = saveGraphAsPNG(graph, imagePath)
	if err != nil {
		fmt.Println("Error saving Graphviz graph as PNG:", err)
		os.Exit(1)
	}
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

func automataToGraphviz(automaton *models.Automata) (*gographviz.Graph, error) {
	graphAst, err := gographviz.ParseString("digraph G {}")
	if err != nil {
		return nil, err
	}

	graph := gographviz.NewGraph()
	if err := gographviz.Analyse(graphAst, graph); err != nil {
		return nil, err
	}

	for _, state := range automaton.States {
		attrs := make(map[string]string)
		attrs["label"] = state.Data
		if state.IsInitial {
			attrs["style"] = "filled"
			attrs["fillcolor"] = "lightblue"
		}
		if state.IsFinal {
			attrs["shape"] = "doublecircle"
		}
		graph.AddNode("G", state.Data, attrs)
	}

	for _, transition := range automaton.Transitions {
		attrs := make(map[string]string)
		attrs["label"] = "\"" + strings.Join(transition.Chars, ", ") + "\""
		graph.AddEdge(transition.Start, transition.End, true, attrs)
	}

	return graph, nil
}

// "D:\\Program Files\\Graphviz\\bin\\dot.exe"
func saveGraphAsPNG(graph *gographviz.Graph, outputPath string) error {
	dot := graph.String()

	cmd := exec.Command("D:\\Program Files\\Graphviz\\bin\\dot.exe", "-Tpng")
	cmd.Stdin = strings.NewReader(dot)
	var output bytes.Buffer
	cmd.Stdout = &output

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error running dot command: %v. Output: %s", err, output.String())
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	if _, err := io.Copy(outputFile, bytes.NewReader(output.Bytes())); err != nil {
		return err
	}

	return nil
}
