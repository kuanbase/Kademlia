package handle

import (
	"Kademlia/pkg/dht"
	"Kademlia/pkg/kencode"
	"net"
)

func Pong(conn net.Conn) error {
	response := kencode.NewEncoder().ResponsePing(conn.LocalAddr().String()).Encode()

	_, err := conn.Write([]byte(response))

	return err
}

func ReturnID(conn net.Conn, id dht.DhtID) error {
	response := kencode.NewEncoder().ResponseGETID(id).Encode()

	_, err := conn.Write([]byte(response))

	return err
}
