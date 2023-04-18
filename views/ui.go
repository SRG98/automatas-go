package views

import (
	"fmt"
	"image/color"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage" // Asegúrate de importar este paquete
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/SRG98/automatas-go/controllers"
)

type UI struct {
	controller            *controllers.Controller
	updateImageFunc       func()
	inputStringsContainer *fyne.Container
}

func NewUI(controller *controllers.Controller) *UI {
	return &UI{
		controller: controller,
	}
}

type ButtonConfig struct {
	Text     string
	OnTapped func()
}

func (ui *UI) RunUI() error {
	app := app.New()
	app.Settings().SetTheme(theme.DefaultTheme())

	win := app.NewWindow("Go-Automats")

	// Crear botones
	buttons := ui.createButtons(win, 8)

	// Asignar una función al primer botón
	// buttons[0].OnTapped = func() {
	// 	ui.showCreateAutomataDialog(win)
	// }

	// buttons[9].OnTapped = func() {
	// 	ui.showCreateAutomataDialog(win)
	// }

	// Crear el contenedor de botones (vertical) con el botón de mostrar imagen
	buttonObjects := make([]fyne.CanvasObject, len(buttons))
	for i, button := range buttons {
		buttonObjects[i] = button
	}
	buttonContainer := container.NewVBox(buttonObjects...)

	// Crear contenedor principal (horizontal) con botones e imagen
	mainContainer := container.NewHBox(buttonContainer)

	win.SetContent(mainContainer)
	win.Resize(fyne.NewSize(800, 600))

	// Mostrar la ventana de la imagen
	ui.updateImageFunc = ui.showImageWindow(app)

	win.ShowAndRun()

	return nil
}

func (ui *UI) createButtons(win fyne.Window, numButtons int) []*widget.Button {
	buttons := make([]*widget.Button, numButtons)

	for i := 0; i < numButtons; i++ {
		var text string
		var onTapped func()

		switch i {
		case 0:
			text = "Crear autómata"
			onTapped = func() {
				ui.showCreateAutomataDialog(win)
			}
		case 1:
			text = "Seleccionar autómata"
			onTapped = func() {
				ui.showAutomatsListDialog(win)
			}
		case 2:
			text = "Crear estado"
			onTapped = func() {
				ui.showStateConfigDialog(win)
			}
		case 3:
			text = "Crear transición"
			onTapped = func() {
				ui.showTransitionConfigDialog(win)
			}
		case 4:
			text = "Ingresar texto"
			onTapped = func() {
				ui.showFileChooserDialog(win)
			}
		case 5:
			text = "Ver cadenas"
			onTapped = func() {
				ui.showCreateAutomataDialog(win)
			}
		case 6:
			text = "Procesar cadenas"
			onTapped = func() {
				ui.validateInputStrings()
			}
		default:
			text = fmt.Sprintf("Botón %d", i+1)
			onTapped = nil
		}

		button := widget.NewButton(text, onTapped)
		buttons[i] = button
	}

	return buttons
}

func (u *UI) showImageWindow(app fyne.App) func() {
	// Obtener la ruta del directorio actual
	currentDir, err := os.Getwd()
	if err != nil {
		log.Printf("No se pudo obtener el directorio actual: %v\n", err)
		return nil
	}

	// Construir la ruta absoluta a la imagen
	imagePath := filepath.Join(currentDir, "data", "graph.png")

	// Cargar la imagen usando la ruta absoluta
	image := loadImage(imagePath)
	image.FillMode = canvas.ImageFillContain

	// Crear una nueva ventana redimensionable
	imageWindow := app.NewWindow("Imagen")
	imageWindow.Resize(fyne.NewSize(600, 400))

	// Agregar un ScrollContainer para contener la imagen
	imageMaxContainer := container.NewMax(image)

	// Agregar un ScrollContainer para contener el container.Max
	scrollContainer := container.NewScroll(imageMaxContainer)

	imageWindow.SetContent(scrollContainer)
	imageWindow.Show()

	return func() {
		newImage := loadImage(imagePath)
		newImage.FillMode = canvas.ImageFillContain
		imageMaxContainer.Objects = []fyne.CanvasObject{newImage}
		imageMaxContainer.Refresh()
		scrollContainer.Refresh()
	}
}

