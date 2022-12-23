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
	statusCommands = make(map[string]func())
)

func init() {
	commands["status"] = func() {
		var cmd string

		if len(os.Args) >= 2 {
			cmd = strings.ToLower(os.Args[1])

			os.Args = append(os.Args[0:1], os.Args[2:]...)
		} else {
			os.Args = []string{os.Args[0]}
		}

		flag.Parse()

		if len(os.Args) < 1 {
			helpStatus()
		}

		c, ok := statusCommands[cmd]
		if !ok {
			helpStatus()
		}

		c()
	}

	statusCommands["connection"] = func() {
		res := request(&hilib.ReqStatus{})

		status, ok := res.(*hilib.ResStatus)
		if !ok {
			fmt.Fprintf(os.Stderr, "Error: unexpected type %T\n", res)
		} else {
			specialPrint(status)

			fmt.Printf("Connection status: %s\n", status.ConnStatus)
			fmt.Printf("Signal strength: %s\n", status.SignalStrengthPercent())
		}
	}
	statusCommands["conn"] = statusCommands["connection"]

	statusCommands["all"] = func() {
		res := request(&hilib.ReqStatus{})

		status, ok := res.(*hilib.ResStatus)
		if !ok {
			fmt.Fprintf(os.Stderr, "Error: unexpected type %T\n", res)
		} else {
			if *mode == "normal" {
				str := "verbose"
				mode = &str
			}

			specialPrint(status)
		}
	}

	commands["ip"] = func() {
		res := request(&hilib.ReqStatus{})
		status, ok := res.(*hilib.ResStatus)
		if !ok {
			fmt.Fprintf(os.Stderr, "Error: unexpected type %T\n", res)
		} else {
			specialPrint(status)

			fmt.Printf("IPv4: %s\nIPv6: %s\n",
				t(status.WanIPAddress, "<not available>"),
				t(status.WanIPv6Address, "<not available>"),
			)
		}
	}

	statusCommands["ip"] = commands["ip"]
}

func t(a, b string) string {
	if a == "" {
		return b
	} else {
		return a
	}
}

func helpStatus() {
	fmt.Fprintf(os.Stderr, "Usage: HiTools status [COMMAND] arguments\nwhere COMMAND := {%s}", listStatusCommands())

	os.Exit(-1)
}

func listStatusCommands() string {
	buf := &bytes.Buffer{}

	fmt.Fprint(buf, " ")

	for k := range statusCommands {
		fmt.Fprintf(buf, "%s ", k)
	}

	return buf.String()
}
