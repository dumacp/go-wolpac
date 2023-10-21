package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/dumacp/go-wolpac/pwacii"
)

var (
	portpath string
	timeout  int
)

func init() {
	flag.StringVar(&portpath, "port", "/dev/ttyUSB0", "path to serial port device")
	flag.IntVar(&timeout, "timeout", 10, "send comannd timeout in second")
}

func main() {

	flag.Parse()
	conf := pwacii.NewConf(portpath)

	fmt.Printf("default Opts: %v\n", conf)

	dev, err := conf.Open()
	if err != nil {
		log.Fatalln(err)
	}

	tick0 := time.NewTimer(5 * time.Second)
	defer tick0.Stop()
	tick := time.NewTicker(time.Duration(timeout) * time.Second)
	defer tick.Stop()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	evts := dev.Events(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case v, ok := <-evts:
			if !ok {
				fmt.Println("close channel")
				return
			}
			fmt.Printf("event: %q\n", v)
		case <-tick0.C:
			v, err := dev.Command(pwacii.OneEntryAllow, "")
			if err != nil {
				fmt.Printf("error command (%q): %s\n", pwacii.OneEntryAllow, err)
				break
			}
			time.Sleep(1 * time.Second)
			fmt.Printf("command (%q) response: %q\n", pwacii.OneEntryAllow, v)
			v1, err := dev.Command(pwacii.StatusRequest, "")
			if err != nil {
				fmt.Printf("error command (%q): %s\n", pwacii.OneEntryAllow, err)
				break
			}
			fmt.Printf("command (%q) response: %q\n", pwacii.OneEntryAllow, v1)
		case <-tick.C:
			v, err := dev.Command(pwacii.OneEntryAllow, "")
			if err != nil {
				fmt.Printf("error command (%q): %s\n", pwacii.OneEntryAllow, err)
				break
			}
			fmt.Printf("command (%q) response: %q\n", pwacii.OneEntryAllow, v)
			time.Sleep(1 * time.Second)
			fmt.Printf("command (%q) response: %q\n", pwacii.OneEntryAllow, v)
			v1, err := dev.Command(pwacii.StatusRequest, "")
			if err != nil {
				fmt.Printf("error command (%q): %s\n", pwacii.StatusRequest, err)
				break
			}
			fmt.Printf("command (%q) response: %q\n", pwacii.StatusRequest, v1)
		}
	}
}