func loadImage(filePath string) *canvas.Image {
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("No se pudo abrir el archivo: %s\n", filePath)
		return canvas.NewImageFromResource(theme.FyneLogo())
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("No se pudo leer el archivo: %s\n", filePath)
		return canvas.NewImageFromResource(theme.FyneLogo())
	}

	res := fyne.NewStaticResource(filepath.Base(filePath), data)
	image := canvas.NewImageFromResource(res)
	if image.Resource == nil {
		log.Printf("No se pudo cargar la imagen desde: %s\n", filePath)
		return canvas.NewImageFromResource(theme.FyneLogo())
	}
	return image
}

func (u *UI) showCreateAutomataDialog(w fyne.Window) {
	dialogContent := container.NewVBox(
		widget.NewLabel("Crear autómata"),
	)

	entry := widget.NewEntry()
	entry.SetPlaceHolder("nombre")
	dialogContent.Add(entry)

	dialog := dialog.NewCustom("Crear autómata", "", dialogContent, w)

	createButton := widget.NewButton("Aceptar", func() {
		automataName := entry.Text
		if u.controller.CreateAutomata(automataName) {
			fmt.Println("autómata creado:", automataName)
			u.updateImageFunc()
			dialog.Hide()
		}
	})

	cancelButton := widget.NewButton("Cancelar", func() {
		dialog.Hide()
	})

	buttons := container.NewHBox(createButton, cancelButton)
	dialogContent.Add(buttons)

	dialog.Show()
}

func (u *UI) showAutomatsListDialog(w fyne.Window) {
	// Lista de strings de ejemplo
	// stringList := []string{"Opción 1", "Opción 2", "Opción 3"}
	automatsList := u.controller.GetAutomatsList()
	stringList := make([]string, len(automatsList))

	for i, auto := range automatsList {
		stringList[i] = auto.GetName()
	}

	dialogContent := container.NewVBox(
		widget.NewLabel("Selecciona una opción"),
	)

	if len(stringList) > 0 {
		dropdown := widget.NewSelect(stringList, func(selected string) {
			position := -1
			for i, item := range stringList {
				if item == selected {
					position = i
					break
				}
			}

			if position >= 0 && u.controller.SelectAutomata(position) {
				fmt.Printf("Posición seleccionada: %d\n", position)
				u.updateImageFunc()
			} else {
				fmt.Println("Elemento no encontrado en la lista")
			}

			fmt.Printf("Opción seleccionada: %s\n", selected)

		})
		dialogContent.Add(dropdown)
	}

	dialog := dialog.NewCustom("Desplegable", "", dialogContent, w)

	closeButton := widget.NewButton("Cerrar", func() {
		dialog.Hide()
	})
	dialogContent.Add(closeButton)

	dialog.Show()
}

func (u *UI) showStateConfigDialog(w fyne.Window) {
	dialogContent := container.NewVBox(
		widget.NewLabel("Configurar estado"),
	)

	dataEntry := widget.NewEntry()
	dataEntry.SetPlaceHolder("data")
	dialogContent.Add(dataEntry)

	initialCheckbox := widget.NewCheck("Inicial", nil)
	dialogContent.Add(initialCheckbox)

	finalCheckbox := widget.NewCheck("Final", nil)
	dialogContent.Add(finalCheckbox)

	dialog := dialog.NewCustom("Configurar estado", "", dialogContent, w)

	acceptButton := widget.NewButton("Aceptar", func() {
		data := dataEntry.Text
		initial := initialCheckbox.Checked
		final := finalCheckbox.Checked
		if u.controller.CreateState(data, initial, final) {
			fmt.Printf("Data: %s, Inicial: %t, Final: %t\n", data, initial, final)
			u.updateImageFunc()
			dialog.Hide()
		}
	})

	cancelButton := widget.NewButton("Cancelar", func() {
		dialog.Hide()
	})

	buttons := container.NewHBox(acceptButton, cancelButton)
	dialogContent.Add(buttons)

	dialog.Show()
}

