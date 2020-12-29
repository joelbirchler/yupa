package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joelbirchler/yupa/pkg/adafruitio"
	"github.com/joelbirchler/yupa/pkg/pms5003"
	"github.com/joelbirchler/yupa/pkg/rgbmatrix"
	"gobot.io/x/gobot/platforms/raspi"
)

var aio adafruitio.FeedSet

func init() {
	aio = adafruitio.FeedSet{Username: os.Getenv("AIO_USERNAME"), Key: os.Getenv("AIO_KEY"), Group: os.Getenv("AIO_GROUP")}
	if aio.Username == "" || aio.Key == "" || aio.Group == "" {
		log.Fatalln("Environment variables AIO_USERNAME, AIO_KEY, and AIO_GROUP are required. See the README for more information.")
	}

	// trap sigterm to attempt to close before exiting
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

	if err := pms5003.Open(pi); err != nil {
		log.Fatalf("unable to setup pms5003: %v", err)
	}

	var errorCount uint16

	for {
		f, err := pms5003.ReadFrame()
		if err != nil {
			log.Printf("error reading frame: %v", err)
			errorCount++
		}

		log.Printf("%+v", f)

		if err := aio.Send(f.Environment25, f.Environment100, errorCount); err != nil {
			log.Printf("sending error: %v", err)
		}

		rgbmatrix.Set(rgbmatrix.Binary(f.Environment25))
		if err := rgbmatrix.Render(); err != nil {
			log.Printf("rendering error: %v", err)
		}

		time.Sleep(time.Second * 30)
	}
}
