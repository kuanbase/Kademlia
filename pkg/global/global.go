package global

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	KademliaDirectoryName      = "kademlia"
	KademliaFilesDirectoryName = "kademlia/files"
	KademliaNodesDirectoryName = "kademlia/nodes"
)

var HomeDirectory string
var KademliaDirectoryPath string
var KademliaNodesPath string
var KademliaFilesPath string

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	HomeDirectory = home
	KademliaDirectoryPath = filepath.Join(home, KademliaDirectoryName)
	KademliaNodesPath = filepath.Join(home, KademliaNodesDirectoryName)
	KademliaFilesPath = filepath.Join(home, KademliaFilesDirectoryName)
}

func ValidateIPAddress(address string) bool {
	ips := strings.Split(address, ".")

	if len(ips) != 4 {
		return false
	}

	for _, s := range ips {
		number, err := strconv.Atoi(s)
		if err != nil {
			return false
		}

		if number < 0 || number > 255 {
			return false
		}
	}

	return true
}
