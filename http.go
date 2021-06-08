package main

import (
	"log"
	"net/http"
	"strings"
)

func httpPush(stringReader *strings.Reader) error {
	resp, err := http.Post(pushJobURL, "application/x-www-form-urlencoded", stringReader)
	if err != nil {
		return (err)
	}
	defer resp.Body.Close()
	if verbose {
		log.Println(resp.Status)
	}
	return nil
}
