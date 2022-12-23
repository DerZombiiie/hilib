package main

import (
	"fmt"
	"github.com/DerZombiiie/HiLib"
	"os"
)

func token() int {
	if *mode == "verbose" {
		fmt.Fprintln(os.Stderr, "Requesting token")
	}

	res := request(&hilib.ReqToken{})
	return res.(*hilib.ResToken).Token
}

func init() {
	commands["token"] = func() {
		res := request(&hilib.ReqToken{})

		tkn, ok := res.(*hilib.ResToken)
		if !ok {
			fmt.Fprintf(os.Stderr, "Error: unexpected type %T\n", res)
		} else {
			specialPrint(tkn)

			fmt.Printf("Token: %d\n", tkn.Token)
		}
	}
}
