package main

import (
	"github.com/summeroch/dht-spider/pkg/dht"
	"log"
	"os"
	"time"
)

func main() {
	f := "./getpeers.log"
	logFile, err := os.OpenFile(f, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("%s", err)
	}
	log.SetOutput(logFile)

	d := dht.New(nil)
	d.OnGetPeersResponse = func(infoHash string, peer *dht.Peer) {
		log.Printf("got peer: %s:%d\n", peer.IP, peer.Port)
	}

	go func() {
		for {
			// ubuntu-20.04.3-desktop-amd64.iso
			err := d.GetPeers("b26c81363ac1a236765385a702aec107a49581b5")
			if err != nil && err != dht.ErrNotReady {
				log.Fatal(err)
			}
			if err == dht.ErrNotReady {
				time.Sleep(time.Second * 1)
				continue
			}
			break
		}
	}()

	d.Run()
}
