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
	AbsInputTextFile = "D:/Código/automatas/goiii/data"
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
	//views.RunUI()
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
		case 5:
			c.readInputFile(inputTextFile)
		case 6:
			err = c.validateString()
		case 7:
			c.GenerateImage()
		case 8:
			c.viewAutomata()
		case 9:
			c.viewStrings()
		case 10:
			c.ProcessInputStrings()
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
	fmt.Println("| 1. Ingresar un autómata")
	fmt.Println("| 2. Cargue un autómata")
	fmt.Println("| 3. Ingresar un estado")
	fmt.Println("| 4. Crear transicion")
	fmt.Println("| 5. Ingresar un texto")
	fmt.Println("| 6. Valide la cadena")
	fmt.Println("| 7. Crear imagen del automata")
	fmt.Println("| 8. Observar el automata")
	fmt.Println("| 9. Ver las cadenas creadas")
	fmt.Println("| 10. Procesar cadenas")
	fmt.Println("| Salir")
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
	for i := range chars {
		chars[i] = strings.TrimSpace(chars[i])
	}

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
	if err == nil {
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
	list := []bool{false, true, false, true, false}
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
