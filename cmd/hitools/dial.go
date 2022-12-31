package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/DerZombiiie/HiLib"
	"os"
	"strings"
)

var (
	dialCommands = make(map[string]func())
)

func init() {
	commands["dial"] = func() {
		var cmd string

		if len(os.Args) >= 2 {
			cmd = strings.ToLower(os.Args[1])

			os.Args = append(os.Args[0:1], os.Args[2:]...)
		} else {
			os.Args = []string{os.Args[0]}
		}

		flag.Parse()

		if len(os.Args) < 1 {
			helpSim()
		}

		c, ok := dialCommands[cmd]
		if !ok {
			helpDial()
		}

		c()
	}

	dialCommands["connect"] = func() {
		res := dial(hilib.DialConnect)

		resp, ok := res.(*hilib.ResDial)
		if !ok {
			fmt.Fprintf(os.Stderr, "Error: unexpected type %T\n", res)
		} else {
			specialPrint(resp)

			fmt.Println(resp.Response)
		}
	}

	commands["connect"] = dialCommands["connect"]
	dialCommands["conn"] = dialCommands["connect"]

	dialCommands["disconnect"] = func() {
		res := dial(hilib.DialDisconnect)

		resp, ok := res.(*hilib.ResDial)
		if !ok {
			fmt.Fprintf(os.Stderr, "Error: unexpected type %T\n", res)
		} else {
			specialPrint(resp)

			fmt.Println(resp.ResString.Response)
		}
	}

	commands["disconnect"] = dialCommands["disconnect"]
	dialCommands["disco"] = dialCommands["disconnect"]
}

func dial(a hilib.DialAction) hilib.Response {
	return request(&hilib.ReqDial{
		Token:  token(),
		Action: a,
	})
}

func helpDial() {
	fmt.Fprintf(os.Stderr, "Usage: HiTools sim [COMMAND] arguments\nwhere COMMAND := {%s}\n", listDialCommands())

	os.Exit(-1)
}

func listDialCommands() string {
	buf := &bytes.Buffer{}

	fmt.Fprint(buf, " ")

	for k := range dialCommands {
		fmt.Fprintf(buf, "%s ", k)
	}

	return buf.String()
}
