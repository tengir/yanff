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

	fmt.Fprintln(step2, "set 2 dst mac", config[*target][0])
	fmt.Fprintln(step2, "set 3 dst mac", config[*target][1])
}
