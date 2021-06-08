package main

import (
	"bufio"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

var (
	pushJobURL string
	verbose    bool
)

func main() {
	// variables declaration
	var mac string
	var url string
	var job string
	var metric string

	// flags declaration using flag package
	flag.StringVar(&mac, "b", "", "Bluetooth mac address used with gatttool to connect and parse data")
	flag.StringVar(&url, "u", "", "pushgateway url to push bpm to prometheus")
	flag.StringVar(&job, "j", "cardiac_frequency", "Specify prometheus job.")
	flag.StringVar(&metric, "m", "bpm", "Specify prometheus metric.")
	flag.BoolVar(&verbose, "v", false, "verbose")

	flag.Parse() // after declaring flags we need to call it

	if mac == "" || url == "" {
		flag.Usage()
		os.Exit(1)
	}

	pushJobURL = fmt.Sprintf("%s/metrics/job/%s", url, job)

	if verbose {
		log.Println("Will push to:", pushJobURL)
	}

	//gatttool -t random -b 01:AB:CD:EF:02:03 --char-write-req --handle=0x0011 --value=0100 --listen
	command := fmt.Sprintf("gatttool -t random -b %s --char-write-req --handle=0x0011 --value=0100 --listen", mac)
	if verbose {
		log.Println("Fetch command:", command)
	}
	commandSplitted := strings.Fields(command)
	cmd := exec.Command(commandSplitted[0], commandSplitted[1:]...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		rawText := scanner.Text()
		s := strings.SplitAfter(rawText, ": ")
		if len(s) > 1 {
			n := strings.Fields(s[len(s)-1])
			bpm, err := hex.DecodeString(n[1])
			if err != nil {
				log.Println(err)
				continue
			}

			// display logs
			if verbose {
				log.Println(int(bpm[0]))
			}

			metricReader := strings.NewReader(fmt.Sprintf("%s %d\n", metric, int(bpm[0])))

			if verbose {
				log.Println(pushJobURL)
			}
			//echo "some_metric 4.16" | curl --data-binary @- http://192.168.1.2:9091/metrics/job/some_job
			if err := httpPush(metricReader); err != nil {
				log.Fatal(err)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
}
