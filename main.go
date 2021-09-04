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

const sleepSeconds = 60
const readsBeforeReset = 60

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

	var errorCount, reads uint16

	for {
		// it takes ~40 seconds for the PMS5003 to be accurrate after startup
		time.Sleep(time.Second * sleepSeconds)

		// read
		f, err := pms5003.ReadFrame()
		if err != nil {
			log.Printf("error reading frame: %v", err)
			errorCount++
		}

		log.Printf("%+v", f)

		// send for telemetry
		if err := aio.Send(f.Environment25, f.Environment100, errorCount); err != nil {
			log.Printf("sending error: %v", err)
		}

		// show on the display
		rgbmatrix.Set(rgbmatrix.Fill(aqi2rgb(f.Environment25), 0x20))
		if err := rgbmatrix.Render(); err != nil {
			log.Printf("rendering error: %v", err)
		}

		// reset
		if reads > readsBeforeReset {
			if err := pms5003.Reset(pi); err != nil {
				log.Printf("resetting error: %v", err)
			}
			reads = 0
		}

		reads++
	}
}

// This a simplistic version of AQI that only uses one datapoint. Typically AQI is
// calculated with multiple datapoints averaged over a 24-hour period.
func aqi2rgb(u uint16) rgbmatrix.Rgb {
	switch {
	case u < 13: // Good
		return rgbmatrix.Rgb{R: 0x00, G: 0xFF, B: 0x00}
	case u < 35: // Moderate
		return rgbmatrix.Rgb{R: 0xFF, G: 0x90, B: 0x00}
	case u < 55: // Unhealthy for senstive groups
		return rgbmatrix.Rgb{R: 0xFF, G: 0x30, B: 0x00}
	case u < 150: // Unhealthy
		return rgbmatrix.Rgb{R: 0xDD, G: 0x00, B: 0x00}
	case u < 250: // Very Unhealthy
		return rgbmatrix.Rgb{R: 0xBB, G: 0x00, B: 0x33}
	}

	// Hazardous
	return rgbmatrix.Rgb{R: 0x33, G: 0x00, B: 0x33}
}
