package handle

import (
	"Kademlia/pkg/dht"
	"Kademlia/pkg/global"
	"Kademlia/pkg/kencode"
	"Kademlia/pkg/peer"
	"encoding/hex"
	"log"
	"net"
	"strconv"
	"strings"
)

// Server 處理信號
func Server(peerNode *peer.PeerNode, conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 10000)

	_, err := conn.Read(buf)
	if err != nil {
		log.Println(err)
	}

	kenCode := kencode.NewDecoder(string(buf)).Decode()

	for i := 0; i < len(kenCode.Commands); i++ {
		switch kenCode.Commands[i] {
		case kencode.PING:
			err = Pong(conn)
			if err != nil {
				global.ErrPrintln(err.Error())
				continue
			}
		case kencode.GETID:
			err = ReturnID(conn, peerNode.DhtNode.ID)
			if err != nil {
				global.ErrPrintln(err.Error())
				continue
			}
		case kencode.FINDNODE:
			id, ok := kenCode.Values[i].(dht.DhtID)

			if !ok {
				global.ErrPrintln("Dht Id Convert Error")
				continue
			}

			FindNode(peerNode, id)

		default:
			global.SystemPrintln("Unknown command: " + kenCode.Commands[i])
		}
	}
}

// Cli 輸入指令
func Cli(peerNode *peer.PeerNode) {
	for {
		command, value := global.CmdInput()

		if command == "" || value == "" {
			global.ErrPrintln("Please don't enter empty command or empty value.")
			continue
		}

		switch command {
		case "ping":
			addr := strings.Split(value, ":")

			if len(addr) != 2 {
				global.ErrPrintln("Please enter the address like <ip>:<port>")
				continue
			}

			ip := addr[0]
			_, err := strconv.Atoi(addr[1])

			if err != nil {
				global.ErrPrintln("Please enter validation port.")
				continue
			}

			if !global.ValidateIPAddress(ip) {
				global.ErrPrintln("Please enter validation IP address.")
				continue
			}

			kenCode, err := Ping(peerNode, value)
			if err != nil {
				global.ErrPrintln(err.Error())
			}

			for i := 0; i < len(kenCode.Commands); i++ {
				switch kenCode.Commands[i] {
				case kencode.PONG:
					address, ok := kenCode.Values[i].(string)

					if !ok {
						global.ErrPrintln("Please enter the address like <ip>:<port>")
						continue
					}

					global.PongPrintln(address)
				default:
					global.SystemPrintln("Unknown command: " + kenCode.Commands[i])
				}
			}
		case "getid":
			addr := strings.Split(value, ":")

			if len(addr) != 2 {
				global.ErrPrintln("Please enter the address like <ip>:<port>")
				continue
			}

			ip := addr[0]
			_, err := strconv.Atoi(addr[1])

			if err != nil {
				global.ErrPrintln("Please enter validation port.")
				continue
			}

			if !global.ValidateIPAddress(ip) {
				global.ErrPrintln("Please enter validation IP address.")
				continue
			}

			kenCode, err := GetID(peerNode, value)

			if err != nil {
				global.ErrPrintln(err.Error())
				continue
			}

			for i := 0; i < len(kenCode.Commands); i++ {
				switch kenCode.Commands[i] {
				case kencode.RETURNID:
					sid, ok := kenCode.Values[i].(string)

					if !ok {
						global.ErrPrintln("Please enter the dht id")
						continue
					}

					_, err := hex.DecodeString(sid)
					if err != nil {
						global.ErrPrintln(err.Error())
						continue
					}

					global.DhtSidPrintln(sid)
				default:
					global.SystemPrintln("Unknown command: " + kenCode.Commands[i])
				}
			}
		case "findnode":
			id, err := hex.DecodeString(value)
			if err != nil {
				global.ErrPrintln(err.Error())
				continue
			}

			kenCode, err := FindNode(peerNode, id)
			if err != nil {
				global.ErrPrintln(err.Error())
				continue
			}

			for i := 0; i < len(kenCode.Commands); i++ {
				switch kenCode.Commands[i] {
				case kencode.RETURNNODE:
					address, ok := kenCode.Values[i].(string)

					if !ok {
						global.ErrPrintln("RETURN NODE ERROR")
						continue
					}

					global.AddressPrintln(address)
				case kencode.RETURNNULL:
					global.AddressPrintln("not found")
				default:
					global.SystemPrintln("Unknown command: " + kenCode.Commands[i])
				}
			}
		}
	}
}
