package main

import (
	"crypto/rand"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
)

const (
	BucketSize            = 160
	DhtIDSize             = 160
	KademliaDirectoryPath = "~/kademlia"
	KademliaFilesPath     = "~/kademlia/files"
)

var BootstrapNode = []string{
	"",
}

type DhtID []byte

func (d DhtID) String() string {
	return string(d)
}

func NewDhtID() (DhtID, error) {
	nodeID, err := generateNodeID()
	if err != nil {
		return nil, err
	}

	return nodeID, nil
}
func (d DhtID) XOR(id DhtID) (DhtID, error) {
	if len(d) != DhtIDSize || len(id) != DhtIDSize {
		return nil, errors.New("invalid DhtID")
	}

	result := make([]byte, len(d))

	for i := 0; i < len(result); i++ {
		result[i] = id[i] ^ d[i]
	}

	return result, nil
}

type DhtNode struct {
	ID       DhtID
	KBuckets [BucketSize][]DhtID
}

func generateNodeID() ([]byte, error) {
	randomBytes := make([]byte, BucketSize)

	_, err := io.ReadFull(rand.Reader, randomBytes)
	if err != nil {
		return nil, err
	}

	hash := sha1.New()
	hash.Write(randomBytes)
	nodeID := hash.Sum(nil)

	return nodeID, nil
}

func NewDhtNode() (*DhtNode, error) {
	nodeID, err := NewDhtID()
	if err != nil {
		return nil, err
	}

	return &DhtNode{ID: nodeID}, nil
}

func (dhtNode *DhtNode) Distance(id DhtID) (int, error) {
	result, err := dhtNode.ID.XOR(id)

	if err != nil {
		return 0, err
	}

	k := DhtIDSize

	for i := 0; i < len(result); i++ {
		if result[i] != id[k] {
			break
		}

		k--
	}

	return k, nil
}

func (dhtNode *DhtNode) addKBucket(id DhtID) error {
	distance, err := dhtNode.Distance(id)
	if err != nil {
		return err
	}

	if dhtNode.KBuckets[distance-1] == nil {
		dhtNode.KBuckets[distance-1] = make([]DhtID, 0)
	}

	dhtNode.KBuckets[distance-1] = append(dhtNode.KBuckets[distance-1], id)

	return nil
}

func (dhtNode *DhtNode) FindKBucketByDhtID(id DhtID) ([]DhtID, error) {
	distance, err := dhtNode.Distance(id)
	if err != nil {
		return nil, err
	}

	if dhtNode.KBuckets[distance-1] == nil {
		return nil, errors.New("not found")
	}

	return dhtNode.KBuckets[distance-1], nil
}

type PeerNode struct {
	DhtNode        DhtNode             // 本節點的Dht Node
	Address        net.Addr            // 本節點的 IP Address
	DhtIDToAddress map[string]net.Addr // 將DhtID 轉換成string 然後映射到對應用 IP Address
}

func NewPeerNode(ip string) (*PeerNode, error) {
	dhtNode, err := NewDhtNode()
	if err != nil {
		return nil, err
	}

	return &PeerNode{DhtNode: *dhtNode, Address: &net.TCPAddr{IP: net.IP(ip)}}, nil
}

func (peerNode *PeerNode) Ping() error {

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

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./kademlia [create <ip> | join <bootstrap ips...> | ping <ip> | store <data> | find_node <DhtID> | find_value <file_hash_key>]")
		return
	}

	fmt.Println("OK")
}
