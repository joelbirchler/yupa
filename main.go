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

type adafruitIo struct {
	username, key, group string
}

func init() {
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
	aio := adafruitIo{os.Getenv("AIO_USERNAME"), os.Getenv("AIO_KEY"), os.Getenv("AIO_GROUP")}
	if aio.username == "" || aio.key == "" || aio.group == "" {
		log.Fatalln("Environment variables AIO_USERNAME, AIO_KEY, and AIO_GROUP are required. See the README for more information.")
	}

	pi := raspi.NewAdaptor()

	if err := rgbmatrix.Open(pi); err != nil {
		log.Fatalf("unable to setup rgbmatrix: %v", err)
	}

	if err := pms5003.Open(pi); err != nil {
		log.Fatalf("unable to setup pms5003: %v", err)
	}

	for {
		f, err := pms5003.ReadFrame()
		if err != nil {
			log.Printf("error reading frame: %v", err)
			// continue
		}

		log.Printf("%+v", f)

		rgbmatrix.Set(rgbmatrix.Binary(f.Environment25))
		if err := rgbmatrix.Render(); err != nil {
			log.Printf("rendering error: %v", err)
		}

		time.Sleep(time.Second)
	}
}
