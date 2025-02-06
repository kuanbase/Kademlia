package handle

import (
	"Kademlia/pkg/global"
	"Kademlia/pkg/kencode"
	"Kademlia/pkg/peer"
	"log"
	"net"
	"strconv"
	"strings"
)

// Server 處理信號
func Server(conn net.Conn) {
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
			response := kencode.NewEncoder().ResponsePing().Encode()
			_, err := conn.Write([]byte(response))
			if err != nil {
				log.Println(err)
			}
		default:
			log.Printf("\rUnknown command: %s\n", kenCode.Commands[i])
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

			err = peerNode.Ping(value)
			if err != nil {
				global.ErrPrintln(err.Error())
			}
		}
	}
}
