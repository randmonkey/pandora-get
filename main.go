package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/randmonkey/pandora-get/pandora"
)

func main() {
	var server, token string
	var spl string
	var splFile string

	flag.StringVar(&server, "s", "", "pandora server")
	flag.StringVar(&token, "t", "", "pandora token")
	flag.StringVar(&splFile, "f", "", "SPL file, empty to read from stdin")
	flag.Parse()

	if splFile != "" {
		f, err := os.Open(splFile)
		if err != nil {
			log.Fatalf("failed to open SPL file %s, error %v", splFile, err)
		}
		buf, err := ioutil.ReadAll(f)
		if err != nil {
			log.Fatalf("failed to read SPL file %s, error %v", splFile, err)
		}
		spl = string(buf)
	}

	interval := 10 * time.Minute
	for {
		pandoraClient := pandora.NewClient(server, token)
		now := time.Now()
		res, err := pandoraClient.GetQueryResult(spl, now.Add(-1*interval), now, 10000, 60*time.Second)
		if err != nil {
			log.Printf("failed to get result from pandora, error %+v", err)
		}
		if len(res) > 0 {
			buf, _ := json.Marshal(res[0])
			fmt.Println(string(buf))
		}

		time.Sleep(interval)
	}

}
