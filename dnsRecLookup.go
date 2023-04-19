// dnsRecLookup
// from: https://stackoverflow.com/questions/59889882/specifying-dns-server-for-lookup-in-go
// modified
// Author: prr azul software
// Date 10 April 2023
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
	var domain, provider string
	var nsRecords []*net.NS
	var domAddrs, txtRecs, txtAcmeRecs []string
	var mxRecs []*net.MX
	numarg := len(os.Args)

//	ctx := context.Background()

	useStr := "usage: dnsRecLookup domain [provider]"

	switch numarg {
	case 1:
		fmt.Printf("insufficient arguments: %s\n", useStr)
		os.Exit(-1)
	case 2:
		domain = os.Args[1]

	case 3:
		domain = os.Args[1]
		provider = os.Args[2]

	default:
		fmt.Printf("too many arguments: %s\n", useStr)
		os.Exit(-1)
	}

	if provider == "" {provider = "cf"}

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

		nsRecs, err := r.LookupNS(context.Background(), domain)
		if err != nil {log.Fatal("failed r.LookUpNs: %v\n", err)}
		nsRecords = nsRecs

		addrs, err := r.LookupHost(context.Background(), domain)
		if err != nil {log.Fatal("failed r.LookUpHost: %v\n", err)}
		domAddrs = addrs

//		addrs, err := r.LookupCNAME(context.Background(), domain)
//		if err != nil {log.Fatal("failed r.LookUpHost: %v\n", err)}
//		domAddrs = addrs

		txtrecs, err :=r.LookupTXT(context.Background(), domain)
		if err != nil {log.Fatal("failed r.LookUpTXT: %v\n", err)}
		txtRecs = txtrecs

//		txtacmerecs, err :=r.LookupTXT(context.Background(), "_acme-challenge")
//		if err != nil {log.Fatal("failed r.LookUpTXT: %v\n", err)}
//		txtAcmeRecs = txtacmerecs

		mxrecs, err := r.LookupMX(context.Background(), domain)
		if err != nil {log.Fatal("failed r.LookUpMX: %v\n", err)}
		mxRecs = mxrecs

	} else {
		nsRecs, err := net.LookupNS(domain)
		if err != nil {log.Fatal("failed net.LookUpNs: %v\n", err)}
		nsRecords = nsRecs
	}

	log.Printf("domain:      %s\n", domain)
	log.Printf("NS provider: %s\n", provider)

	PrintAddrs(&domAddrs)
	PrintTxtRecs(&txtRecs)
	PrintTxtRecs(&txtAcmeRecs)
	PrintMxRecs(mxRecs)
	PrintNsRecs(nsRecords)
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
