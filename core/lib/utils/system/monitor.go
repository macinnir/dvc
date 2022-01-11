package system

import (
	"fmt"
	"time"
)

type Monitor struct {
	interval     time.Duration
	counter      int
	pings        chan struct{}
	ticker       *time.Ticker
	pingCallback func(counter int)
}

func NewMonitor(seconds int) *Monitor {

	fmt.Println("NewMonitor()")

	m := &Monitor{
		pings: make(chan struct{}),
		pingCallback: func(counter int) {
			fmt.Println("Sending ping #", counter)
		},
		interval: time.Second * time.Duration(seconds),
	}

	// Don't start the monitor unless the seconds interval is specified
	if m.interval > 0 {
		m.ticker = time.NewTicker(m.interval)
		go m.Start()
	}

	return m
}

func (m *Monitor) SetPingCallback(cb func(counter int)) {
	m.pingCallback = cb
}

func (m *Monitor) ping() {
	fmt.Println("Ping()")
	// Send struct to pings
	m.pings <- struct{}{}
}

func (m *Monitor) Start() {

	fmt.Println("Start()")

	// interrupt := make(chan os.Signal, 1)
	// signal.Notify(interrupt, os.Interrupt)

	// loop:
	for {
		select {
		case t := <-m.ticker.C:
			fmt.Println("Tick at ", t.String())
			go m.ping()
		case <-m.pings:
			m.counter++
			m.pingCallback(m.counter)

			// case <-interrupt:
			// 	fmt.Println("Interrupt!!!!")
			// 	break loop
		}
	}

	// fmt.Println("Got here!")
}

func (m *Monitor) Stop() {

	fmt.Println("Stopping the monitor...")

	if m.ticker != nil {
		m.ticker.Stop()
	}

	fmt.Println("monitor stopped!")

}
