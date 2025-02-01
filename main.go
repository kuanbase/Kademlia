package main

import (
	"Kademlia/pkg/global"
	"Kademlia/pkg/kencode"
	"Kademlia/pkg/peer"
	"bytes"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
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

	go func() {
		var command, value string

		for {
			fmt.Print("Kademlia Commander> ")
			_, err = fmt.Scanf("%s %s", &command, &value)

			if err != nil {
				log.Println("???? " + err.Error())
			}

			switch command {
			case "ping":
				addr := strings.Split(value, ":")
				ip := addr[0]
				_, err := strconv.Atoi(addr[1])

				if err != nil {
					log.Println(err)
					continue
				}

				if !global.ValidateIPAddress(ip) {
					log.Println("Kademlia: Invalid IP address, Please Try it Again")
					continue
				}

				err = peerNode.Ping(value)
				if err != nil {
					log.Println(err)
				}
			}

			if err != nil {
				log.Println(err)
			}
		}
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
		}

		// 接收信號，並返回信號
		go func() {
			defer conn.Close()

			var buffer bytes.Buffer

			_, err := buffer.ReadFrom(conn)
			if err != nil {
				log.Println(err)
			}

			kenCode := kencode.NewDecoder(buffer.String()).Decode()

			for i := 0; i < len(kenCode.Commands); i++ {
				switch kenCode.Commands[i] {
				case "PING":
					fmt.Println("Kademlia> " + conn.RemoteAddr().String() + " PONG")
					response := kencode.NewEncoder().ResponsePing().Encode()
					_, err := conn.Write([]byte(response))
					if err != nil {
						log.Println(err)
					}
				default:
					log.Printf("Unknown command: %s", kenCode.Commands[i])
				}
			}
		}()
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

	bootstrapNodeFilename := filepath.Join(global.KademliaNodesPath, "bootstrap_nodes.txt")

	_, err := os.Stat(bootstrapNodeFilename)
	if os.IsNotExist(err) {
		_, _ = os.Create(bootstrapNodeFilename)
	}

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
