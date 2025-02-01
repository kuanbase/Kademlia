package peer

import (
	"Kademlia/pkg/dht"
	"Kademlia/pkg/global"
	"Kademlia/pkg/kencode"
	"bytes"
	"encoding/json"
	"errors"
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
	DhtIDToAddress map[string]net.Addr // 將DhtID 轉換成string 然後映射到對應用 IP Address
}

func NewPeerNode(address string) (*PeerNode, error) {
	dhtNode, err := dht.NewDhtNode()
	if err != nil {
		return nil, err
	}

	addr := strings.Split(address, ":")
	ip := addr[0]
	port, err := strconv.Atoi(addr[1])

	if err != nil {
		return nil, err
	}

	netIP := net.ParseIP(ip)

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
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	defer conn.Close()

	err = conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		return err
	}

	msg := kencode.NewEncoder().Ping(address).Encode()

	_, err = conn.Write([]byte(msg))
	if err != nil {
		return err
	}

	return nil
}

func (peerNode *PeerNode) Store() error {

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
