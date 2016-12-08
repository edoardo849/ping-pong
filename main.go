package main

import (
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p"

	log "github.com/Sirupsen/logrus"
)

// https://github.com/ethereum/go-ethereum/wiki/Peer-to-Peer
func main() {

	nodekey, _ := crypto.GenerateKey()
	srv := p2p.Server{
		Config: p2p.Config{
			MaxPeers:   10,
			PrivateKey: nodekey,
			Name:       "my node name",
			ListenAddr: ":30300",
			Protocols:  []p2p.Protocol{MyProtocol()},
		},
	}
	if err := srv.Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	info := srv.NodeInfo()

	log.WithFields(log.Fields{

		"ip":         info.IP,
		"enode":      info.Enode,
		"id":         info.ID,
		"listenAddr": info.ListenAddr,
		"name":       info.Name,
		"ports":      info.Ports,
		"protocols":  info.Protocols,
	}).Info("NodeInfo")

	select {}
}

// MyProtocol is a custom protocol
func MyProtocol() p2p.Protocol {
	return p2p.Protocol{ // 1.
		Name:    "MyProtocol", // 2.
		Version: 1,            // 3.
		Length:  1,            // 4.
		Run:     msgHandler,   // 5.
	}
}

const messageID = 0 // 1.

// Message is a message
type Message string // 2.

func msgHandler(peer *p2p.Peer, ws p2p.MsgReadWriter) error {

	for {
		msg, err := ws.ReadMsg() // 3.
		if err != nil {          // 4.
			return err // if reading fails return err which will disconnect the peer.
		}

		var myMessage [1]Message
		err = msg.Decode(&myMessage) // 5.
		if err != nil {
			// handle decode error
			continue
		}

		switch myMessage[0] {
		case "foo":
			err := p2p.SendItems(ws, messageID, "bar") // 6.
			if err != nil {
				return err // return (and disconnect) error if writing fails.
			}
		default:
			fmt.Println("recv:", myMessage)
		}
	}

	//return nil
}
