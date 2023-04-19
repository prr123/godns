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

	numArgs := len(os.Args)

	useStr := "usage is: \"nsRecCf domain\""
	if numArgs < 2 {
		fmt.Println(useStr)
		fmt.Printf("insufficent args!\n")
		os.Exit(-1)
	}
	if numArgs >2 {
		fmt.Println(useStr)
		fmt.Printf("too many args!\n")
		os.Exit(-1)
	}


	domain := os.Args[1]

    r := &net.Resolver{
        PreferGo: true,
        Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
            d := net.Dialer{
                Timeout: time.Millisecond * time.Duration(10000),
            }
            return d.DialContext(ctx, network, "1.1.1.1:53")
        },
    }

	log.Printf("Using Cloudflare 1.1.1.1 to resolve ns recs for %s\n", domain)

    nsRecords, err := r.LookupNS(context.Background(), domain)
	if err != nil {log.Fatal("failed r.LookUpNs: %v\n", err)}

	log.Printf("number of ns records: %d\n", len(nsRecords))

	for i:=0; i< len(nsRecords); i++ {
		log.Printf("rec [%d]: %s\n", i +1, nsRecords[i].Host)
	}

	log.Printf("success\n")
}
