package adafruitio

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"k8s.io/apimachinery/pkg/util/json"
)

// FeedSet is a group of feeds with username and key fields for authentication
// https://io.adafruit.com/api/docs/#groups
type FeedSet struct {
	Username, Key, Group string
}

type feed struct {
	Key   string `json:"key"`
	Value uint16 `json:"value"`
}

type createData struct {
	Feeds []feed `json:"feeds"`
}

// Send sends pm data to adafruit.io
func (f FeedSet) Send(pm25, pm100, e uint16) error {
	client := http.Client{
		Timeout: time.Second * 5,
	}

	url := fmt.Sprintf("https://io.adafruit.com/api/v2/%s/groups/%s/data", f.Username, f.Group)

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

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-AIO-Key", f.Key)

	if _, err := client.Do(req); err != nil {
		return err
	}
	return nil
}
