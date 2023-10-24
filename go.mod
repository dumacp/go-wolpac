module github.com/dumacp/go-wolpac

go 1.19

replace github.com/asynkron/protoactor-go => ../../asynkron/protoactor-go

require (
	github.com/eiannone/keyboard v0.0.0-20220611211555-0d226195f203
	github.com/looplab/fsm v1.0.1
	github.com/nsf/termbox-go v1.1.1
	github.com/tarm/serial v0.0.0-20180830185346-98f6abe2eb07
	golang.org/x/sys v0.12.0
)

require github.com/mattn/go-runewidth v0.0.9 // indirect

replace github.com/brian-armstrong/gpio => ../../brian-armstrong/gpio
