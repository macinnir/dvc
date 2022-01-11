package system

import (
	"os"
	"os/signal"
	"syscall"
)

type InterruptHandler struct {
	signalChan chan os.Signal
}

func NewInterruptHandler() *InterruptHandler {
	i := &InterruptHandler{
		signalChan: make(chan os.Signal),
	}
	signal.Notify(i.signalChan, os.Interrupt, syscall.SIGTERM)
	return i
}

func (i *InterruptHandler) OnInterrupt(cb func()) {
	select {
	case <-i.signalChan:
		cb()
	}
}
