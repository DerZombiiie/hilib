package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/DerZombiiie/HiLib"
	"os"
	"strings"
	"time"
)

var (
	statCommands = make(map[string]func())
)

func init() {
	commands["statistics"] = func() {
		var cmd string

		if len(os.Args) >= 2 {
			cmd = strings.ToLower(os.Args[1])

			os.Args = append(os.Args[0:1], os.Args[2:]...)
		} else {
			os.Args = []string{os.Args[0]}
		}

		flag.Parse()

		if len(os.Args) < 1 {
			helpStat()
		}

		c, ok := statCommands[cmd]
		if !ok {
			helpStat()
		}

		c()
	}
	commands["stat"] = commands["statistics"]

	statCommands["all"] = func() {
		res := request(&hilib.ReqTrafficStats{})

		status, ok := res.(*hilib.ResTrafficStats)
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

	statCommands["data"] = func() {
		res := request(&hilib.ReqTrafficStats{})

		status, ok := res.(*hilib.ResTrafficStats)
		if !ok {
			fmt.Fprintf(os.Stderr, "Error: unexpected type %T\n", res)
		} else {
			specialPrint(status)

			fmt.Printf("Total Upload:   %s\n", hilib.ByteCountIEC(status.TotalUpload))
			fmt.Printf("Total Download: %s\n", hilib.ByteCountIEC(status.TotalDownload))
		}
	}

	statCommands["speed"] = func() {
		res := request(&hilib.ReqTrafficStats{})

		status, ok := res.(*hilib.ResTrafficStats)
		if !ok {
			fmt.Fprintf(os.Stderr, "Error: unexpected type %T\n", res)
		} else {
			specialPrint(status)

			fmt.Printf("Upload speed:   %s/s\n", hilib.ByteCountIEC(status.CurrentUploadRate))
			fmt.Printf("Download speed: %s/s\n", hilib.ByteCountIEC(status.CurrentDownloadRate))
		}
	}

	statCommands["time"] = func() {
		res := request(&hilib.ReqTrafficStats{})

		status, ok := res.(*hilib.ResTrafficStats)
		if !ok {
			fmt.Fprintf(os.Stderr, "Error: unexpected type %T\n", res)
		} else {
			specialPrint(status)

			fmt.Printf("CurrentConnectTime: %s\n", (time.Second * time.Duration(status.CurrentConnectTime)).String())
			fmt.Printf("TotalConnectTime:   %s\n", (time.Second * time.Duration(status.TotalConnectTime)).String())
		}
	}
}

func helpStat() {
	fmt.Fprintf(os.Stderr, "Usage: HiTools statistics [COMMAND] arguments\nwhere COMMAND := {%s}\n", listStatCommands())

	os.Exit(-1)
}

func listStatCommands() string {
	buf := &bytes.Buffer{}

	fmt.Fprint(buf, " ")

	for k := range statCommands {
		fmt.Fprintf(buf, "%s ", k)
	}

	return buf.String()
}
