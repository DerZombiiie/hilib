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
	smsCommands = make(map[string]func())
)

func init() {
	commands["sms"] = func() {
		var cmd string

		if len(os.Args) >= 2 {
			cmd = strings.ToLower(os.Args[1])

			os.Args = append(os.Args[0:1], os.Args[2:]...)
		} else {
			os.Args = []string{os.Args[0]}
		}

		flag.Parse()

		if len(os.Args) < 1 {
			helpSMS()
		}

		c, ok := smsCommands[cmd]
		if !ok {
			helpSMS()
		}

		c()
	}

	smsCommands["list"] = func() {
		res := request(&hilib.ReqSmsList{
			Token: token(),

			PageIndex:       1,
			ReadCount:       20,
			BoxType:         1,
			SortType:        0,
			Ascending:       0,
			UnreadPreferred: 0,
		})

		list, ok := res.(*hilib.ResSmsList)
		if !ok {
			fmt.Fprintf(os.Stderr, "Error: unexpected type %T\n", res)
		} else {
			specialPrint(list)

			println("TESTSETET")
		}
	}
}

func helpSMS() {
	fmt.Fprintf(os.Stderr, "Usage: HiTools sms [COMMAND] arguments\nwhere COMMAND := {%s}\n", listSMSCommands())

	os.Exit(-1)
}

func listSMSCommands() string {
	buf := &bytes.Buffer{}

	fmt.Fprint(buf, " ")

	for k := range smsCommands {
		fmt.Fprintf(buf, "%s ", k)
	}

	return buf.String()
}