func (u *UI) showTransitionConfigDialog(w fyne.Window) {
	dialogContent := container.NewVBox(
		widget.NewLabel("Configurar transición"),
	)

	startEntry := widget.NewEntry()
	startEntry.SetPlaceHolder("inicio")
	dialogContent.Add(startEntry)

	endEntry := widget.NewEntry()
	endEntry.SetPlaceHolder("termino")
	dialogContent.Add(endEntry)

	charsEntry := widget.NewEntry()
	charsEntry.SetPlaceHolder("caracteres (a,b,...)")
	dialogContent.Add(charsEntry)

	dialog := dialog.NewCustom("Configurar estado", "", dialogContent, w)

	acceptButton := widget.NewButton("Aceptar", func() {
		start := startEntry.Text
		end := endEntry.Text
		chars := charsEntry.Text
		if u.controller.CreateTransition(start, end, chars) {
			fmt.Printf("Inicio: %s, Termino: %s, Chars: %s\n", start, end, chars)
			u.updateImageFunc()
			dialog.Hide()
		}

	})

	cancelButton := widget.NewButton("Cancelar", func() {
		dialog.Hide()
	})

	buttons := container.NewHBox(acceptButton, cancelButton)
	dialogContent.Add(buttons)

	dialog.Show()
}

func (u *UI) showFileChooserDialog(win fyne.Window) {
	fileDialog := dialog.NewFileOpen(
		func(file fyne.URIReadCloser, err error) {
			if err == nil && file != nil {
				defer file.Close()

				if filepath.Ext(file.URI().Path()) == ".txt" {
					data, err := ioutil.ReadAll(file)
					if err != nil {
						log.Printf("No se pudo leer el archivo: %v\n", err)
						return
					}

					lines := strings.Split(string(data), "\n")

					// Guardar las cadenas en el controlador
					u.controller.SetInputStrings(lines)

					// Llamar a displayInputStrings para mostrar las cadenas
					u.displayInputStrings(win)
					// u.createStringCards(u.controller.GetInputStrings())

				} else {
					dialog.ShowInformation("Error", "Solo se permiten archivos .txt", win)
				}
			}
		},
		win,
	)
	predefinedPath := controllers.AbsInputTextFile
	location, err := storage.ListerForURI(storage.NewFileURI(predefinedPath))
	if err != nil {
		log.Printf("No se pudo convertir la ruta en ListableURI: %v\n", err)
		return
	}

	fileDialog.SetLocation(location)
	fileDialog.SetFilter(storage.NewExtensionFileFilter([]string{".txt"}))
	fileDialog.Show()
}

func (u *UI) createStringCards(inputStrings []string) *fyne.Container {
	cards := make([]fyne.CanvasObject, len(inputStrings))

	for i, inputString := range inputStrings {
		// Crear una etiqueta para la cadena
		label := widget.NewLabel(inputString)

		// Crear un contenedor vacío para el ícono
		iconContainer := container.NewMax()

		// Crear un contenedor horizontal que tenga la etiqueta y el contenedor de íconos
		cardContent := container.NewHBox(label, iconContainer)

		// Crear un borde alrededor de la card
		border := canvas.NewRectangle(color.Black)
		border.SetMinSize(fyne.NewSize(
			cardContent.MinSize().Width+50, cardContent.MinSize().Height+20,
		))

		// Crear un contenedor para la card que incluye el borde y el contenido
		card := container.NewMax(border, cardContent)
		cards[i] = card
	}

	// Crear un contenedor vertical que contenga todas las cards
	stringCards := container.NewVBox(cards...)

	return stringCards
}

