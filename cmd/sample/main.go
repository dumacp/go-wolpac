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
)

func init() {
	flag.StringVar(&portpath, "port", "/dev/ttyUSB2", "path to serial port device")
}

func main() {

	flag.Parse()
	conf := pwacii.NewConf(portpath)

	fmt.Printf("default Opts: %v\n", conf)

	dev, err := conf.Open()
	if err != nil {
		log.Fatalln(err)
	}

	tick := time.NewTicker(10 * time.Second)
	defer tick.Stop()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	evts := dev.Events(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case v := <-evts:
			fmt.Printf("event: %q\n", v)
		case <-tick.C:
			v, err := dev.Command(pwacii.OneEntryAllow, "")
			if err != nil {
				fmt.Printf("error command (%q): %s\n", pwacii.OneEntryAllow, err)
				break
			}
			fmt.Printf("command (%q) response: %q\n", pwacii.OneEntryAllow, v)
		}
	}
}
