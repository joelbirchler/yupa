package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joelbirchler/yupa/pkg/pms5003"
	"github.com/joelbirchler/yupa/pkg/rgbmatrix"
	"gobot.io/x/gobot/platforms/raspi"
)

func init() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		rgbmatrix.Close()
		pms5003.Close()
		os.Exit(1)
	}()
}

func main() {
	pi := raspi.NewAdaptor()

	if err := rgbmatrix.Open(pi); err != nil {
		log.Fatalf("unable to setup rgbmatrix: %v", err)
	}

	for i := 0; i < 600; i++ {
		rgbmatrix.Set(rgbmatrix.Binary(uint16(i)))

		if err := rgbmatrix.Render(); err != nil {
			log.Printf("rendering error: %v", err)
		}

		time.Sleep(250 * time.Millisecond)
	}

	if err := pms5003.Open(pi); err != nil {
		log.Fatalf("unable to setup pms5003: %v", err)
	}

	for {
		f, err := pms5003.ReadFrame()
		if err != nil {
			log.Fatalf("error reading frame: %v", err)
		}

		log.Printf("%+v", f)

		time.Sleep(time.Second)
	}
}
