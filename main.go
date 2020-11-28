package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joelbirchler/yupa/pkg/pms5003"
	"github.com/joelbirchler/yupa/pkg/rgbmatrix"
	"gobot.io/x/gobot/platforms/raspi"
	"k8s.io/apimachinery/pkg/util/json"
)

type adafruitIo struct {
	username, key, group string
}

type feed struct {
	Key   string `json:"key"`
	Value uint16 `json:"value"`
}

type createData struct {
	Feeds []feed `json:"feeds"`
}

func (a adafruitIo) send(pm25, pm100, e uint16) error {
	client := http.Client{
		Timeout: time.Second * 5,
	}

	url := fmt.Sprintf("https://io.adafruit.com/api/v2/%s/groups/%s/data", a.username, a.group)

	foo := createData{
		Feeds: []feed{
			{"Environment25", pm25},
			{"Environment100", pm100},
			{"Error", e},
		},
	}

	data, err := json.Marshal(foo)
	if err != nil {
		return err
	}

	log.Printf("obj: %v", foo)
	log.Printf("json: %s", data)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-AIO-Key", a.key)

	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	log.Println(string(body)) // FIXME: we can ignore the body

	return nil
}

var aio adafruitIo

func init() {
	aio = adafruitIo{os.Getenv("AIO_USERNAME"), os.Getenv("AIO_KEY"), os.Getenv("AIO_GROUP")}
	if aio.username == "" || aio.key == "" || aio.group == "" {
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

		if err := aio.send(f.Environment25, f.Environment100, errorCount); err != nil {
			log.Printf("sending error: %v", err)
		}

		rgbmatrix.Set(rgbmatrix.Binary(f.Environment25))
		if err := rgbmatrix.Render(); err != nil {
			log.Printf("rendering error: %v", err)
		}

		time.Sleep(time.Second * 30)
	}
}
