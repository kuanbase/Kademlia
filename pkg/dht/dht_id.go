package dht

import "errors"

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
