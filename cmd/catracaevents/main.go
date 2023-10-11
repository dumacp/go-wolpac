package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/dumacp/go-wolpac/catraca"
)

var (
	gpioInput1 int
	gpioInput2 int
	pathSignal string
)

func init() {
	flag.IntVar(&gpioInput1, "gpioinput1", 86, "número de la GPIO para sensor 1")
	flag.IntVar(&gpioInput2, "gpioinput2", 85, "número de la GPIO para sensor 2")
	flag.StringVar(&pathSignal, "pathsignal", "output1", "path (/sys/class/leds/<PATH>) GPIO para enviar el comando de permtir paso")
}

func main() {
	flag.Parse()

	fmt.Printf("default opts: %s\n", catraca.DefaultsOptions())
	dev := catraca.New(catraca.WithInputSysfsT1(gpioInput1), catraca.WithInputSysfsT2(gpioInput2))

	fmt.Printf("actual opts: %s\n", dev.Opts)

	if err := dev.Open(); err != nil {
		log.Fatalf("open error: %s", err)
	}
	defer dev.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ch, err := dev.Events(ctx, true)
	if err != nil {
		log.Fatalf("events error: %s", err)
	}

	tick := time.NewTicker(10 * time.Second)
	defer tick.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
			fmt.Println("send allow Entrance")
			if err := dev.OneEntrance(); err != nil {
				log.Printf("error allow: %s", err)
				return
			}
		case v := <-ch:
			log.Printf("new event: %v", v)
		}
	}
}
