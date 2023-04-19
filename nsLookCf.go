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
)

func main() {
    r := &net.Resolver{
        PreferGo: true,
        Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
            d := net.Dialer{
                Timeout: time.Millisecond * time.Duration(10000),
            }
            return d.DialContext(ctx, network, "1.1.1.1:53")
        },
    }
    nsRecords, err := r.LookupNS(context.Background(), "azulsoftware.eu")
	if err != nil {log.Fatal("failed r.LookUpNs: %v\n", err)}

	log.Printf("number of ns records: %d\n", len(nsRecords))

	for i:=0; i< len(nsRecords); i++ {
		fmt.Printf("rec [%d]: %s\n", i +1, nsRecords[i].Host)
	}

}
