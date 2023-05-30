package controllers

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/SRG98/automatas-go/logic"
	"github.com/SRG98/automatas-go/models"
)

var (
	inputJSONFile    = "data/automatas.json"
	inputTextFile    = "data/cadenas.txt"
	AbsInputTextFile = "C:/Código/automatas/goiii/data"
	outputImagePath  = "data/graph.png"
)

type Controller struct {
	selectedAutomata *models.Automata
	AutomatsList     []*models.Automata
	Determiner       *logic.Determiner
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
			var name string
			fmt.Println("Ingrese el nombre del automata")
			fmt.Scanln(&name)
			c.CreateAutomata(name)
		case 2:
			var val int
			fmt.Println("Ingrese el valor")
			fmt.Scanln(&val)
			c.SelectAutomata(val)
		case 3:
			var name string
			fmt.Println("Ingrese el estado")
			fmt.Scanln(&name)
			var ini bool
			fmt.Println("Ingrese si es estado inicial")
			fmt.Scanln(&ini)
			var fin bool
			fmt.Println("Ingrese si es estado final")
			fmt.Scanln(&fin)
			c.CreateState(name, ini, fin)
		case 4:
			var sal string
			fmt.Println("Ingrese el estado de salida")
			fmt.Scanln(&sal)
			var lle string
			fmt.Println("Ingrese el estado de llegada")
			fmt.Scanln(&lle)
			var val string
			fmt.Println("Ingrese la transicion")
			fmt.Scanln(&val)
			c.CreateTransition(sal, lle, val)
		case 6:
			c.readInputFile(inputTextFile)
		case 5:
			c.validateString()
		case 7:
			c.GenerateImage()
		case 8:
			c.viewAutomata()
		case 9:
			c.viewStrings()
		case 10:
			c.validateStringsFromFile(inputTextFile)
		case 0:
			return nil
		default:
			fmt.Println("Opción no válida. Por favor, intente de nuevo.")
			continue
		}

		if err != nil {
			fmt.Println("Error:", err)
		}
		_, _ = reader.ReadString('\n')
		_, _ = fmt.Scanln()
	}
}

func (c *Controller) showMenu() {
	fmt.Println("\n-- PROYECTO AUTOMATAS 2023-1 --")
	fmt.Println("Bienvenido al sistema")
	fmt.Println("Seleccione la accion que desea realizar : ")
	fmt.Println("Opt                                       ")
	fmt.Println("(1) Ingresar un autómata")
	fmt.Println("(2) Cargue un autómata")
	fmt.Println("(3) Ingresar un estado")
	fmt.Println("(4) Crear transicion")
	fmt.Println("(5) Valide la cadena")
	fmt.Println("(6) Cargar cadenas")
	fmt.Println("(7) Crear imagen del automata")
	fmt.Println("(8) Observar el automata")
	fmt.Println("(9) Ver las cadenas creadas")
	fmt.Println("(10) Validar cadenas cargas desde txt")
	fmt.Println("Salir")
}

func (c *Controller) readOption() (int, error) {
	var option int
	_, err := fmt.Scan(&option)
	return option, err
}

func (c *Controller) CreateAutomata(name string) bool {
	if name == "" {
		fmt.Println("Campo vacío")
		return false
	}

	auto := models.NewAutomaton()
	auto.SetName(name)

	// Añadir el objeto Automata a la lista de autómatas (c.automataList) y seleccionarlo como el autómata actual (c.selectedAutomaton).
	c.AutomatsList = append(c.AutomatsList, auto)
	c.selectedAutomata = auto

	fmt.Println("Autómata creado y guardado.")

	// FUNCIONALIDAD PELIGROSA
	c.SelectAutomata(len(c.AutomatsList) - 1)

	// Guardar el autómata en el archivo JSON
	if c.writeJSONFile(inputJSONFile, auto) {
		fmt.Println("Json guardado correctamente")
		c.GenerateImage()
		return true

	}
	return false
}

func (c *Controller) AddAutomata(automata *models.Automata) bool {
	if automata == nil {
		fmt.Println("Automata nil")
		return false
	}

	// Añadir el objeto Automata a la lista de autómatas (c.automataList) y seleccionarlo como el autómata actual (c.selectedAutomaton).
	c.AutomatsList = append(c.AutomatsList, automata)
	c.selectedAutomata = automata

	fmt.Println("Autómata creado y guardado exitosamente.")

	// FUNCIONALIDAD PELIGROSA
	c.SelectAutomata(len(c.AutomatsList) - 1)

	// Guardar el autómata en el archivo JSON
	if c.writeJSONFile(inputJSONFile, automata) {
		fmt.Println("Json guardado")
		c.GenerateImage()
		return true
	}
	return false
}

