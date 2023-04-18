package controllers

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/SRG98/automatas-go/models"
)

const (
	inputJSONFile    = "data/auto.json"
	inputTextFile    = "data/input.txt"
	AbsInputTextFile = "C:/Users/oh/3D Objects/Computación/Automatas/JavaScript/autovii/"
	outputImagePath  = "data/graph.png"
)

type Controller struct {
	selectedAutomata *models.Automata
	AutomatsList     []*models.Automata
	inputStrings     []string
	function         *models.Function
}

func NewController() *Controller {
	return &Controller{
		AutomatsList: make([]*models.Automata, 0),
		inputStrings: make([]string, 0),
	}
}

func (c *Controller) GetAutomatsList() []*models.Automata {
	return c.AutomatsList
}

func (c *Controller) SetInputStrings(strings []string) {
	c.inputStrings = strings
}

func (c *Controller) GetInputStrings() []string {
	return c.inputStrings
}

func (c *Controller) Run() error {
	// views.RunUI()
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
			// err = c.CreateAutomata()
		case 2:
			// err = c.selectAutomata()
		case 3:
			// err = c.createState()
		case 4:
			// err = c.createTransition()
		case 5:
			err = c.readInputFile(inputTextFile)
		case 6:
			err = c.validateString()
		case 7:
			// err = c.generateImage()
		case 8:
			err = c.viewAutomata()
		case 9:
			err = c.viewStrings()
		case 10:
			// err = c.processInputStrings()
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

func (c *Controller) CreateAutomata(name string) bool {
	if name == "" {
		fmt.Println("Nombre vacío")
		return false
	}

	auto := models.NewAutomaton()
	auto.SetName(name)

	// Añadir el objeto Automata a la lista de autómatas (c.automataList) y seleccionarlo como el autómata actual (c.selectedAutomaton).
	c.AutomatsList = append(c.AutomatsList, auto)
	c.selectedAutomata = auto

	fmt.Println("Autómata creado y guardado exitosamente.")

	// FUNCIONALIDAD PELIGROSA
	c.SelectAutomata(len(c.AutomatsList) - 1)

	// Guardar el autómata en el archivo JSON
	if c.writeJSONFile(inputJSONFile, auto) {
		fmt.Println("Json guardado")
		c.GenerateImage()
		return true
	}
	return false
}

func (c *Controller) SelectAutomata(index int) bool {
	fmt.Println("POS", index)

	if len(c.AutomatsList) == 0 {
		fmt.Println("no hay autómatas disponibles")
		return false
	}

	if index == -1 {
		fmt.Println("Negative Access")
	}

	if index >= len(c.AutomatsList) {
		fmt.Println("índice de autómata inválido")
		return false
	}

	c.selectedAutomata = c.AutomatsList[index]
	if c.function == nil {
		fmt.Println("No FN")
		c.function = models.NewFunction(c.selectedAutomata)
	} else {
		fmt.Println("has fn")
		c.function.SetAutomata(c.selectedAutomata)
	}

	fmt.Printf("Autómata '%s' seleccionado exitosamente.\n", c.selectedAutomata.Name)
	c.GenerateImage()
	return true
}

func (c *Controller) CreateState(data string, isInitial bool, isFinal bool) bool {
	if c.selectedAutomata == nil {
		fmt.Println("ningún autómata seleccionado")
		return false
	}

	if data == "" {
		return false
	}

	if c.selectedAutomata.NewState(data, isInitial, isFinal) {
		fmt.Println("estado creado exitosamente.")
		// RETURN FINAL
		// Guardar el autómata en el archivo JSON
		c.GenerateImage()
		return c.writeJSONFile(inputJSONFile, c.selectedAutomata)
	}
	fmt.Println("estado ya existente.")
	return false
}

func (c *Controller) CreateTransition(start string, end string, charsStr string) bool {
	if start == "" || end == "" || charsStr == "" {
		fmt.Println("Blank form")
		return false
	}

	if c.selectedAutomata == nil {
		fmt.Print("ningún autómata seleccionado")
		return false
	}

	if len(c.selectedAutomata.States) == 0 {
		fmt.Print("el autómata no tiene estados")
		return false
	}

	if !c.selectedAutomata.ExistState(start) || !c.selectedAutomata.ExistState(end) {
		fmt.Println("el estado final o de inicio no existe en el autómata")
		return false
	}

	chars := strings.Split(charsStr, ",")

	if c.selectedAutomata.NewTransition(start, end, chars) {
		fmt.Print("Nueva transición creada")
		// Guardar el autómata en el archivo JSON
		c.GenerateImage()
		return c.writeJSONFile(inputJSONFile, c.selectedAutomata)
	}
	return false
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

func (c *Controller) GenerateImage() error {
	if c.selectedAutomata == nil {
		return fmt.Errorf("ningún autómata seleccionado")
	}

	outputPath := outputImagePath

	err := CreateImage(c.selectedAutomata, outputPath)
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

func (c *Controller) ProcessInputStrings() ([]bool, error) {
	// list := []bool{false, true, false, true, false}
	// return list, nil

	validations := []bool{}

	if c.selectedAutomata == nil {
		return validations, fmt.Errorf("ningún autómata seleccionado")
	}

	if len(c.inputStrings) == 0 {
		return validations, fmt.Errorf("no hay cadenas de entrada para procesar")
	}
	fmt.Println("List: ", c.inputStrings)
	// listRet []bool{}

	for _, inputString := range c.inputStrings {
		// Procesa y valida la cadena de entrada con el autómata seleccionado
		// Puedes reemplazar esta parte con la lógica adecuada de validación de cadenas
		// c.function.SetString(inputString)
		inputString = strings.TrimSpace(inputString)
		// fmt.Println(inputString)
		isValid := c.function.Validate(inputString)

		validations = append(validations, isValid)
	}
	return validations, nil
}

func clearInputBuffer(reader *bufio.Reader) {
	for {
		_, _, err := reader.ReadRune()
		if err != nil || reader.Buffered() == 0 {
			break
		}
	}
}
