package lib

import (
	"fmt"
	"log"
	"time"
)

func LogAdd(t time.Time, message string, a ...interface{}) {
	fmt.Printf("[+] %fs %s\n", time.Since(t).Seconds(), fmt.Sprintf(message, a...))
}

func LogRemove(t time.Time, message string, a ...interface{}) {
	fmt.Printf("[-] %fs %s\n", time.Since(t).Seconds(), fmt.Sprintf(message, a...))
}

func LogFatal(e error) {
	log.Fatalf("[!] %s\n", e.Error())
}