func (c *Controller) NormalizeAutomata() error {
	normal := logic.NewDeterminer()

	if c.selectedAutomata == nil {
		return fmt.Errorf("ningún autómata seleccionado")
	}

	normal.SetAutomata(c.selectedAutomata)
	newAutomata := normal.Determine()

	// if err != nil {
	// 	return nil, fmt.Errorf("error al normalizar el autómata: %v", err)
	// }

	c.AutomatsList = append(c.AutomatsList, newAutomata)
	c.selectedAutomata = newAutomata

	c.SelectAutomata(len(c.AutomatsList) - 1)

	// fmt.Println("Autómata creado y guardado exitosamente.")

	// FUNCIONALIDAD PELIGROSA
	fmt.Print(c.selectedAutomata.ToString())

	// Guardar el autómata en el archivo JSON
	// if c.writeJSONFile(inputJSONFile, newAuto) {
	// 	fmt.Println("Json guardado")
	// }

	// if newAutomata == nil {
	// 	return fmt.Errorf("no se pudo determinar el autómata")
	// }

	return nil
}

func (c *Controller) GenerateAutomata() (*models.Automata, error) {
	newAuto := c.createAutoI()

	c.AutomatsList = append(c.AutomatsList, newAuto)
	c.selectedAutomata = newAuto

	c.SelectAutomata(len(c.AutomatsList) - 1)

	fmt.Print(c.selectedAutomata.ToString())

	return newAuto, nil
}

func (c *Controller) createAutoI() *models.Automata {
	newAuto := models.NewAutomaton()

	newAuto.SetName("AFND")

	newAuto.NewState("A", true, false)
	newAuto.NewState("B", false, false)
	newAuto.NewState("C", false, false)
	newAuto.NewState("D", false, false)
	newAuto.NewState("E", false, false)
	newAuto.NewState("F", false, false)
	newAuto.NewState("G", false, false)
	newAuto.NewState("H", false, false)
	newAuto.NewState("I", false, false)
	newAuto.NewState("J", false, true)

	newAuto.NewTransition("A", "B", []string{"a"})
	newAuto.NewTransition("B", "C", []string{"_"})
	newAuto.NewTransition("C", "D", []string{"a"})
	newAuto.NewTransition("C", "D", []string{"_"})
	newAuto.NewTransition("D", "E", []string{"_"})
	newAuto.NewTransition("E", "F", []string{"_"})
	newAuto.NewTransition("E", "H", []string{"_"})
	newAuto.NewTransition("F", "G", []string{"a"})
	newAuto.NewTransition("H", "I", []string{"b"})
	newAuto.NewTransition("G", "J", []string{"_"})
	newAuto.NewTransition("I", "J", []string{"_"})

	return newAuto
}

func (c *Controller) SelectAutomata(index int) bool {

	if len(c.AutomatsList) == 0 {
		fmt.Println("No se han creado automatas")
		return false
	}

	if index == -1 {
		fmt.Println("No se puede acceder a esta posicion")
	}

	if index >= len(c.AutomatsList) {
		fmt.Println("índice de autómata inválido")
		return false
	}

	c.selectedAutomata = c.AutomatsList[index]
	if c.function == nil {
		fmt.Println(" ")
		c.function = models.NewFunction(c.selectedAutomata)
	} else {
		fmt.Println(" ")
		c.function.SetAutomata(c.selectedAutomata)
	}

	fmt.Printf("Autómata '%s' seleccionado exitosamente.\n", c.selectedAutomata.Name)
	c.GenerateImage()
	return true

}

func (c *Controller) CreateState(data string, isInitial bool, isFinal bool) bool {
	if c.selectedAutomata == nil {
		fmt.Println("ningún autómata ha sido seleccionado")
		return false
	}

	if data == "" {
		return false
	}

	if c.selectedAutomata.NewState(data, isInitial, isFinal) {
		fmt.Println("estado creado correctamente.")

		// Guardar el autómata en el archivo JSON
		c.GenerateImage()
		return c.writeJSONFile(inputJSONFile, c.selectedAutomata)
	}
	fmt.Println("estado ya existente.")
	return false
}

