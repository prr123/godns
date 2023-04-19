// acmeDnsRec.go
// Author: prr azul software
// Date: 19 April 2023
// copyright prr, azul software
//
package main

import (
    "context"
    "net"
    "time"
	"log"
	"fmt"
	"os"
)

func main() {
	var txtRecs []string

	numarg := len(os.Args)

//	ctx := context.Background()

	useStr := "usage: acmeDnsRec domain [provider]"

	if numarg < 2 {
		fmt.Println(useStr)
		fmt.Printf("insufficient arguments!\n")
		os.Exit(-1)
	}

	if numarg > 3 {
		fmt.Println(useStr)
		fmt.Printf("too many arguments: %s\n", useStr)
		os.Exit(-1)
	}

	provider := ""
	domain := os.Args[1]

	if numarg == 3 {
		provider = os.Args[2]
	}

	domain = "_acme-challenge." + domain

	log.Printf("TXT Records for domain: %s\n", domain)
	if provider == "cf" {
		r := &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := net.Dialer{
					Timeout: time.Millisecond * time.Duration(10000),
				}
				return d.DialContext(ctx, network, "1.1.1.1:53")
			},
		}


		txtrecs, err :=r.LookupTXT(context.Background(), domain)
		if err != nil {log.Fatalf("failed r.LookUpTXT: %v\n", err)}
		txtRecs = txtrecs

	} else {
		txtrecs, err := net.LookupTXT(domain)
		if err != nil {log.Fatalf("failed net.LookUpTxT: %v\n", err)}
		txtRecs = txtrecs
	}

	log.Printf("domain:      %s\n", domain)
	log.Printf("NS provider: %s\n", provider)

	PrintTxtRecs(&txtRecs)
}

func PrintNsRecs(nsRecords []*net.NS) {
	fmt.Printf("***************** NS Records: %d ***********\n", len(nsRecords))
	hosts := make([]string, len(nsRecords))
	for i:=0; i< len(nsRecords); i++ {
		tmp := []byte(nsRecords[i].Host)
		hosts[i] = string(tmp[:len(tmp)-1])
		fmt.Printf("NS Rec [%d]: %s\n", i +1, hosts[i])
	}
	fmt.Printf("********************************************\n")
}

func PrintMxRecs(mxRecs []*net.MX) {
	fmt.Printf("***************** MX records: %d ***********\n", len(mxRecs))
	hosts := make([]string, len(mxRecs))
	for i:=0; i< len(mxRecs); i++ {
		tmp := []byte(mxRecs[i].Host)
		hosts[i] = string(tmp[:len(tmp)-1])
		fmt.Printf("Mx Rec [%d]: %s Pref %d\n", i +1, hosts[i], mxRecs[i].Pref)
	}
	fmt.Printf("********************************************\n")
}

func PrintAddrs(addrs *[]string) {
	fmt.Printf("***************** domain addr: %d ***********\n", len(*addrs))
	for i:=0; i< len(*addrs); i++ {
		fmt.Printf("adr[%d]: %s\n", i+1, (*addrs)[i])
	}
	fmt.Printf("********************************************\n")
}

func PrintTxtRecs(txtrecs *[]string) {
	fmt.Printf("***************** TXT Records: %d ***********\n", len(*txtrecs))
	for i:=0; i< len(*txtrecs); i++ {
		fmt.Printf("rec[%d]: %s\n", i+1, (*txtrecs)[i])
	}
	fmt.Printf("********************************************\n")
}
