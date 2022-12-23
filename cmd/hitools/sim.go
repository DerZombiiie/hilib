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
	simCommands = make(map[string]func())
)

func init() {
	commands["sim"] = func() {
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

		c, ok := simCommands[cmd]
		if !ok {
			helpSim()
		}

		c()

	}

	simCommands["status"] = func() {
		res := request(&hilib.ReqSimStatus{})

		status, ok := res.(*hilib.ResSimStatus)
		if !ok {
			fmt.Fprintf(os.Stderr, "Error: unexpected type %T\n", res)
		} else {
			specialPrint(status)

			fmt.Printf("Sim status: %d\n", status.SimState) //.String())
		}
	}

	simCommands["lock"] = func() {
		res := request(&hilib.ReqSimLock{})

		status, ok := res.(*hilib.ResSimLock)
		if !ok {
			fmt.Fprintf(os.Stderr, "Error: unexpected type %T\n", res)
		} else {
			specialPrint(status)

			fmt.Printf("Simlock:      %s\n", status.SimLockEnable.OnOff()) //.String())
			fmt.Printf("Unlock tries: %d\n", status.SimLockRemainTimes)    //.String())
		}
	}

	statusCommands["sim"] = simCommands["status"]
}

func helpSim() {
	fmt.Fprintf(os.Stderr, "Usage: HiTools sim [COMMAND] arguments\nwhere COMMAND := {%s}\n", listSimCommands())

	os.Exit(-1)
}

func listSimCommands() string {
	buf := &bytes.Buffer{}

	fmt.Fprint(buf, " ")

	for k := range simCommands {
		fmt.Fprintf(buf, "%s ", k)
	}

	return buf.String()
}
