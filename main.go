package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/awalterschulze/gographviz"
)

type State struct {
	Data     string   `json:"data"`
	IsStart  bool     `json:"isStart"`
	IsEnd    bool     `json:"isEnd"`
	Adjacent []string `json:"adjacent"`
}

type Transition struct {
	Start string   `json:"start"`
	End   string   `json:"end"`
	Chars []string `json:"chars"`
}

type Automaton struct {
	Name        string       `json:"name"`
	Alphabet    string       `json:"alphabet"`
	States      []State      `json:"states"`
	Transitions []Transition `json:"transitions"`
}

func main() {
	// Leer el archivo JSON
	jsonData, err := os.ReadFile("./data/auto.json")
	if err != nil {
		fmt.Println("Error al leer el archivo:", err)
		os.Exit(1)
	}

	// Decodificar el JSON
	var automaton Automaton
	err = json.Unmarshal(jsonData, &automaton)
	if err != nil {
		fmt.Println("Error al decodificar el JSON:", err)
		os.Exit(1)
	}

	// Crear un grafo con gographviz
	graphAst, _ := gographviz.ParseString("digraph G {}")
	graph := gographviz.NewGraph()
	gographviz.Analyse(graphAst, graph)

	// Añadir estados al grafo
	for _, state := range automaton.States {
		graph.AddNode("G", state.Data, nil)
	}

	// Añadir transiciones al grafo
	for _, transition := range automaton.Transitions {
		for _, char := range transition.Chars {
			attrs := make(map[string]string)
			attrs["label"] = char
			graph.AddEdge(transition.Start, transition.End, true, attrs)
		}
	}

	// Generar el código DOT del grafo
	dotOutput := graph.String()
	fmt.Println(dotOutput)

	// Para visualizar el grafo, copia y pega el código DOT en un
	// visualizador en línea de Graphviz, como https://dreampuf.github.io/GraphvizOnline/
}
