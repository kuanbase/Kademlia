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
	DhtNode        dht.DhtNode            // 本節點的Dht Node
	Address        net.TCPAddr            // 本節點的 IP Address
	DhtIDToAddress map[string]net.TCPAddr // 將DhtID轉換成string 然後映射到對應用 IP Address
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

	if netIP == nil {
		return nil, errors.New("invalid ip")
	}

	data, err := os.ReadFile(global.BootstrapNodeFilePath)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")

	dhtIDToAddress := make(map[string]net.TCPAddr)

	for _, line := range lines {
		if line == "" {
			break
		}

		s := strings.Split(line, " ")

		if len(s) != 2 {
			return nil, errors.New("bootstape file format error")
		}

		bootstapeNodeAddress := strings.TrimSpace(s[0])

		address := strings.Split(bootstapeNodeAddress, ":")

		if len(address) != 2 {
			return nil, errors.New("bootstape file format error")
		}

		bootstapeIP := strings.TrimSpace(address[0])
		bootstapePort := strings.TrimSpace(address[1])
		port, err := strconv.Atoi(bootstapePort)
		if err != nil {
			return nil, errors.New("bootstape file format error")
		}

		bootstapeSid := strings.TrimSpace(s[1])

		bootstapeId, err := hex.DecodeString(bootstapeSid)
		if err != nil {
			return nil, err
		}

		bootstapeNetIp := net.ParseIP(bootstapeIP)
		if bootstapeNetIp == nil {
			return nil, errors.New("boostape invalid ip")
		}

		dhtIDToAddress[dht.DhtID(bootstapeId).ToString()] = net.TCPAddr{IP: bootstapeNetIp, Port: port}

		err = dhtNode.AddKBucket(dht.DhtID(bootstapeId))
		if err != nil {
			return nil, err
		}
	}

	global.SystemPrintln(dhtNode)

	return &PeerNode{DhtNode: *dhtNode, Address: net.TCPAddr{IP: netIP, Port: port}, DhtIDToAddress: dhtIDToAddress}, nil
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

func (peerNode *PeerNode) AddNode(id dht.DhtID, ip string, port int) error {
	netIP := net.ParseIP(ip)
	if netIP == nil {
		return errors.New("invalid ip")
	}

	peerNode.DhtIDToAddress[id.ToString()] = net.TCPAddr{IP: netIP, Port: port}

	return nil
}
