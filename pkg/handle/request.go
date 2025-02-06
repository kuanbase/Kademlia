package handle

import (
	"Kademlia/pkg/dht"
	"Kademlia/pkg/kencode"
	"Kademlia/pkg/peer"
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

func Store(peerNode *peer.PeerNode, id dht.DhtID) error {

	return nil
}

func FindNode(peerNode *peer.PeerNode) error {

	return nil
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
