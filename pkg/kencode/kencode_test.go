package kencode

import (
	"fmt"
	"testing"
)

func TestEncoder(t *testing.T) {
	encoder := NewEncoder()
	s := encoder.Ping("192.168.1.35:8090").Encode()
	fmt.Println(s)
	decoder := NewDecoder(s)
	kenCode := decoder.Decode()
	fmt.Println(kenCode)
}
