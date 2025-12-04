package main

import (
	"fmt"
	"log"
	"net"
	"os/exec"

	"github.com/songgao/water"
)

const (
	SERVER_ADDR = "192.168.1.13:8000" // replace with your Kali IP
	CLIENT_IP   = "10.2.0.2/24"
)

var simpleKey = []byte("alloftheabove")
var addNum byte = 10

func main() {
	// 1️⃣ Create TUN interface
	cfg := water.Config{DeviceType: water.TUN}
	iface, err := water.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Client TUN device created:", iface.Name())

	// 2️⃣ Assign IP
	err = execCommand("ifconfig", iface.Name(), "10.2.0.2", "10.2.0.2", "netmask", "255.255.255.0", "up")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Client TUN set up with IP:", CLIENT_IP)

	// 3️⃣ Connect to server
	conn, err := net.Dial("tcp", SERVER_ADDR)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to VPN server")

	fmt.Println("If ping fails, add route: sudo route add 10.2.0.0/24 -interface", iface.Name())

	// 4️⃣ TUN → TCP
	go func() {
		buf := make([]byte, 1500)
		for {
			n, err := iface.Read(buf)
			if err != nil {
				log.Println("TUN read error:", err)
				return
			}
			conn.Write(simpleEncrypt(buf[:n]))
		}
	}()

	// 5️⃣ TCP → TUN
	buf := make([]byte, 1500)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Println("TCP read error:", err)
			return
		}
		iface.Write(simpleDecrypt(buf[:n]))
	}
}

// helper to run shell commands
func execCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stderr = nil
	cmd.Stdout = nil
	return cmd.Run()
}

// simple encryption: add + XOR
func simpleEncrypt(data []byte) []byte {
	res := make([]byte, len(data))
	for i, b := range data {
		res[i] = (b + addNum) ^ simpleKey[i%len(simpleKey)]
	}
	return res
}

// simple decryption: reverse of encryption
func simpleDecrypt(data []byte) []byte {
	res := make([]byte, len(data))
	for i, b := range data {
		res[i] = (b ^ simpleKey[i%len(simpleKey)]) - addNum
	}
	return res
}
