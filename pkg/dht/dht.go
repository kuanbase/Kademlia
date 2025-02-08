package dht

import (
	"crypto/rand"
	"crypto/sha1"
	"io"
)

const (
	BucketBitsSize  = 160
	BucketBytesSize = BucketBitsSize / 8
	DhtIDBytesSize  = BucketBitsSize / 8
)

type DhtNode struct {
	ID       DhtID
	KBuckets [BucketBitsSize][]DhtID
}

func generateNodeID() ([]byte, error) {
	randomBytes := make([]byte, DhtIDBytesSize)

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

	k := BucketBitsSize

	for i := 0; i < len(result); i++ {
		x := result[i] ^ id[i]

		for j := 0; j < 8; j++ {
			b := x & 1

			if b == 0 {
				k--
			} else {
				return k, nil
			}

			x >>= 1
		}
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
