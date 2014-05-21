// subscription package implement a subscription service.
// Client connect to the servive and should send a unique ID as a
// 'CRLF' terminated string.
// Each subscription request is then routed back to a receiver listening
// on SubscriptionChan, which receive the SubscriberID and it's tcp connection.
package subscription

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

// SubscriptionRequest associate a client identified by an ID
// with it's tcp connection
type SubscriptionRequest struct {
	SubscriberID string
	Conn         *net.TCPConn
}

// Subscription server listen for client connection
// and broadcast them through SubscriptionChan has SubscriptionRequest
type SubscriptionServer struct {
	Port             int
	SubscriptionChan chan *SubscriptionRequest
}

func (s *SubscriptionServer) handleSubscriptionRequest(conn *net.TCPConn) {
	ID, err := bufio.NewReader(conn).ReadString('\n')
	ID = ID[:len(ID)-1]

	if err != nil {
		log.Printf("Cannot read payload: %v", err)
		return
	}

	s.SubscriptionChan <- &SubscriptionRequest{ID, conn}
}

func (s *SubscriptionServer) Listen() {
	addr, _ := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%v", s.Port))

	ln, err := net.ListenTCP("tcp", addr)

	defer ln.Close()

	if err != nil {
		log.Fatalf("*** Cannot start Subscription server: %v", err)
	}

	log.Printf("=== Subscription server listening to %v", s.Port)

	for {
		conn, err := ln.AcceptTCP()

		if err != nil {
			log.Printf("Cannot read from socket: %v", err)
			continue
		}

		go s.handleSubscriptionRequest(conn)
	}
}

func New(port int, subscriptionChan chan *SubscriptionRequest) *SubscriptionServer {
	return &SubscriptionServer{
		Port:             port,
		SubscriptionChan: subscriptionChan,
	}
}
