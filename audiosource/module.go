package audiosource

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
)

type AudioSource struct {
}

func New() (*AudioSource, error) {
	m := &AudioSource{}

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		return nil, fmt.Errorf("error creating audiosource listener: %w", err)
	}
	server := grpc.NewServer()
	// audiosourcepb.RegisterAudioSourceServer()
	server.Serve(listener)

	return m, nil
}
