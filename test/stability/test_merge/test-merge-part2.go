// Copyright 2017 Intel Corporation.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"

	"github.com/intel-go/yanff/flow"
	"github.com/intel-go/yanff/packet"
	"github.com/intel-go/yanff/test/stability/stabilityCommon"
)

var (
	inport1 uint
	inport2 uint
	outport uint

	fixMACAddrs func(*packet.Packet, flow.UserContext)
)

// Main function for constructing packet processing graph.
func main() {
	flag.UintVar(&inport1, "inport1", 0, "port for 1st receiver")
	flag.UintVar(&inport2, "inport2", 1, "port for 2nd receiver")
	flag.UintVar(&outport, "outport", 0, "port for sender")
	configFile := flag.String("config", "config.json", "Specify config file name")
	target := flag.String("target", "nntsat01g4", "Target host name from config file")
	flag.Parse()

	// Init YANFF system at requested number of cores.
	config := flow.Config{
		CPUCoresNumber: 16,
	}
	flow.SystemInit(&config)
	stabilityCommon.InitCommonState(*configFile, *target)
	fixMACAddrs = stabilityCommon.ModifyPacket[outport].(func(*packet.Packet, flow.UserContext))

	// Receive packets from 0 and 1 ports
	inputFlow1 := flow.SetReceiver(uint8(inport1))
	inputFlow2 := flow.SetReceiver(uint8(inport2))

	outputFlow := flow.SetMerger(inputFlow1, inputFlow2)
	flow.SetHandler(outputFlow, fixPackets, nil)
	flow.SetSender(outputFlow, uint8(outport))

	// Begin to process packets.
	flow.SystemStart()
}

func fixPackets(pkt *packet.Packet, ctx flow.UserContext) {
	if stabilityCommon.ShouldBeSkipped(pkt) {
		return
	}
	fixMACAddrs(pkt, ctx)
}
