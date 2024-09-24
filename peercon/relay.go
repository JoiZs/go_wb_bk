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

type RelayServer struct {
	listenAddr string
	node       *host.Host
}

func InitRelay(addr string) *RelayServer {
	return &RelayServer{
		listenAddr: addr,
	}
}

func (r *RelayServer) Start() {
	node, err := libp2p.New(libp2p.ListenAddrStrings(r.listenAddr))
	if err != nil {
		log.Fatal(err)
	}

	r.node = &node
	r.SubTopic()

	_, err = relay.New(node)
	if err != nil {
		log.Fatal(err)
	}

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

func (r *RelayServer) SubTopic() {
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
