package server

import (
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/ethanmidgley/the-sequel/in-memory/db"
	"github.com/ethanmidgley/the-sequel/in-memory/handlers"
	"github.com/ethanmidgley/the-sequel/in-memory/pkg/resp"
)

type Server struct {
	wg         sync.WaitGroup
	listener   net.Listener
	shutdown   chan struct{}
	connection chan net.Conn
}

func New(address string) (*Server, error) {

	l, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}

	return &Server{
		listener:   l,
		shutdown:   make(chan struct{}),
		connection: make(chan net.Conn),
	}, nil

}

func (s *Server) AcceptConnections() {
	defer s.wg.Done()

	for {
		select {
		case <-s.shutdown:
			return
		default:
			conn, err := s.listener.Accept()
			if err != nil {
				if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
					continue
				}
				continue
			}
			s.connection <- conn
		}
	}
}

func (s *Server) HandleConnections() {
	defer s.wg.Done()

	for {
		select {
		case <-s.shutdown:
			return
		case conn := <-s.connection:
			go s.HandleConnection(conn)
		}
	}
}

func (s *Server) HandleConnection(conn net.Conn) {
	defer conn.Close()

	for {

		r := resp.New(conn)
		value, err := r.Read()

		if err == io.EOF {
			return
		}

		writer := resp.NewWriter(conn)

		if err != nil {
			writer.Write(resp.Value{Typ: "error", Str: "ERROR PARSING DATA"})
			continue
		}

		if value.Typ != "array" {
			writer.Write(resp.Value{Typ: "error", Str: "INVALID REQUEST, EXPECTED ARRAY"})
			continue
		}

		if len(value.Array) == 0 {
			writer.Write(resp.Value{Typ: "error", Str: "INVALID REQUEST, EXPECTED ARRAY LENGTH > 0"})
			continue
		}

		command := strings.ToUpper(value.Array[0].Bulk)
		args := value.Array[1:]

		handler, ok := handlers.Handlers[command]
		if !ok {
			writer.Write(resp.Value{Typ: "error", Str: "INVALID REQUEST, NO COMMAND"})
			continue
		}

		if command == "SET" {
			db.AOF.Write(value)
		}

		v := handler(args)
		writer.Write(v)

	}

}

func (s *Server) Start() {
	s.wg.Add(2)
	go s.AcceptConnections()
	go s.HandleConnections()
}

func (s *Server) Stop() {
	close(s.shutdown)
	s.listener.Close()

	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return
	case <-time.After(time.Second):
		fmt.Println("Timed out waiting for connections to finish.")
		return
	}

}
