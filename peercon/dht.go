package peercon

import (
	"context"
	"fmt"
	"os"
	"sync"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	dryrout "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	dutil "github.com/libp2p/go-libp2p/p2p/discovery/util"
)

func InitDHT(ctx context.Context, h host.Host) *dht.IpfsDHT {
	kdht, err := dht.New(ctx, h)
	if err != nil {
		panic(err)
	}

	if err = kdht.Bootstrap(ctx); err != nil {
		panic(err)
	}

	var wg sync.WaitGroup

	for _, peerAddr := range dht.DefaultBootstrapPeers {
		peerInfo, _ := peer.AddrInfoFromP2pAddr(peerAddr)
		wg.Add(1)

		go func() {
			defer wg.Done()

			if err := h.Connect(ctx, *peerInfo); err != nil {
				fmt.Print("Bootstrap err: ", err)
			}
		}()
	}

	wg.Wait()

	return kdht
}

func DiscoverPeer(ctx context.Context, h host.Host) {
	kdht := InitDHT(ctx, h)
	routingDiscovery := dryrout.NewRoutingDiscovery(kdht)

	topicVal := os.Getenv("topic")

	dutil.Advertise(ctx, routingDiscovery, topicVal)

	isConnect := false

	for !isConnect {
		fmt.Println("Looking for peers")
		peerCh, err := routingDiscovery.FindPeers(ctx, topicVal)
		if err != nil {
			panic(err)
		}

		for peer := range peerCh {
			if peer.ID == h.ID() {
				continue
			}

			err := h.Connect(ctx, peer)
			if err != nil {
				fmt.Printf("Failed connecting to %s, error: %s\n", peer.ID, err)
			} else {
				fmt.Println("Connected to peer: ", peer.ID)
				isConnect = true
			}

		}
	}
	fmt.Println("Completed peer discovery...")
}
