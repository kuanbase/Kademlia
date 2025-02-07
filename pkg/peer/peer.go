package peer

import (
	"Kademlia/pkg/dht"
	"Kademlia/pkg/global"
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type PeerNode struct {
	DhtNode        dht.DhtNode         // 本節點的Dht Node
	Address        net.TCPAddr         // 本節點的 IP Address
	DhtIDToAddress map[string]net.Addr // 將DhtID轉換成string 然後映射到對應用 IP Address
	// History        history.History     // 節點的一切歷史行為
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

	data, err := os.ReadFile(global.BootstrapNodeFilePath)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")

	for _, line := range lines {
		s := strings.Split(line, " ")

		if len(s) != 2 {
			return nil, errors.New("bootstape file format error")
		}

		sid := strings.TrimSpace(s[1])

		id, err := hex.DecodeString(sid)
		if err != nil {
			return nil, err
		}

		dhtNode.AddKBucket(id)
	}

	if netIP == nil {
		return nil, errors.New("invalid ip")
	}

	global.SystemPrintln(dhtNode)

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
