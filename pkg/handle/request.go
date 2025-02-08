package handle

import (
	"Kademlia/pkg/dht"
	"Kademlia/pkg/kencode"
	"Kademlia/pkg/peer"
	"errors"
	"net"
	"time"
)

func Ping(peerNode *peer.PeerNode, address string) (*kencode.KenCode, error) {
	conn, err := net.DialTimeout("tcp", address, 10*time.Second)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	msg := kencode.NewEncoder().Ping(address).Encode()

	_, err = conn.Write([]byte(msg))
	if err != nil {
		return nil, err
	}

	buf := make([]byte, 10000)

	// conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	_, err = conn.Read(buf)
	if err != nil {
		return nil, err
	}

	kenCode := kencode.NewDecoder(string(buf)).Decode()

	return kenCode, nil
}

func GetID(peerNode *peer.PeerNode, address string) (*kencode.KenCode, error) {
	conn, err := net.DialTimeout("tcp", address, 10*time.Second)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	msg := kencode.NewEncoder().GetID(address).Encode()

	_, err = conn.Write([]byte(msg))
	if err != nil {
		return nil, err
	}

	buf := make([]byte, 10000)

	_, err = conn.Read(buf)
	if err != nil {
		return nil, err
	}

	kenCode := kencode.NewDecoder(string(buf)).Decode()

	return kenCode, nil
}

func Store(peerNode *peer.PeerNode, data []byte) (*kencode.KenCode, error) {
	// hash := sha1.New()
	// _, err := hash.Write(data)
	// if err != nil {
	// 	return err
	// }

	// storeCode := hash.Sum(data)

	return nil, nil
}

func FindNode(peerNode *peer.PeerNode, id dht.DhtID) (*kencode.KenCode, error) {
	distance, err := peerNode.DhtNode.Distance(id)
	if err != nil {
		return nil, err
	}

	min := distance

	for i := 0; i < len(peerNode.DhtNode.KBuckets); i++ {
		if peerNode.DhtNode.KBuckets[i] != nil && min > i {
			min = i
			break
		}
	}

	nodeIds := peerNode.DhtNode.KBuckets[min]

	if nodeIds == nil {
		msg := kencode.NewEncoder().ResponseFindNode("").Encode()
		kenCode := kencode.NewDecoder(msg).Decode()
		return kenCode, nil
	}

	for _, nodeID := range nodeIds {
		nodeAddr, ok := peerNode.DhtIDToAddress[nodeID.ToString()]
		if ok {
			if nodeID.Equal(id) {
				msg := kencode.NewEncoder().ResponseFindNode(nodeAddr.String()).Encode()
				kenCode := kencode.NewDecoder(msg).Decode()
				return kenCode, nil
			}

			conn, err := net.DialTimeout("tcp", nodeAddr.String(), 10*time.Second)
			if err != nil {
				return nil, err
			}
			defer conn.Close()

			msg := kencode.NewEncoder().FindNode(nodeID).Encode()

			_, err = conn.Write([]byte(msg))
			if err != nil {
				return nil, err
			}

			buf := make([]byte, 10000)

			conn.SetReadDeadline(time.Now().Add(10 * time.Second))
			_, err = conn.Read(buf)
			if err != nil {
				return nil, err
			}

			kenCode := kencode.NewDecoder(string(buf)).Decode()

			return kenCode, nil
		}
	}

	return nil, errors.New("not found")
}

func FindValue(peerNode *peer.PeerNode) error {

	return nil
}

func Download(peerNode *peer.PeerNode) error {

	return nil
}

func Upload(peerNode *peer.PeerNode) error {

	return nil
}
