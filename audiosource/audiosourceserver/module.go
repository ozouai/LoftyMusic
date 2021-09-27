package audiosourceserver

import (
	"fmt"
	"net"

	"github.com/ozouai/loftymusic/audiosource/audiosourcepb"
	"google.golang.org/grpc"
)

type AudioSourceServer struct {
	audiosourcepb.UnimplementedAudioSourceServer
	router *router
}

func New() (*AudioSourceServer, error) {
	m := &AudioSourceServer{
		router: &router{
			audioRouter: map[string][]*subRouter{},
		},
	}

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		return nil, fmt.Errorf("error creating audiosource server listener: %w", err)
	}
	server := grpc.NewServer()
	audiosourcepb.RegisterAudioSourceServer(server, m)
	server.Serve(listener)

	return m, nil
}

func (m *AudioSourceServer) ControlChannel(server audiosourcepb.AudioSource_ControlChannelServer) error {
	err := server.Send(&audiosourcepb.ControlChannelRequest{
		Request: &audiosourcepb.ControlChannelRequest_Register{Register: &audiosourcepb.Register_Request{Uuid: ""}},
	})
	if err != nil {
		return fmt.Errorf("error sending register request: %w", err)
	}
	resp, err := server.Recv()
	if err != nil {
		return fmt.Errorf("error receiving register request response: %w", err)
	}
	var audioType string
	if registerResp, ok := resp.GetResponse().(*audiosourcepb.ControlChannelResponse_Register); !ok {
		return fmt.Errorf("error received the wrong type of reply for the register command")
	} else {
		audioType = registerResp.Register.GetSourceType()
	}
	m.router.Add(audioType, server)
	return nil
}
