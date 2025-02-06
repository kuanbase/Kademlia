package main

import (
	"Kademlia/pkg/global"
	"Kademlia/pkg/handle"
	"Kademlia/pkg/peer"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
)

var BootstrapNode = []string{
	"",
}

var create = flag.String("create", "", "create a new node")
var join = flag.String("join", "", "join nodes network")
var ping = flag.String("ping", "", "ping node")
var store = flag.String("store", "", "store data")
var findNode = flag.String("find_node", "", "find a node")
var findValue = flag.String("find_value", "", "find a value")
var run = flag.String("run", "", "run node")

func Run(peerNode *peer.PeerNode) {
	listener, err := net.ListenTCP("tcp", &peerNode.Address)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	fmt.Println("Kademlia: Listening on " + peerNode.Address.String())

	// 輸入指令
	go handle.Cli(peerNode)

	// 監聽每一個連接
	for {
		conn, err := listener.Accept()
		if err != nil {
			global.ErrPrintln(err.Error())
		}

		// 處理信號
		go handle.Server(conn)
	}
}

func main() {
	_ = os.Mkdir(global.KademliaDirectoryPath, 0777)
	_ = os.Mkdir(global.KademliaNodesPath, 0777)
	_ = os.Mkdir(global.KademliaFilesPath, 0777)

	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, "Usage: \n")
		_, _ = fmt.Fprintf(os.Stderr, "	./kademlia [create <ip> "+
			"| join <bootstrap_nodes filepath> "+
			"| store <data>] "+
			"| find_node <DhtID> "+
			"| find_value <FileHash> "+
			"| run <DhtID>\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *create != "" {
		fmt.Println("Create node:", *create)
		node, err := peer.NewPeerNode(*create)
		if err != nil {
			log.Fatal(err)
		}

		*create = strings.Replace(*create, ":", "_", -1)

		filePath := filepath.Join(global.KademliaNodesPath, *create)

		data, err := node.Marshal()
		if err != nil {
			log.Fatal(err)
		}

		f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		_, err = f.Write(data)
		if err != nil {
			log.Fatal(err)
		}

		return
	}

	if *join != "" {
		fmt.Println("Join node:", *join)
		return
	}

	if *ping != "" {
		fmt.Println("Ping node:", *ping)
		return
	}

	if *store != "" {
		fmt.Println("Store node:", *store)
		return
	}

	if *findNode != "" {
		fmt.Println("Find node:", *findNode)
		return
	}

	if *findValue != "" {
		fmt.Println("Find value:", *findValue)
		return
	}

	if *run != "" {
		fmt.Println("Run node:", *run)

		*run = strings.Replace(*run, ":", "_", -1)

		node, err := peer.NewPeerNodeByPeerFile(*run)

		if err != nil {
			log.Fatal(err)
		}

		Run(node)

		return
	}
}
