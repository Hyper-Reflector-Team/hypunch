package main

import (
	"encoding/json"
	"log"
	"net"
	"sync"

	"github.com/google/uuid"
)

const HOLE_PUNCH_SERVER_PORT int = 33333 // This can be any port you want to run the server on
const HOST_IP string = "0.0.0.0"

// A Peer is a client outside of the server.
type Peer struct {
	UID     string `json:"uid"`
	Address string `json:"address"`
	Port    int    `json:"port"`
}

// Message is an incoming JSON message
type Message struct {
	UID     string `json:"uid"`
	PeerUID string `json:"peerUid"`
	Kill    bool   `json:"kill"`
}

// MatchPayload is sent to both peers
type MatchPayload struct {
	Peer    Peer   `json:"peer"`
	MatchID string `json:"matchId"`
}

var (
	users = make(map[string]Peer)
	mu    sync.Mutex
)

func main() {
	addr := net.UDPAddr{
		Port: HOLE_PUNCH_SERVER_PORT,
		IP:   net.ParseIP(HOST_IP),
	}

	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		log.Fatal("UDP server listening error:", err)
	}
	defer conn.Close()

	log.Println("UDP holepunching server listening on", addr.String())

	buf := make([]byte, 1024)

	for {
		n, remoteAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Println("Failed to read message:", err)
			continue
		}

		go handleMessage(conn, buf[:n], remoteAddr)
	}
}

// Handle recieve a message on the server
func handleMessage(conn *net.UDPConn, data []byte, remote *net.UDPAddr) {
	var msg Message
	log.Println("Server recieved a message")
	if err := json.Unmarshal(data, &msg); err != nil {
		log.Println("Invalid JSON:", err)
		return
	}

	if msg.UID == "" || msg.PeerUID == "" {
		log.Println("Invalid message format")
		return
	}

	mu.Lock()
	defer mu.Unlock()

	// Remove user from map
	if msg.Kill {
		log.Printf("Removing users: %s & %s\n", msg.UID, msg.PeerUID)
		delete(users, msg.UID)
		delete(users, msg.PeerUID)
		return
	}

	// Send message off to user
	users[msg.UID] = Peer{
		UID:     msg.UID,
		Address: remote.IP.String(),
		Port:    remote.Port,
	}

	log.Printf("Stored %s at %s:%d\n", msg.UID, remote.IP, remote.Port)

	if peer, exists := users[msg.PeerUID]; exists {
		sendMatchData(conn, users[msg.UID], peer)
	}
}

func sendMatchData(conn *net.UDPConn, a, b Peer) {
	// Track the delay if needed.
	// start := time.Now()

	matchID := uuid.New().String()

	msgToA, _ := json.Marshal(MatchPayload{
		Peer:    b,
		MatchID: matchID,
	})

	msgToB, _ := json.Marshal(MatchPayload{
		Peer:    a,
		MatchID: matchID,
	})

	clientSend(conn, msgToA, a)
	clientSend(conn, msgToB, b)

	// Read the delay if needed
	//duration := time.Since(start)
	// log.Printf("Match data sent to %s and %s in %s\n", a.UID, b.UID, duration)
}

// Send the required port and address information to both parties
func clientSend(conn *net.UDPConn, msg []byte, peer Peer) {
	addr := net.UDPAddr{
		IP:   net.ParseIP(peer.Address),
		Port: peer.Port,
	}
	_, err := conn.WriteToUDP(msg, &addr)
	if err != nil {
		log.Printf("Error sending information to %s: %v\n", peer.UID, err)
	} else {
		log.Printf("Sent peer information to %s\n", peer.UID)
	}
}
