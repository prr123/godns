// nsLookup
// from: https://stackoverflow.com/questions/59889882/specifying-dns-server-for-lookup-in-go
// modified
// Author: prr axul software
// Date 2 April 2023
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
	numarg := len(os.Args)

	useStr := "usage: nsLookup domain [provider]"

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
	} else {
		nsRecs, err := net.LookupNS(domain)
		if err != nil {log.Fatal("failed net.LookUpNs: %v\n", err)}
		nsRecords = nsRecs
	}

	log.Printf("domain:      %s\n", domain)
	log.Printf("NS provider: %s\n", provider)

	log.Printf("number of ns records: %d\n", len(nsRecords))
	hosts := make([]string, len(nsRecords))
	for i:=0; i< len(nsRecords); i++ {
		tmp := []byte(nsRecords[i].Host)
		hosts[i] = string(tmp[:len(tmp)-1])
		log.Printf("rec [%d]: %s\n", i +1, hosts[i])
	}

}
