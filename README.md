# Go Point-to-Point VPN

This project creates a **very simple point-to-point VPN** using Go + TUN interfaces.

- **Server runs on Kali Linux (Virtual Machine)**
- **Client runs on macOS**

Traffic is sent between TUN interfaces over a TCP connection.

Both devices will get internal VPN IPs:

| Device | Interface | IP |
|--------|-----------|------|
| **Kali server** | `vpn0` or `tun0` | `10.2.0.1/24` |
| **macOS client** | `utunX` | `10.2.0.2/24` |

After setup, both sides can **ping each other through VPN**.

---

## 📌 Requirements

### On Kali (server)
- Go installed (`go version`)
- `github.com/songgao/water` library

### On macOS (client)
- Go installed
- `github.com/songgao/water` library

---

## 📁 File Structure

```
vpn_server.go
vpn_client.go
```

---

## 🚀 Setup Instructions

### 1. Install Go (if missing)

#### Kali
```bash
sudo apt update
sudo apt install golang-go -y
```

#### macOS
```bash
brew install go
```

### 2. 📦 Install TUN library (both sides)

Inside each project folder:

```bash
go mod init vpn
go get github.com/songgao/water
```

---

## 🖥️ SERVER SETUP (KALI)

Build the server:

```bash
go build -o vpnServer vpn_server.go
```

Run the server as root:

```bash
sudo ./vpnServer
```

You should see:

```
Server TUN device created: vpn0
Server TUN set up with IP 10.2.0.1/24
Server listening on :8000
```

---

## 💻 CLIENT SETUP (macOS)


Build the client:

```bash
go build -o vpnClient vpn_client.go
```

Run client as root with server IP and port flags:

```bash
sudo ./vpnClient -ip <SERVER_IP> -p <PORT>
```

✅ Example:

```bash
sudo ./vpnClient -ip 192.168.1.12 -p 8000
```

If successful you'll see:

```
Client TUN device created: utun4
Client TUN set up with IP 10.2.0.2/24
Connected to VPN server: 192.168.1.12:8000
```

### 🔧 If macOS didn't configure IP automatically

Run (replace `utun4` with your actual interface name):

```bash
sudo ifconfig utun4 10.2.0.2 10.2.0.1 netmask 255.255.255.0 up
```

---

## 🧪 TESTING (Ping)

### From macOS → Kali

```bash
ping 10.2.0.1
```

### From Kali → macOS

```bash
ping 10.2.0.2
```

### 🚑 If ping fails (macOS)

Add a route manually:

```bash
sudo route add 10.2.0.0/24 -interface <interface-name>
```

Example:

```bash
sudo route add 10.2.0.0/24 -interface utun4
```

---

## 🎉 If both pings work, your Go VPN is successfully running!
