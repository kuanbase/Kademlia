package handle

import (
	"Kademlia/pkg/dht"
	"Kademlia/pkg/kencode"
	"encoding/hex"
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

	sid := hex.EncodeToString(id)

	response := kencode.NewEncoder().ResponseGETID(sid).Encode()

	fmt.Printf("\rDebug> after encode: %v\n", response)

	_, err := conn.Write([]byte(response))

	return err
}
