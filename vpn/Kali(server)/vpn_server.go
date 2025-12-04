package main

import (
	"fmt"
	"log"
	"net"
	"os/exec"

	"github.com/songgao/water"
)

func main() {
	// 1) Create TUN
	cfg := water.Config{
		DeviceType: water.TUN,
	}
	iface, err := water.New(cfg)
	if err != nil {
		log.Fatal("TUN create:", err)
	}
	fmt.Println("Server TUN:", iface.Name())

	// 2) Assign IP using ip command (Kali)
	//    ifconfig sometimes breaks on modern systems
	err = exec.Command("ip", "addr", "add", "10.0.0.1/24", "dev", iface.Name()).Run()
	if err != nil {
		log.Fatal("IP assign:", err)
	}
	exec.Command("ip", "link", "set", "dev", iface.Name(), "up").Run()

	fmt.Println("Server TUN configured: 10.0.0.1/24")

	// 3) Start TCP listener
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

	// 4) Forward TUN → TCP
	go func() {
		buf := make([]byte, 1500)
		for {
			n, err := iface.Read(buf)
			if err != nil {
				log.Println("TUN read:", err)
				return
			}
			_, err = conn.Write(buf[:n])
			if err != nil {
				log.Println("TCP write:", err)
				return
			}
		}
	}()

	// 5) Forward TCP → TUN
	buf := make([]byte, 1500)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Println("TCP read:", err)
			return
		}
		_, err = iface.Write(buf[:n])
		if err != nil {
			log.Println("TUN write:", err)
			return
		}
	}
}


