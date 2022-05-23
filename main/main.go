package main

import (
	"flag"
	"fmt"
	"net/http"
	"io/ioutil"
	"os"
	"strings"
	"time"
	"encoding/json"
)


type Response struct {
	VersionHash string `json:"Version hash"`
	BuildDate   string `json:"Build date"`
}

func main() {
	var url = flag.String("url", "http://localhost/", "URL to poll")
	var responseCode = flag.Int("code", 200, "Response code to wait for")
	var timeout = flag.Int("timeout", 2000, "Timeout before giving up in ms")
	var interval = flag.Int("interval", 200, "Interval between polling in ms")
	var localhost = flag.String("localhost", "", "Ip address to use for localhost")
	var gitHash = flag.String("gitHash", "", "githash for build")

	flag.Parse()

	fmt.Printf("Polling URL `%s` Git hash `%s` for response code %d for up to %d ms at %d ms intervals\n", *url, *gitHash ,*responseCode, *timeout, *interval)
	startTime := time.Now()
	timeoutDuration := time.Duration(*timeout) * time.Millisecond
	sleepDuration := time.Duration(*interval) * time.Millisecond

	if *localhost!="" && strings.Contains(*url, "localhost") {
		*url = strings.ReplaceAll(*url, "localhost", *localhost)
	}
	for {
		res, err := http.Get(*url)


        body, err := ioutil.ReadAll(res.Body) // response body is []byte
         fmt.Println(string(body))
         // snippet only
         var result Response
         if err := json.Unmarshal(body, &result); err != nil {   // Parse []byte to go struct pointer
             fmt.Println("Can not unmarshal JSON")
         }
         fmt.Printf("Hash: %s, Data: %s", result.VersionHash, result.BuildDate)

		if err == nil && res.StatusCode == *responseCode &&  result.VersionHash == *gitHash {
			fmt.Printf("Response header: %v", res)
			os.Exit(0)
		}
		time.Sleep(sleepDuration)
		elapsed := time.Now().Sub(startTime)
		if elapsed > timeoutDuration {
			fmt.Printf("Timed out\n")
			os.Exit(1)
		}
	}
}
// PrettyPrint to print struct in a readable way
func PrettyPrint(i interface{}) string {
    s, _ := json.MarshalIndent(i, "", "\t")
    return string(s)
}