package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/DerZombiiie/HiLib"
	"os"
	"reflect"
	"strings"
)

var (
	basePath = flag.String("base-path", "http://192.168.8.1/", "Base url for http-api")
	mode     = flag.String("mode", "normal", "Set mode, can be script, normal, verbose and raw")
	field    = flag.String("field", "", "Only works with in script mode")
)

var (
	commands = make(map[string]func())

	base string
)

var config *hilib.Config

func listCommands() string {
	buf := &bytes.Buffer{}

	fmt.Fprint(buf, " ")

	for k := range commands {
		fmt.Fprintf(buf, "%s ", k)
	}

	return buf.String()
}

func main() {
	if len(os.Args) < 2 {
		help()
	}

	cmd := strings.ToLower(os.Args[1])
	if len(os.Args) >= 2 {
		if len(os.Args) < 3 {
			os.Args = append(os.Args[0:1], os.Args[2:]...)
		} else {
			os.Args = append(os.Args[0:1], os.Args[3:]...)
		}
	} else {
		os.Args = []string{os.Args[0]}
	}

	flag.Parse()

	// check mode
	_, ok := map[string]struct{}{"script": {}, "normal": {}, "verbose": {}, "raw": {}}[*mode]
	if !ok {
		fmt.Fprintln(os.Stderr, "mode argument can only be script, normal, verbose and raw")
		os.Exit(-1)
	}

	config = hilib.NewConfig(*basePath)

	c, ok := commands[cmd]
	if !ok {
		help()
	} else {
		c()
	}
}

func request(r hilib.Request) hilib.Response {
	if *mode == "verbose" {
		fmt.Fprintf(os.Stderr, "Requesting: %s%s\n", config.BaseURL, r.ReqPath())
	}
	res, err := r.Request(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error during Request: %s\n", err)
		os.Exit(-1)
	}

	return res
}
func help() {
	fmt.Fprintf(os.Stderr,
		`HiTools beta

Usage: hitools [COMMAND] arguments
where  COMMAND := {%s}
`, listCommands())
	os.Exit(-1)
}

func specialPrint(data hilib.Response) {
	switch *mode {
	case "verbose":
		str := fmt.Sprintf("%+v", data)

		fmt.Println(
			"\nResponse:\n  " + strings.ReplaceAll(str[2:len(str)-1], " ", "\n  "),
		)

	case "script":
		if *field == "" {
			fmt.Printf("Error: no field specified!")
			os.Exit(-1)
		}

		f := reflect.ValueOf(data).Elem().FieldByName(*field)
		if !f.IsValid() {
			fmt.Printf("Error: field '%s' can't be found!\n", *field)
			os.Exit(-1)
		}

		fmt.Printf("%v", f)

	case "raw":
		fmt.Printf("%s", data.Raw())

	default:
		return
	}

	os.Exit(0)
}