func (c *Controller) CreateTransition(start string, end string, charsStr string) bool {
	if start == "" || end == "" || charsStr == "" {
		fmt.Println("Campo Vacio")
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
	for i := range chars {
		chars[i] = strings.TrimSpace(chars[i])
	}

	if c.selectedAutomata.NewTransition(start, end, chars) {
		fmt.Print("Nueva transición creada")

		c.GenerateImage()
		return c.writeJSONFile(inputJSONFile, c.selectedAutomata)
	}
	return true
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
		return fmt.Errorf("Error al generar la imagen: %v", err)
	}

	fmt.Println("Imagen generada exitosamente en", outputPath)
	return nil
}

func (c *Controller) viewAutomata() error {
	if c.selectedAutomata == nil {
		return fmt.Errorf("Ningún autómata seleccionado")
	}

	fmt.Println("Información del autómata seleccionado:")
	fmt.Println("Nombre:", c.selectedAutomata.Name)

	fmt.Println("Estados:")

	for _, state := range c.selectedAutomata.States {
		fmt.Printf("- %s (Estado de salida: %v, Estado de llegada: %v)\n", state.Data, state.IsInitial, state.IsFinal)
	}

	fmt.Println("Transiciones:")

	for _, transition := range c.selectedAutomata.Transitions {
		fmt.Printf("- %s → %s (transicion: %s)\n", transition.Start, transition.End, strings.Join(transition.Chars, ", "))
	}
	if len(c.selectedAutomata.Transitions) <= 2 {
		println("Automata incompleto")
	} else {
		println("Automata completo")
	}

	return nil
}

func (c *Controller) viewStrings() error {

	if len(c.inputStrings) == 0 {
		return fmt.Errorf("no hay cadenas de entrada para mostrar")
	}

	fmt.Println("Cadenas de Cargadas:")
	for i, inputString := range c.inputStrings {
		fmt.Printf("%d. %s\n", i+1, inputString)
	}
	return nil
}

func (c *Controller) ProcessInputStrings() ([]bool, error) {

	if c.selectedAutomata == nil {
		return nil, fmt.Errorf("Ningun automata seleccionado")
	}
	if len(c.inputStrings) == 0 {
		return nil, fmt.Errorf("No hay cadenas para procesar")
	}
	validations := make([]bool, len(c.inputStrings))

	for i, inputString := range c.inputStrings {
		inputString = strings.TrimSpace(inputString)
		isValid := c.function.Validate(inputString)
		validations[i] = isValid
		print(validations)
	}

	return validations, nil
	/*list := []bool{false, true, false, true, false}
	return list, nil

	validations := []bool{}

	if c.selectedAutomata == nil {
		return validations, fmt.Errorf("ningún autómata seleccionado")
	}

	if len(c.inputStrings) == 0 {
		return validations, fmt.Errorf("no hay cadenas de entrada para procesar")
	}
	fmt.Println("List: ", c.inputStrings)

	for _, inputString := range c.inputStrings {

		// Procesa y valida la cadena de entrada con el autómata seleccionado
		inputString = strings.TrimSpace(inputString)

		isValid := c.function.Validate(inputString)

		validations = append(validations, isValid)
	}
	return list, nil*/
}

func clearInputBuffer(reader *bufio.Reader) {
	for {
		_, _, err := reader.ReadRune()
		if err != nil || reader.Buffered() == 0 {
			break
		}
	}
}
func (c *Controller) validateStringsFromFile(filepath string) error {
	if c.selectedAutomata == nil {
		return fmt.Errorf("ningún autómata seleccionado")
	}

	// Leer el contenido del archivo
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}

	// Dividir el contenido en una lista de cadenas
	inputStrings := strings.Split(string(content), "\n")

	// Validar cada cadena
	for _, inputStr := range inputStrings {
		inputStr = strings.TrimSpace(inputStr)
		if inputStr == "" {
			continue
		}

		isValid := c.function.Validate(inputStr)
		if isValid {
			fmt.Printf("La cadena \"%s\" es válida.\n", inputStr)
		} else {
			fmt.Printf("La cadena \"%s\" es inválida.\n", inputStr)
		}
	}

	return nil
}
