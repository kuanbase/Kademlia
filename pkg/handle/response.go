package handle

import (
	"Kademlia/pkg/dht"
	"Kademlia/pkg/kencode"
	"fmt"
	"net"
)

func Pong(conn net.Conn) error {
	response := kencode.NewEncoder().ResponsePing(conn.LocalAddr().String()).Encode()

	_, err := conn.Write([]byte(response))

	return err
}

func ReturnID(conn net.Conn, id dht.DhtID) error {
	fmt.Printf("\rDebug> before encode: %v\n", id)

	response := kencode.NewEncoder().ResponseGETID(id).Encode()

	fmt.Printf("\rDebug> after encode: %v\n", id)

	_, err := conn.Write([]byte(response))

	return err
}
