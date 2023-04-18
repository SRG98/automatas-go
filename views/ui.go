package views

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/SRG98/automatas-go/controllers"
)

type UI struct {
	controller      *controllers.Controller
	updateImageFunc func()
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

	// Crear 10 botones
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
				ui.showCreateAutomataDialog(win)
			}
		case 5:
			text = "Ver cadenas"
			onTapped = func() {
				ui.showCreateAutomataDialog(win)
			}
		case 6:
			text = "Procesar cadenas"
			onTapped = func() {
				ui.showCreateAutomataDialog(win)
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
