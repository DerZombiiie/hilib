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
		res := request(&hilib.ReqSMSList{
			Token: token(),

			PageIndex:       1,
			ReadCount:       20,
			BoxType:         hilib.BoxInbox,
			SortType:        0,
			Ascending:       0,
			UnreadPreferred: 0,
		})

		list, ok := res.(*hilib.ResSMSList)
		if !ok {
			fmt.Fprintf(os.Stderr, "Error: unexpected type %T\n", res)
		} else {
			specialPrint(list)

			fmt.Printf("Got %d messages in 1 request(s)\n", len(list.Messages))

			for k := range list.Messages {
				fmt.Printf("%s\n---\n", formatMessage(list.Messages[k]))

			}
		}
	}

	smsCommands["listen"] = func() {
		msgch, errch := hilib.SMSChan()

		for {
			select {
			case msg := <-msgch:
				fmt.Printf("%s\n---\n", formatMessage(msg))

			case err := <-errch:
				fmt.Printf("ERROR: %s\n", err)
			}
		}
	}

	smsCommands["send"] = func() {
		if len(os.Args) < 3 {
			fmt.Printf("%v\n", os.Args)
			fmt.Fprintf(os.Stderr, "Usage: HiTools sms send \"<message>\" <phone1> <phone2> <phone3>\n")
			os.Exit(-1)
		}

		res := request(&hilib.ReqSendSMS{
			Token: token(),

			Phones:  os.Args[2:],
			Content: os.Args[1],
		})

		resp, ok := res.(*hilib.ResSendSMS)
		if !ok {
			fmt.Fprintf(os.Stderr, "Error: unexpected type %T\n", res)
		} else {
			specialPrint(resp)

			fmt.Println(resp.Response)
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

func formatMessage(msg hilib.SMSMessage) string {
	buf := &bytes.Buffer{}

	fmt.Fprintf(buf, "From: %s; date %s (%s)\n", msg.Phone, msg.Date, msg.Type.String())
	fmt.Fprintf(buf, "Content:\n")
	fmt.Fprintf(buf, "  "+strings.ReplaceAll(msg.Content, "\n", "\n  "))

	return buf.String()
}
