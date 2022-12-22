package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

func printUsage() {
	fmt.Fprintf(os.Stderr, `Usage: jwt-decode [JWT]
Reads from STDIN if JWT is not provided as an argument.

Flags:
`)
	flag.PrintDefaults()
	os.Exit(1)
}

func printJson(str string) {
	var buf bytes.Buffer
	err := json.Indent(&buf, []byte(str), "", "  ")
	if err != nil {
		fmt.Println(str)
		return
	}
	fmt.Println(buf.String())
}

func main() {
	help := flag.Bool("h", false, "print help")
	flag.Parse()
	if *help {
		printUsage()
		return
	}
	if flag.NArg() > 1 {
		printUsage()
		return
	}
	input := ""
	if flag.NArg() == 1 {
		input = flag.Arg(0)
	} else {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			panic(err)
		}
		input = strings.TrimSpace(string(data))
	}
	parts := strings.Split(input, ".")
	if len(parts) != 3 {
		fmt.Fprintln(os.Stderr, "invalid format:", input)
		os.Exit(1)
		return
	}
	encodedHeader := parts[0]
	encodedPayload := parts[1]
	signature := parts[2]
	headerAsBytes, err := base64.RawURLEncoding.DecodeString(encodedHeader)
	if err != nil {
		fmt.Fprintln(os.Stderr, "header:", encodedHeader, err)
		os.Exit(1)
		return
	}
	payloadAsBytes, err := base64.RawURLEncoding.DecodeString(encodedPayload)
	if err != nil {
		fmt.Fprintln(os.Stderr, "payload:", encodedPayload, err)
		os.Exit(1)
		return
	}
	printJson(string(headerAsBytes))
	fmt.Println("#")
	printJson(string(payloadAsBytes))
	fmt.Println("#")
	fmt.Println(signature)
}
