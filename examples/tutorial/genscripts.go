package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	// Parse arguments
	pktgenDir := flag.String("pgdir", "../../dpdk", "Specify directory where to put pktgen script files")
	configFile := flag.String("config", "config.json", "Specify config file name")
	target := flag.String("target", "", "Target host name from config file")
	pktgen := flag.String("pktgen", "", "Pktgen host name from config file")
	flag.Parse()

	// Read config
	err := readConfig(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	// Write script file for step2
	step2, err := os.Create(*pktgenDir + string(os.PathSeparator) + "step2.pg")
    if err != nil {
		log.Fatal(err)
	} else {
		defer step2.Close()
	}

	fmt.Fprintln(step2, "set 0 dst mac", config[*target][0])
	fmt.Fprintln(step2, "set 1 dst mac", config[*target][1])
	fmt.Fprintln(step2, `set 0 size 500
set 1 size 500
`)

	// Write script file for step4
	step4, err := os.Create(*pktgenDir + string(os.PathSeparator) + "step4.pg")
    if err != nil {
		log.Fatal(err)
	} else {
		defer step4.Close()
	}

	fmt.Fprintln(step4, "range 0 dst mac start", config[*target][0])
	fmt.Fprintln(step4, "range 1 dst mac start", config[*target][1])
	fmt.Fprintln(step4, `range 0 size start 500
range 0 size min 500
range 0 size max 500
range 0 dst port start 50
range 0 dst port min 50
range 0 dst port max 60
range 0 dst port inc 1
enable 0 range
`)

	// Write pktgen script file for step5
	step5, err := os.Create(*pktgenDir + string(os.PathSeparator) + "step5.pg")
    if err != nil {
		log.Fatal(err)
	} else {
		defer step5.Close()
	}

	fmt.Fprintln(step5, "range 0 dst mac start", config[*target][0])
	fmt.Fprintln(step5, "range 1 dst mac start", config[*target][1])
	fmt.Fprintln(step5, `range 0 size start 500
range 0 size min 500
range 0 size max 500
range 0 src ip start 111.2.0.0
range 0 src ip min 111.2.0.0
range 0 src ip max 111.2.0.3
range 0 src ip inc 0.0.0.1
enable 0 range
`)

	// Write pktgen script for NAT
	natpg, err := os.Create(*pktgenDir + string(os.PathSeparator) + "nat.pg")
    if err != nil {
		log.Fatal(err)
	} else {
		defer natpg.Close()
	}

	fmt.Fprintf(natpg, `range 0 dst mac start %s
range 1 dst mac start %s
range 0 size start 500
range 0 size min 500
range 0 size max 500
range 1 size start 500
range 1 size min 500
range 1 size max 500
range 0 src ip start 192.168.1.2
range 0 src ip min 192.168.1.2
range 0 src ip max 192.168.1.12
range 0 src ip inc 0.0.0.1
range 0 dst ip start 10.1.1.2
range 0 dst ip inc 0.0.0.0
range 0 src port start 1234
range 0 src port inc 0
range 1 src ip start 10.1.1.2
range 1 src ip inc 0.0.0.0
range 1 dst ip start 10.1.1.1
range 1 dst ip inc 0.0.0.0
range 1 dst port start 1024
range 1 dst port inc 0
enable 0 range
`, config[*target][0], config[*target][1])

	// Write NAT config file
	natjson, err := os.Create("nat.json")
    if err != nil {
		log.Fatal(err)
	} else {
		defer natjson.Close()
	}

	fmt.Fprintf(natjson, `{
    "private-port": {
        "index": 0,
        "dst_mac": "%s",
        "subnet": "192.168.1.1/24"
    },
    "public-port": {
        "index": 1,
        "dst_mac": "%s",
        "subnet": "10.1.1.1"
    }
}
`, config[*pktgen][0], config[*pktgen][1])
}
