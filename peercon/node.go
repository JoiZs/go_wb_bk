package peercon

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/config"
	"github.com/libp2p/go-libp2p/core/host"
	peerstore "github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/relay"
	quicTransport "github.com/libp2p/go-libp2p/p2p/transport/quic"
	webrtc "github.com/libp2p/go-libp2p/p2p/transport/webrtc"
	webtransport "github.com/libp2p/go-libp2p/p2p/transport/webtransport"
)

type LibNode struct {
	listenAddr []string
	node       *host.Host
}

func InitNode(addr ...string) *LibNode {
	return &LibNode{
		listenAddr: addr,
	}
}

func (r *LibNode) Start() {
	var confOpt []config.Option

	confOpt = append(confOpt,
		libp2p.Transport(quicTransport.NewTransport),
		libp2p.Transport(webtransport.New),
		libp2p.Transport(webrtc.New),
		libp2p.ListenAddrStrings(r.listenAddr...),
	)

	node, err := libp2p.New(confOpt...)
	if err != nil {
		log.Fatal(err)
	}

	r.node = &node
	r.SubTopic()
	r.MakeRelay()

	relayInfo := peerstore.AddrInfo{
		ID:    node.ID(),
		Addrs: node.Addrs(),
	}

	_, err = peerstore.AddrInfoToP2pAddrs(&relayInfo)
	if err != nil {
		log.Fatal(err)
	}
	for _, addr := range node.Addrs() {
		log.Printf("Listening on: %s/p2p/%s\n", addr.String(), node.ID())
	}
	// fmt.Println("Addrs: ", relayAddrs)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	<-ch
	if err := node.Close(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Shut down server....")
}

func (r *LibNode) MakeRelay() {
	_, err := relay.New(*r.node)
	if err != nil {
		log.Fatal(err)
	}
}

func (r *LibNode) SubTopic() {
	ctx := context.Background()
	topicVal := os.Getenv("topic")

	ps, err := pubsub.NewGossipSub(ctx, *r.node)
	if err != nil {
		log.Fatal(err)
	}

	tp, err := ps.Join(topicVal)
	if err != nil {
		log.Fatal(err)
	}

	sub, err := tp.Subscribe()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(sub.Topic())
}
