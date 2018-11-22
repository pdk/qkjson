package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/pdk/qkjson/parser"
)

var (
	outFile string
)

func init() {
	flag.StringVar(&outFile, "output", "", "output file name")
	flag.StringVar(&outFile, "o", "", "output file name")
}

func main() {
	flag.Parse()

	data := parser.ParseArgs(flag.Args())
	if data == nil {
		log.Fatalf("cannot produce JSON from nada")
	}

	writeResult(outFile, data)
}

func writeResult(outFile string, data interface{}) {

	result, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		log.Fatalf("cannot make JSON: %s", err)
	}

	if outFile == "" {
		_, err := os.Stdout.Write(result)
		if err != nil {
			log.Fatalf("cannot write to stdout: %s", err)
		}

		return
	}

	err = ioutil.WriteFile(outFile, result, os.ModePerm)
	if err != nil {
		log.Fatalf("cannot write result %s: %s", outFile, err)
	}
}
