package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"bufio"
	"log"
	"strings"
	"encoding/hex"
	"net/http"
)

func main() {
	// variables declaration
	var mac string
	var url string
	var job string
	var metric string
	var verbose bool

	// flags declaration using flag package
	flag.StringVar(&mac, "b", "", "Bluetooth mac address used with gatttool to connect and parse data")
	flag.StringVar(&url, "u", "", "pushgateway url to push bpm to prometheus")
	flag.StringVar(&job, "j", "cardiac_frequency", "Specify prometheus job.")
	flag.StringVar(&metric, "m", "bpm", "Specify prometheus metric.")
	flag.BoolVar(&verbose, "v", false, "verbose")

	flag.Parse()  // after declaring flags we need to call it

	if mac == "" || url == "" {
		flag.Usage()
		os.Exit(0)
	}

	if verbose {
		fmt.Print(url, "/metrics/job/", job, "\n")
	}

	//gatttool -t random -b 01:AB:CD:EF:02:03 --char-write-req --handle=0x0011 --value=0100 --listen
	if verbose {
		fmt.Println("gatttool -t random -b", mac, "--char-write-req --handle=0x0011 --value=0100 --listen")
	}

	cmd := exec.Command("/usr/bin/gatttool", "-t", "random", "-b", mac, "--char-write-req", "--handle=0x0011", "--value=0100", "--listen")

	stdout, _ := cmd.StdoutPipe()
	cmd.Start()

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		raw_text := scanner.Text()
		s := strings.SplitAfter(raw_text, ": ")
		if len(s) > 1 {
			n := strings.Fields(s[len(s)-1])
			bpm, err := hex.DecodeString(n[1])
			if err != nil {
				log.Println(err)
			}

			// display logs
			if verbose {
				fmt.Println(int(bpm[0]))
			}

			metric_str := strings.NewReader(fmt.Sprintln(metric, int(bpm[0])))

			if verbose {
				fmt.Print(url, "/metrics/job/", job, "\n")
			}
			//echo "some_metric 4.16" | curl --data-binary @- http://192.168.1.2:9091/metrics/job/some_job
			req, err := http.NewRequest("POST", fmt.Sprintf("%s/metrics/job/%s", url, job), metric_str)
			if err != nil {
				panic(err)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				panic(err)
			}
			if verbose {
				fmt.Println(resp.Status)
			}
			defer resp.Body.Close()
		}
	}

	if err := scanner.Err(); err != nil {
			log.Println(err)
	}
}
