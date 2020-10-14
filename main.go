package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

var (
	flagOutput = flag.String("o", "", "output file")
	flagHeader = flag.Bool("header", false, "print HTTP-header")
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		log.Fatal("Please enter only one URL")
	}
	url := args[0]
	if !validateURL(url) {
		log.Fatalf("URL not valid: %s\n", url)
	}

	if resp, err := http.Get(url); err != nil {
		log.Fatal(err)
	} else {
		defer resp.Body.Close()

		var w io.Writer
		w = os.Stdout
		if *flagHeader {
			for k, v := range resp.Header {
				fmt.Fprintf(w, "%s :\n", k)
				for i, l := range v {
					fmt.Fprintf(w, "  %03d: %s \n", i+1, l)
				}
			}
			os.Exit(0)
		}

		if *flagOutput != "" {
			if err := os.MkdirAll(filepath.Dir(*flagOutput), 0755); err != nil {
				log.Fatal(err)
			}

			if f, err := os.OpenFile(*flagOutput, os.O_RDWR|os.O_CREATE, 0755); err != nil {
				log.Fatal(err)
			} else {
				defer f.Close()
				w = f
			}
		}

		io.Copy(w, resp.Body)
	}
}

func validateURL(s string) bool {
	_, err := url.ParseRequestURI(s)
	return err == nil
}