package main

import (
	"fmt"
	"time"

	"github.com/nsf/termbox-go"
)

var (
	selectedIndex = 1
	menuItems     = []string{
		"Hora actual:", // Esta opción se actualiza dinámicamente
		"Opción 1",
		"Opción 2",
		"Salir",
	}
)

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	go drawMenu()

	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

loop:
	for {
		select {
		case ev := <-eventQueue:
			switch ev.Type {
			case termbox.EventKey:
				switch ev.Key {
				case termbox.KeyArrowDown:
					if selectedIndex < len(menuItems)-1 {
						selectedIndex++
					}
				case termbox.KeyArrowUp:
					if selectedIndex > 1 {
						selectedIndex--
					}
				case termbox.KeyEnter:
					handleSelection(selectedIndex)
					if selectedIndex == len(menuItems)-1 { // Si seleccionó "Salir"
						break loop
					}
				}
				drawMenu()
			}
		case <-time.Tick(1 * time.Second):
			drawMenu()
		}
	}
}

func drawMenu() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	menuItems[0] = fmt.Sprintf("Hora actual: %s", time.Now().Format("15:04:05"))
	printMenu(menuItems)
	termbox.Flush()
}

func printMenu(items []string) {
	for y, line := range items {
		color := termbox.ColorDefault
		if y == selectedIndex {
			color = termbox.ColorGreen // Resaltar opción seleccionada
		}
		for x, ch := range line {
			termbox.SetCell(x, y*2, ch, color, termbox.ColorDefault)
		}
	}
}

func handleSelection(index int) {
	switch index {
	case 1:
		showMessage("Seleccionaste la Opción 1!")
	case 2:
		showMessage("Seleccionaste la Opción 2!")
	case 3:
		showMessage("Adiós!")
	}
	if index != 3 { // Si no seleccionó "Salir"
		drawMenu() // Volver a mostrar el menú
	}
}

func showMessage(message string) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	for x, ch := range message {
		_, h := termbox.Size()
		termbox.SetCell(x, h/2, ch, termbox.ColorWhite, termbox.ColorDefault)
	}
	termbox.Flush()
	time.Sleep(2 * time.Second)
}
