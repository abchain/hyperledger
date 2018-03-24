package main

import (
	"encoding/json"
	"flag"
	"fmt"
	token "hyperledger.abchain.org/chaincode/generaltoken"
	reg "hyperledger.abchain.org/chaincode/registrar"
	"math/big"
	"os"
)

func main() {

	short := flag.Bool("short", false, "Short mode: only print the arguments in JSON format")
	showhelp := flag.Bool("help", false, "Show help info")

	var manager string
	var regmanager string

	flag.StringVar(&manager, "admin", "Admin", "Set string for manager role")
	flag.StringVar(&regmanager, "regmanager", "", "Set string for registrar manager role")

	flag.Parse()

	if *showhelp {

		flag.PrintDefaults()

	} else if flag.NArg() == 0 {

		fmt.Println("Must provide the total token number")
		os.Exit(1)

	} else {

		total, ok := big.NewInt(0).SetString(flag.Arg(0), 0)
		if !ok {

			fmt.Println("Invalid total count")
			os.Exit(1)
		}

		var args []string

		args, err := token.CCDeploy(total, args)
		if err != nil {
			fmt.Println("gen token parameter fail", err)
			os.Exit(1)
		}

		if regmanager == "" {
			regmanager = manager
		}

		args, err = reg.CCDeploy(manager, regmanager, args)
		if err != nil {
			fmt.Println("gen registrar parameter fail", err)
			os.Exit(1)
		}

		retbyte, err := json.Marshal(args)
		if err != nil {
			fmt.Println("encoding arguments fail", err)
			os.Exit(1)
		}

		if *short {
			fmt.Print(string(retbyte))
		} else {
			fmt.Println("Deployment helper for AECC")
			fmt.Println("")
			fmt.Println("Deploy for", total, "tokens")
			fmt.Println("Deploy specified manager and regmanager as", manager, regmanager)
			fmt.Println("Encoding parameter as:")
			fmt.Println("  ", string(retbyte))
			fmt.Println("--------------------------------------------")
		}
	}

}
