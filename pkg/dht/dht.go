package dht

import (
	"crypto/rand"
	"crypto/sha1"
	"errors"
	"io"
)

const (
	BucketSize = 160
	DhtIDSize  = 160
)

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

func (dhtNode *DhtNode) AddKBucket(id DhtID) error {
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
