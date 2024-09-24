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
	"github.com/libp2p/go-libp2p/core/host"
	peerstore "github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/relay"
)

type LibNode struct {
	listenAddr string
	node       *host.Host
}

func InitNode(addr string) *LibNode {
	return &LibNode{
		listenAddr: addr,
	}
}

func (r *LibNode) Start() {
	node, err := libp2p.New(libp2p.ListenAddrStrings(r.listenAddr))
	if err != nil {
		log.Fatal(err)
	}

	r.node = &node
	r.SubTopic()

	relayInfo := peerstore.AddrInfo{
		ID:    node.ID(),
		Addrs: node.Addrs(),
	}

	relayAddrs, err := peerstore.AddrInfoToP2pAddrs(&relayInfo)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Addrs: ", relayAddrs)

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
