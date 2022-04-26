package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/google/uuid"

	"github.com/kepkin/interview-price-monitor/api"
)

func makeShortMonitorReq(endpoint string) error {
	payload := api.MonitorTask {
		Frequency: "1s",
		Start: time.Now().Add(time.Second*5),
		Stop:  time.Now().Add(time.Minute),
		Target: "bitcoin",
	}

	id := uuid.New().String()

	payloadBytes, err := json.Marshal(payload)

	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", endpoint + id, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(bodyBytes))
		return fmt.Errorf("returned statucs code: %v", resp.StatusCode)
	}

	fmt.Println("New task: ", id)
	return nil
}

func main() {
	var args struct {
		Server string
	}
	arg.MustParse(&args)
	err := makeShortMonitorReq(args.Server)
	if err != nil {
		log.Fatal(err)
	}
}
