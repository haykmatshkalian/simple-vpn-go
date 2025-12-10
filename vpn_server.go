package main

import (
	"fmt"
	"log"
	"net"
	"os/exec"

	"github.com/songgao/water"
)

var simpleKey = []byte("mysecret")
var addNum byte = 3

func main() {
	// 1️⃣ Create TUN
	cfg := water.Config{DeviceType: water.TUN}
	iface, err := water.New(cfg)
	if err != nil {
		log.Fatal("TUN create:", err)
	}
	fmt.Println("Server TUN:", iface.Name())

	// 2️⃣ Assign IP (same subnet as client)
	serverIP := "10.2.0.1/24"
	err = exec.Command("ip", "addr", "add", serverIP, "dev", iface.Name()).Run()
	if err != nil {
		log.Fatal("IP assign:", err)
	}
	exec.Command("ip", "link", "set", "dev", iface.Name(), "up").Run()
	fmt.Println("Server TUN configured:", serverIP)

	// 3️⃣ Enable IP forwarding (Linux)
	exec.Command("sysctl", "-w", "net.ipv4.ip_forward=1").Run()

	// 4️⃣ Start TCP listener
	ln, err := net.Listen("tcp", "0.0.0.0:8000")
	if err != nil {
		log.Fatal("Listen:", err)
	}
	fmt.Println("Server listening on 0.0.0.0:8000 ...")

	conn, err := ln.Accept()
	if err != nil {
		log.Fatal("Accept:", err)
	}
	fmt.Println("Client connected!")

	// 5️⃣ Forward TUN → TCP
	go func() {
		buf := make([]byte, 1500)
		for {
			n, err := iface.Read(buf)
			if err != nil {
				log.Println("TUN read:", err)
				return
			}
			conn.Write(simpleEncrypt(buf[:n]))
		}
	}()

	// 6️⃣ Forward TCP → TUN
	buf := make([]byte, 1500)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Println("TCP read:", err)
			return
		}
		iface.Write(simpleDecrypt(buf[:n]))
	}
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
