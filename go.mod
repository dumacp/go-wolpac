module github.com/dumacp/go-wolpac

go 1.19

replace github.com/asynkron/protoactor-go => ../../asynkron/protoactor-go

require (
	github.com/brian-armstrong/gpio v0.0.0-00010101000000-000000000000
	github.com/nsf/termbox-go v1.1.1
	github.com/tarm/serial v0.0.0-20180830185346-98f6abe2eb07
	periph.io/x/conn/v3 v3.7.0
)

require (
	github.com/eiannone/keyboard v0.0.0-20220611211555-0d226195f203 // indirect
	github.com/mattn/go-runewidth v0.0.9 // indirect
	golang.org/x/sys v0.12.0 // indirect
)

replace github.com/brian-armstrong/gpio => ../../brian-armstrong/gpio
