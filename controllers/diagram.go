package controllers

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/SRG98/automatas-go/models"
	"github.com/awalterschulze/gographviz"
)

func CreateImage(automaton *models.Automata, outputPath string) error {
	graph, err := automataToGraphviz(automaton)
	if err != nil {
		return fmt.Errorf("error generating Graphviz graph: %v", err)
	}

	err = saveGraphAsPNG(graph, outputPath)
	if err != nil {
		return fmt.Errorf("error saving Graphviz graph as PNG: %v", err)
	}

	// err = convertPNGtoJPG("data/graph.png", "data/graph.jpg")
	// if err != nil {
	// 	fmt.Println("Error al convertir la imagen:", err)
	// } else {
	// 	fmt.Println("Imagen convertida exitosamente")
	// }

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

func saveGraphAsPNG(graph *gographviz.Graph, outputPath string) error {
	dot := graph.String()

	cmd := exec.Command("dot", "-Tpng")
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

// func convertPNGtoJPG(pngFilePath, jpgFilePath string) error {
// 	// Abrir archivo PNG
// 	pngFile, err := os.Open(pngFilePath)
// 	if err != nil {
// 		return fmt.Errorf("error abriendo archivo PNG: %v", err)
// 	}
// 	defer pngFile.Close()

// 	// Decodificar archivo PNG
// 	pngImage, err := png.Decode(pngFile)
// 	if err != nil {
// 		return fmt.Errorf("error decodificando archivo PNG: %v", err)
// 	}

// 	// Crear archivo JPG
// 	jpgFile, err := os.Create(jpgFilePath)
// 	if err != nil {
// 		return fmt.Errorf("error creando archivo JPG: %v", err)
// 	}
// 	defer jpgFile.Close()

// 	// Codificar la imagen PNG en formato JPEG y guardarla en el archivo JPG
// 	var quality = 75 // Define la calidad de la imagen JPEG (entre 1 y 100)
// 	err = jpeg.Encode(jpgFile, pngImage, &jpeg.Options{Quality: quality})
// 	if err != nil {
// 		return fmt.Errorf("error codificando imagen en formato JPG: %v", err)
// 	}

// 	return nil
// }
