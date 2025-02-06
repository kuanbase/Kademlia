package peer

import (
	"Kademlia/pkg/dht"
	"Kademlia/pkg/global"
	"Kademlia/pkg/history"
	"Kademlia/pkg/kencode"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type PeerNode struct {
	DhtNode        dht.DhtNode         // 本節點的Dht Node
	Address        net.TCPAddr         // 本節點的 IP Address
	DhtIDToAddress map[string]net.Addr // 將DhtID轉換成string 然後映射到對應用 IP Address
	History        history.History     // 節點的一切歷史行為
}

func NewPeerNode(address string) (*PeerNode, error) {
	dhtNode, err := dht.NewDhtNode()
	if err != nil {
		return nil, err
	}

	addr := strings.Split(address, ":")

	if len(addr) != 2 {
		return nil, errors.New("please enter the name like <ip>:<port>")
	}

	ip := addr[0]
	port, err := strconv.Atoi(addr[1])

	if err != nil {
		return nil, err
	}

	netIP := net.ParseIP(ip)

	// data, err := os.ReadFile(global.BootstrapNodeFilePath)

	// lines := strings.Split(string(data), "\n")

	// for _, line := range lines {
	// 	dhtNode.AddKBucket()
	// }

	if netIP == nil {
		return nil, errors.New("invalid ip")
	}

	return &PeerNode{DhtNode: *dhtNode, Address: net.TCPAddr{IP: netIP, Port: port}}, nil
}

func NewPeerNodeByPeerFile(filename string) (*PeerNode, error) {
	var peerNode PeerNode

	nodePath := filepath.Join(global.KademliaNodesPath, filename)

	peerFileData, err := os.ReadFile(nodePath)
	if err != nil {
		return nil, err
	}

	err = peerNode.Unmarshal(peerFileData)
	if err != nil {
		return nil, err
	}

	return &peerNode, nil
}

func (peerNode *PeerNode) Marshal() ([]byte, error) {
	var buf bytes.Buffer

	err := json.NewEncoder(&buf).Encode(*peerNode)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (peerNode *PeerNode) Unmarshal(data []byte) error {
	buf := bytes.NewBuffer(data)

	if peerNode == nil {
		return errors.New("nil peerNode")
	}

	err := json.NewDecoder(buf).Decode(peerNode)
	if err != nil {
		return err
	}

	return nil
}

func (peerNode *PeerNode) Ping(address string) error {
	conn, err := net.DialTimeout("tcp", address, 10*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()

	msg := kencode.NewEncoder().Ping(address).Encode()

	_, err = conn.Write([]byte(msg))
	if err != nil {
		return err
	}

	buf := make([]byte, 10000)

	// conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	_, err = conn.Read(buf)
	if err != nil {
		return err
	}

	kenCode := kencode.NewDecoder(string(buf)).Decode()

	for i := 0; i < len(kenCode.Commands); i++ {
		switch kenCode.Commands[i] {
		case kencode.PONG:
			fmt.Printf("\r%s> %s: PONG\n", conn.RemoteAddr().String(), kenCode.Values[i])
		default:
			log.Printf("\rUnknown command: %s\n", kenCode.Commands[i])
		}
	}

	return nil
}

func (peerNode *PeerNode) Store(id dht.DhtID) error {

	return nil
}

func (peerNode *PeerNode) FindNode() error {

	return nil
}

func (peerNode *PeerNode) FindValue() error {

	return nil
}

func (peerNode *PeerNode) Download() error {

	return nil
}

func (peerNode *PeerNode) Upload() error {

	return nil
}
