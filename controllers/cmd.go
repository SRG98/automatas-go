package controllers

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/SRG98/automatas-go/models"
	"github.com/SRG98/automatas-go/views"
)

const (
	inputJSONFile   = "data/auto.json"
	inputTextFile   = "data/input.txt"
	outputImagePath = "data/graph.png"
)

type Controller struct {
	selectedAutomata *models.Automata
	automataList     []*models.Automata
	inputStrings     []string
	function         *models.Function
}

func NewController() *Controller {
	return &Controller{
		automataList: make([]*models.Automata, 0),
		inputStrings: make([]string, 0),
	}
}

func (c *Controller) Run() error {
	reader := bufio.NewReader(os.Stdin)
	for {
		c.showMenu()
		option, err := c.readOption()
		clearInputBuffer(reader)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		switch option {
		case 1:
			err = c.createAutomata()
		case 2:
			err = c.selectAutomata()
		case 3:
			err = c.createState()
		case 4:
			err = c.createTransition()
		case 5:
			err = c.readInputFile(inputTextFile)
		case 6:
			err = c.validateString()
		case 7:
			err = c.generateImage()
		case 8:
			err = c.viewAutomata()
		case 9:
			err = c.viewStrings()
		case 10:
			err = c.processInputStrings()
		case 0:
			return nil
		default:
			fmt.Println("Opción no válida. Por favor, intente de nuevo.")
			continue
		}

		if err != nil {
			fmt.Println("Error:", err)
		}
		_, _ = reader.ReadString('\n') // Agrega esta línea
		// _, _ = fmt.Scanln() // Agrega esta línea
	}
}

func (c *Controller) showMenu() {
	fmt.Println("\n-------------------------------------------------")
	fmt.Println("| 01. Crear autómata | 02. Seleccionar autómata |")
	fmt.Println("| 03. Crear estado   | 04. Crear transición     |")
	fmt.Println("| 05. Ingresar texto | 06. Validar cadena       |")
	fmt.Println("| 07. Generar imagen | 08. Ver autómata         |")
	fmt.Println("| 09. Ver cadenas    | 10. Procesar cadenas     |")
	// fmt.Println("| 00.                | 00.                      |")
	fmt.Println("| 00. Salir          |                          |")
	fmt.Println("-------------------------------------------------")
	fmt.Print("Seleccione una opción: ")
}

func (c *Controller) readOption() (int, error) {
	var option int
	_, err := fmt.Scan(&option)
	return option, err
}

func (c *Controller) createAutomata() error {
	// Solicitar el nombre del autómata y crear un nuevo objeto Automata con ese nombre.
	fmt.Print("Ingrese el nombre del autómata: ")
	var name string
	_, err := fmt.Scan(&name)
	if err != nil {
		return err
	}

	auto := models.NewAutomaton()
	auto.SetName(name)

	// Añadir el objeto Automata a la lista de autómatas (c.automataList) y seleccionarlo como el autómata actual (c.selectedAutomaton).
	c.automataList = append(c.automataList, auto)
	c.selectedAutomata = auto

	fmt.Println("Autómata creado exitosamente.")

	c.function = models.NewFunction(c.selectedAutomata)

	// Guardar el autómata en el archivo JSON
	err = c.writeJSONFile(inputJSONFile, auto)
	if err != nil {
		return err
	}
	return nil
}

func (c *Controller) selectAutomata() error {
	if len(c.automataList) == 0 {
		return fmt.Errorf("no hay autómatas disponibles")
	}

	fmt.Println("Seleccione el índice del autómata que desea seleccionar:")
	for i, auto := range c.automataList {
		fmt.Printf("%d. %s\n", i+1, auto.Name)
	}

	fmt.Print("Índice: ")

	var index int
	_, err := fmt.Scanf("%d", &index)
	if err != nil {
		return err
	}

	if index < 1 || index > len(c.automataList) {
		return fmt.Errorf("índice de autómata inválido")
	}

	c.selectedAutomata = c.automataList[index-1]
	c.function = models.NewFunction(c.selectedAutomata)
	fmt.Printf("Autómata '%s' seleccionado exitosamente.\n", c.selectedAutomata.Name)
	return nil
}

func (c *Controller) createState() error {
	if c.selectedAutomata == nil {
		return fmt.Errorf("ningún autómata seleccionado")
	}

	var data, isInitialStr, isFinalStr string
	var isInitial, isFinal bool

	fmt.Print("Ingrese el nombre del estado: ")
	_, err := fmt.Scan(&data)
	if err != nil {
		return err
	}

	// Solicitar si el estado es inicial
	fmt.Print("Es este estado inicial? (s/n): ")
	_, err = fmt.Scan(&isInitialStr)
	if err != nil {
		return err
	}
	isInitial = (isInitialStr == "s")

	// Solicitar si el estado es final
	fmt.Print("Es este estado final? (s/n): ")
	_, err = fmt.Scan(&isFinalStr)
	if err != nil {
		return err
	}
	isFinal = (isFinalStr == "s")

	c.selectedAutomata.NewState(data, isInitial, isFinal)
	fmt.Println("Estado creado exitosamente.")

	// Guardar el autómata en el archivo JSON
	err = c.writeJSONFile(inputJSONFile, c.selectedAutomata)
	if err != nil {
		return err
	}
	return nil
}

