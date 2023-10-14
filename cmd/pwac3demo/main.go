package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/dumacp/go-wolpac/pwaciii"
	"github.com/eiannone/keyboard"
)

var (
	portpath string
)

func init() {
	flag.StringVar(&portpath, "port", "/dev/ttyUSB2", "path to serial port device")
}

func main() {
	flag.Parse()

	fmt.Printf("default opts: %s\n", pwaciii.DefaultsOptions())

	dev := pwaciii.New(pwaciii.WithPort(portpath), pwaciii.WithEstadoBloque(pwaciii.EntradacontroladaSalidaLibre))

	fmt.Printf("final opts: %s\n", dev.Opts)

	if err := dev.Open(); err != nil {
		log.Fatalln("open error: ", err)
	}

	tick0 := time.NewTimer(1 * time.Second)
	defer tick0.Stop()

	tickevents := time.NewTicker(3 * time.Second)
	defer tickevents.Stop()

	quit := make(chan struct{})
	enter := make(chan struct{})

	// Crea una goroutine para capturar la entrada del teclado
	go func() {
		err := keyboard.Open()
		if err != nil {
			log.Fatal(err)
		}
		defer keyboard.Close()

		for {
			_, key, err := keyboard.GetKey()
			if err != nil {
				log.Fatal(err)
			}
			if key == keyboard.KeyEnter {
				select {
				case enter <- struct{}{}:
				default:
				}
				break
			}
			if key == keyboard.KeyEsc {
				close(quit)
				return
			}
		}
	}()

	for {
		select {
		case <-quit:
			fmt.Println("quit")
			return
		case <-enter:
			fmt.Println("send entrance allow")
			if err := dev.OneEntrance(); err != nil {
				log.Println("error send allow: ", err)
			}
		case <-tickevents.C:
			if res, err := dev.SolicitaEstadoEntrada(); err != nil {
				log.Println("error info entrada: ", err)
			} else {
				log.Printf("informacion de entrada: %s\n", res)
			}
			if res, err := dev.SolicitaEstadoSalida(); err != nil {
				log.Println("error info salida: ", err)
			} else {
				log.Printf("informacion de salida: %s\n", res)
			}
			if in, out, err := dev.SolicitaContadores(); err != nil {
				log.Println("error contadores: ", err)
			} else {
				log.Printf("contadores: (in) %d,  (out) %d\n", in, out)
			}
		}
	}
}
