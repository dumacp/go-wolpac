package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/dumacp/go-wolpac/gpiosysfs"
)

var (
	gpioInput1 int
	gpioInput2 int
	pathSignal string
)

func init() {
	flag.IntVar(&gpioInput1, "gpioinput1", 85, "número de la GPIO para sensor 1")
	flag.IntVar(&gpioInput2, "gpioinput2", 86, "número de la GPIO para sensor 2")
	flag.StringVar(&pathSignal, "pathsignal", "output1", "path (/sys/class/leds/<PATH>) GPIO para enviar el comando de permtir paso")
}

func main() {
	flag.Parse()

	cmdEnable := fmt.Sprintf("echo 1 > /sys/class/leds/%s/brightness", pathSignal)
	cmdDisable := fmt.Sprintf("echo 0 > /sys/class/leds/%s/brightness", pathSignal)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pin1, err := gpiosysfs.OpenPin(gpioInput1)
	if err != nil {
		log.Fatalln(err)
	}
	pin2, err := gpiosysfs.OpenPin(gpioInput2)
	if err != nil {
		log.Fatalln(err)
	}

	ch1, err := pin1.EpollEvents(ctx, true, gpiosysfs.Falling)
	if err != nil {
		log.Fatalln(err)
	}
	ch2, err := pin2.EpollEvents(ctx, true, gpiosysfs.Falling)
	if err != nil {
		log.Fatalln(err)
	}

	tick := time.NewTicker(10 * time.Second)
	defer tick.Stop()

	activeStep := false
	funcCommand := func(cmd string) ([]byte, error) {
		if out, err := exec.Command("/bin/sh", "-c", cmd).Output(); err != nil {
			return out, err
		} else if len(out) > 0 {
			return out, nil
		}
		return nil, nil
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
			if out, err := funcCommand(cmdEnable); err != nil {
				fmt.Printf("error comand: %q, err: %s, output: %s\n", cmdEnable, err, out)
			}
			activeStep = true
		case v, ok := <-ch1:
			if !ok {
				fmt.Println("ch1 close")
				return
			}
			if activeStep {
				activeStep = false
				if out, err := funcCommand(cmdDisable); err != nil {
					fmt.Printf("error comand: %q, err: %s, output: %s\n", cmdDisable, err, out)
				}
			}
			fmt.Printf("input: %d, rising: %v, timestamp: %s\n", v.Pin, v.RisingEdge, v.Time)
		case v, ok := <-ch2:
			if !ok {
				fmt.Println("ch2 close")
				return
			}
			if activeStep {
				activeStep = false
				if out, err := funcCommand(cmdDisable); err != nil {
					fmt.Printf("error comand: %q, err: %s, output: %s\n", cmdDisable, err, out)
				}
			}
			fmt.Printf("input: %d, rising: %v, timestamp: %s\n", v.Pin, v.RisingEdge, v.Time)
		}
	}
}
