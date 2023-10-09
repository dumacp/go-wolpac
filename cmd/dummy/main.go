package main

import (
	"context"
	"fmt"
	"log"
)

func main() {

	data := []byte(`@
!CFXXXXXXXXXXXX

!I1

!I1
!TMP





`)

	s := NewSerial(data)

	conf := NewConf(s)

	fmt.Printf("initial conf: %q\n", conf.OptsToString())

	dev, err := conf.Open()
	if err != nil {
		log.Fatalln(err)
	}

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	ch := dev.Events(ctx)

	for v := range ch {
		fmt.Printf("event: %v\n", v)
	}

}
