package main

import (
	"github.com/JoiZs/wb2/initializer"
	"github.com/JoiZs/wb2/peercon"
)

func main() {
	initializer.Init()

	relayN := peercon.InitNode("/ip4/127.0.0.1/tcp/6565")

	relayN.Start()
	relayN.MakeRelay()
}