func (u *UI) displayInputStrings(w fyne.Window) {
	inputStrings := u.controller.GetInputStrings()

	// Llamar a createStringCards para crear las tarjetas
	stringCards := u.createStringCards(inputStrings)

	// Agregar el contenedor a la UI
	if u.inputStringsContainer == nil {
		u.inputStringsContainer = stringCards
		mainContainer := w.Content().(*fyne.Container)
		mainContainer.Add(stringCards)
		mainContainer.Refresh()
	} else {
		u.inputStringsContainer.Objects = stringCards.Objects
		u.inputStringsContainer.Refresh()
	}
}

func (u *UI) validateInputStrings() {
	validations, error := u.controller.ProcessInputStrings()
	fmt.Println("validations: ", validations)

	if u.inputStringsContainer == nil || len(u.inputStringsContainer.Objects) == 0 {
		log.Println("No hay íconos para actualizar")
		return
	}

	if error != nil {
		fmt.Println("Vals error?", error.Error())
		return
	}

	for i, card := range u.inputStringsContainer.Objects {
		cardContainer := card.(*fyne.Container) // Obtener el contenedor de la tarjeta

		border := cardContainer.Objects[0].(*canvas.Rectangle) // Obtener el objeto borde (rectangle)

		border.FillColor = color.RGBA{R: 209, G: 93, B: 35, A: 255} // Rojo
		if validations[i] {
			border.FillColor = color.RGBA{R: 119, G: 209, B: 35, A: 255} // Verde
		}

		content := cardContainer.Objects[1].(*fyne.Container) // Obtener el contenedor de contenido (etiqueta y contenedor de íconos)
		iconContainer := content.Objects[1].(*fyne.Container) // Obtener el contenedor del ícono

		// Dependiendo de si la validación es verdadera o falsa, agregar el ícono correspondiente
		if validations[i] {
			icon := canvas.NewImageFromResource(theme.ConfirmIcon()) // Ícono de validación
			iconContainer.Objects = []fyne.CanvasObject{icon}
		} else {
			icon := canvas.NewImageFromResource(theme.CancelIcon()) // Ícono de cancelación
			iconContainer.Objects = []fyne.CanvasObject{icon}
		}

		border.Refresh()
		iconContainer.Refresh()
	}
}

// func (u *UI) displayInputStrings(w fyne.Window) {
// 	inputStrings := u.controller.GetInputStrings()

// 	// Crear una lista de objetos Canvas para las cadenas y los círculos
// 	items := make([]fyne.CanvasObject, 0, len(inputStrings)*2)

// 	for _, str := range inputStrings {
// 		// Crear un marco para la cadena
// 		frame := widget.NewLabel(str)

// 		// Crear un círculo gris
// 		circle := canvas.NewCircle(color.Gray{})

// 		circle.StrokeWidth = 2
// 		circle.StrokeColor = color.Gray{}

// 		// Agregar la cadena y el círculo a la lista de objetos
// 		items = append(items, frame, circle)
// 	}

// 	// Crear un encabezado "Cadenas"
// 	header := widget.NewLabel("Cadenas")
// 	items = append([]fyne.CanvasObject{header}, items...)

// 	// Crear un contenedor con los objetos
// 	inputStringsContainer := container.NewVBox(items...)

// 	// Agregar el contenedor a la UI
// 	if u.inputStringsContainer == nil {
// 		u.inputStringsContainer = inputStringsContainer
// 		mainContainer := w.Content().(*fyne.Container)
// 		mainContainer.Add(inputStringsContainer)
// 		mainContainer.Refresh()
// 	} else {
// 		u.inputStringsContainer.Objects = inputStringsContainer.Objects
// 		u.inputStringsContainer.Refresh()
// 	}
// }
