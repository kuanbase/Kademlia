package global

import (
	"Kademlia/pkg/dht"
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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
var BootstrapNodeFilePath string

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	HomeDirectory = home
	KademliaDirectoryPath = filepath.Join(home, KademliaDirectoryName)
	KademliaNodesPath = filepath.Join(home, KademliaNodesDirectoryName)
	KademliaFilesPath = filepath.Join(home, KademliaFilesDirectoryName)

	_ = os.Mkdir(KademliaDirectoryPath, 0777)
	_ = os.Mkdir(KademliaNodesPath, 0777)
	_ = os.Mkdir(KademliaFilesPath, 0777)

	BootstrapNodeFilePath = filepath.Join(KademliaNodesPath, "bootstrap_nodes.txt")

	_, err = os.Stat(BootstrapNodeFilePath)
	if os.IsNotExist(err) {
		_, _ = os.Create(BootstrapNodeFilePath)
	}

	bootstrapNodes := []string{""}

	f, err := os.OpenFile(BootstrapNodeFilePath, os.O_WRONLY|os.O_TRUNC, 0777)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	_, err = f.Write([]byte(strings.Join(bootstrapNodes, "\n")))
	if err != nil {
		log.Fatal(err)
	}
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

func ValidatePort(port int) bool {
	if port < 0 || port > 65535 {
		return false
	}

	return true
}

func Println(s string) {
	fmt.Printf("\r%s\n", s)
}

func CmdInput() (string, string) {
	fmt.Print("Kademlia Commander> ")
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)

	if line == "Exit" || line == "exit" {
		ExitPrintln("Ok!")
		os.Exit(0)
	}

	s := strings.Split(line, " ")

	if len(s) != 2 {
		return "", ""
	}

	command, value := strings.TrimSpace(s[0]), strings.TrimSpace(s[1])

	return command, value
}

func ErrPrintln(e string) {
	fmt.Printf("\rError> %s\n", e)
}

func ExitPrintln(s string) {
	fmt.Printf("\rExit> %s\n", s)
}

func SystemPrintln(s any) {
	fmt.Printf("\rSystem> %v\n", s)
}

func PongPrintln(address string) {
	fmt.Printf("\r%s> PONG\n", address)
}

func DhtIdPrintln(id dht.DhtID) {
	fmt.Printf("\rReturn ID> %v\n", id)
}

func getDefaultGatewayInterface() (string, error) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("route", "print")
	case "darwin":
		cmd = exec.Command("route", "-n", "get", "default")
	case "linux":
		cmd = exec.Command("ip", "route", "show", "default")
	default:
		return "", fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error executing command: %v, output: %s", err, string(out))
	}

	output := string(out)

	switch runtime.GOOS {
	case "windows":
		lines := strings.Split(output, "\n")
		for _, line := range lines {
			if strings.Contains(line, " 0.0.0.0 ") {
				fields := strings.Fields(line)
				if len(fields) > 3 {
					return fields[3], nil // Interface
				}
			}
		}
	case "darwin":
		lines := strings.Split(output, "\n")
		for _, line := range lines {
			if strings.Contains(line, "interface:") {
				parts := strings.Split(line, ":")
				if len(parts) > 1 {
					return strings.TrimSpace(parts[1]), nil
				}
			}
		}
	case "linux":
		lines := strings.Split(output, "\n")
		for _, line := range lines {
			if strings.Contains(line, "dev") {
				fields := strings.Fields(line)
				if len(fields) > 4 {
					return fields[4], nil // Interface
				}
			}
		}

	}

	return "", fmt.Errorf("could not determine default gateway interface")
}

func GetWifiIPV4Address() {
	ifaceName, err := getDefaultGatewayInterface()
	if err != nil {
		fmt.Println("Error getting default interface:", err)
		return
	}

	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		fmt.Printf("找不到介面 %s: %v\n", ifaceName, err)
		return
	}

	addrs, err := iface.Addrs()
	if err != nil {
		fmt.Println("Error getting addresses for interface:", err)
		return
	}

	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			fmt.Println("IP Address:", ipNet.IP.String())
			return //找到第一個IPv4地址就結束
		}
	}

	fmt.Printf("介面 %s 沒有 IPv4 位址\n", ifaceName)
}
