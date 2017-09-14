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
	target := flag.String("target", "dcomp01", "Target host name from config file")
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

	// Write script file for step4
	step4, err := os.Create(*pktgenDir + string(os.PathSeparator) + "step4.pg")
    if err != nil {
		log.Fatal(err)
	} else {
		defer step4.Close()
	}

	fmt.Fprintln(step4, "range 0 dst mac start", config[*target][0])
	fmt.Fprintln(step4, "range 1 dst mac start", config[*target][1])
	fmt.Fprintln(step4, `range 0 dst port start 50
range 0 dst port min 50
range 0 dst port max 60
range 0 dst port inc 1
enable 0 range
`)

	// Write script file for step5
	step5, err := os.Create(*pktgenDir + string(os.PathSeparator) + "step5.pg")
    if err != nil {
		log.Fatal(err)
	} else {
		defer step5.Close()
	}

	fmt.Fprintln(step5, "range 0 dst mac start", config[*target][0])
	fmt.Fprintln(step5, "range 1 dst mac start", config[*target][1])
	fmt.Fprintln(step5, `range 0 src ip start 111.2.0.0
range 0 src ip min 111.2.0.0
range 0 src ip max 111.2.0.3
range 0 src ip inc 0.0.0.1
enable 0 range
`)
}
