package dht

import (
	"encoding/hex"
	"errors"
)

type DhtID []byte

func (d DhtID) ToString() string {
	sid := hex.EncodeToString(d)
	return sid
}

func NewDhtID() (DhtID, error) {
	nodeID, err := generateNodeID()
	if err != nil {
		return nil, err
	}

	return nodeID, nil
}
func (d DhtID) XOR(id DhtID) (DhtID, error) {
	if len(d) != DhtIDBytesSize || len(id) != DhtIDBytesSize {
		return nil, errors.New("invalid DhtID")
	}

	result := make([]byte, len(d))

	for i := 0; i < len(result); i++ {
		result[i] = id[i] ^ d[i]
	}

	return result, nil
}

func (d DhtID) Equal(id DhtID) bool {
	if len(d) != len(id) {
		return false
	}

	for i := 0; i < len(d); i++ {
		if d[i] != id[i] {
			return false
		}
	}

	return true
}
