package main

import (
	"fmt"
	"github.com/DerZombiiie/HiLib"
	"os"
)

func init() {
	commands["signal"] = func() {
		res := request(&hilib.ReqSignalStatus{})

		signal, ok := res.(*hilib.ResSignalStatus)
		if !ok {
			fmt.Fprintf(os.Stderr, "Error: unexpected type %T\n", res)
		} else {
			specialPrint(signal)

			fmt.Printf("RSSI: %s\n", signal.RSSI)
		}
	}

	commands["location"] = func() {
		res := request(&hilib.ReqSignalStatus{})

		signal, ok := res.(*hilib.ResSignalStatus)
		if !ok {
			fmt.Fprintf(os.Stderr, "Error: unexpected type %T\n", res)
		} else {
			specialPrint(signal)

			fmt.Printf("Cell_id: %d\n", signal.Cell_id)
		}
	}
}
