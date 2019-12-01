package http

import (
	"net"
	"net/http"
)

// Server for http
type Server struct {
	port string
	http *http.Server
	Mux  *http.ServeMux
}

// Run this http server
func (s *Server) Run() error {
	ln, err := net.Listen("tcp", "127.0.0.1:"+s.port)
	if err != nil {
		return err
	}

	if err := s.http.Serve(ln); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

// Close this http server
func (s *Server) Close() error {
	return s.http.Close()
}

// NewServer for http
func NewServer(port string) (*Server, error) {
	m := http.NewServeMux()
	s := &Server{
		port: port,
		http: &http.Server{
			Handler: m,
		},
		Mux: m,
	}

	return s, nil
}
