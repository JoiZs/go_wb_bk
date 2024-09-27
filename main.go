package main

import (
	"github.com/JoiZs/wb2/initializer"
	"github.com/JoiZs/wb2/peercon"
)

func main() {
	initializer.Init()

	relayN := peercon.InitNode("/ip4/0.0.0.0/udp/9096/webrtc-direct", "/ip4/0.0.0.0/udp/9095/quic-v1", "/ip4/0.0.0.0/udp/9095/quic-v1/webtransport", "/ip4/127.0.0.1/tcp/6565")

	relayN.Start()
}