func (c *Controller) createTransition() error {
	if c.selectedAutomata == nil {
		return fmt.Errorf("ningún autómata seleccionado")
	}

	if len(c.selectedAutomata.States) == 0 {
		return fmt.Errorf("el autómata no tiene estados")
	}

	var start, end, charsStr string

	fmt.Print("Ingrese el estado de inicio: ")
	_, err := fmt.Scan(&start)
	if err != nil {
		return err
	}

	if !c.selectedAutomata.ExistState(start) {
		return fmt.Errorf("el estado de inicio no existe en el autómata")
	}

	fmt.Print("Ingrese el estado final: ")
	_, err = fmt.Scan(&end)
	if err != nil {
		return err
	}

	if !c.selectedAutomata.ExistState(end) {
		return fmt.Errorf("el estado final no existe en el autómata")
	}

	fmt.Print("Ingrese los caracteres de la transición (separados por comas): ")
	_, err = fmt.Scan(&charsStr)
	if err != nil {
		return err
	}
	chars := strings.Split(charsStr, ",")

	err = c.selectedAutomata.NewTransition(start, end, chars)
	if err != nil {
		return err
	}

	// Guardar el autómata en el archivo JSON
	err = c.writeJSONFile(inputJSONFile, c.selectedAutomata)
	if err != nil {
		return err
	}
	return nil
}

func (c *Controller) validateString() error {
	if c.selectedAutomata == nil {
		return fmt.Errorf("ningún autómata seleccionado")
	}

	fmt.Print("Ingrese la cadena a validar: ")
	reader := bufio.NewReader(os.Stdin)
	inputStr, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	inputStr = strings.TrimSpace(inputStr)

	isValid := c.function.Validate(inputStr)
	if isValid {
		fmt.Println("La cadena es válida.")
	} else {
		fmt.Println("La cadena es inválida.")
	}

	return nil
}

func (c *Controller) generateImage() error {
	if c.selectedAutomata == nil {
		return fmt.Errorf("ningún autómata seleccionado")
	}

	outputPath := outputImagePath

	err := views.GenerateImage(c.selectedAutomata, outputPath)
	if err != nil {
		return fmt.Errorf("error al generar la imagen: %v", err)
	}

	fmt.Println("Imagen generada exitosamente en", outputPath)
	return nil
}

func (c *Controller) viewAutomata() error {
	if c.selectedAutomata == nil {
		return fmt.Errorf("ningún autómata seleccionado")
	}

	fmt.Println("Información del autómata seleccionado:")
	fmt.Println("Nombre:", c.selectedAutomata.Name)
	// fmt.Println("Alfabeto:", strings.Join(c.selectedAutomata.Alphabet, ", "))
	fmt.Println("Estados:")

	for _, state := range c.selectedAutomata.States {
		fmt.Printf("- %s (inicial: %v, final: %v)\n", state.Data, state.IsInitial, state.IsFinal)
	}

	fmt.Println("Transiciones:")

	for _, transition := range c.selectedAutomata.Transitions {
		fmt.Printf("- %s → %s (símbolos: %s)\n", transition.Start, transition.End, strings.Join(transition.Chars, ", "))
	}

	return nil
}

func (c *Controller) viewStrings() error {

	if len(c.inputStrings) == 0 {
		return fmt.Errorf("no hay cadenas de entrada para mostrar")
	}

	fmt.Println("Cadenas de entrada:")
	for i, inputString := range c.inputStrings {
		fmt.Printf("%d. %s\n", i+1, inputString)
	}

	return nil
}

func (c *Controller) processInputStrings() error {
	if c.selectedAutomata == nil {
		return fmt.Errorf("ningún autómata seleccionado")
	}

	if len(c.inputStrings) == 0 {
		return fmt.Errorf("no hay cadenas de entrada para procesar")
	}
	fmt.Println("List: ", c.inputStrings)

	for _, inputString := range c.inputStrings {
		// Procesa y valida la cadena de entrada con el autómata seleccionado
		// Puedes reemplazar esta parte con la lógica adecuada de validación de cadenas
		// c.function.SetString(inputString)
		inputString = strings.TrimSpace(inputString)
		fmt.Println(inputString)
		isValid := c.function.Validate(inputString)
		if isValid {
			fmt.Println("La cadena", inputString, "es válida para este autómata.")
		} else {
			fmt.Println("La cadena", inputString, "no es válida para este autómata.")
		}
	}

	return nil
}

func clearInputBuffer(reader *bufio.Reader) {
	for {
		_, _, err := reader.ReadRune()
		if err != nil || reader.Buffered() == 0 {
			break
		}
	}
}