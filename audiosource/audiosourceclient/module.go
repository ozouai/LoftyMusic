package audiosourceclient

import (
	"context"
	"fmt"
	"sync"

	"github.com/ozouai/loftymusic/audiosource/audiosourcepb"
	"google.golang.org/grpc"
)

type AudioSourceClient struct {
	handler   AudioSourceHandler
	channel   audiosourcepb.AudioSource_ControlChannelClient
	wg        sync.WaitGroup
	audioType string
}

type AudioSourceHandler interface {
	Search(ctx context.Context, req *audiosourcepb.Search_Request) (*audiosourcepb.Search_Response, error)
}

func New(ctx context.Context, serverAddr string, audioType string, handler AudioSourceHandler) (*AudioSourceClient, error) {
	m := &AudioSourceClient{
		audioType: audioType,
	}
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("error error dialing server: %w", err)
	}
	client := audiosourcepb.NewAudioSourceClient(conn)
	channel, err := client.ControlChannel(ctx)
	if err != nil {
		return nil, fmt.Errorf("error opening control channel: %w", err)
	}
	m.channel = channel
	m.wg.Add(1)
	go m.loop(ctx)
	return m, nil
}

func (m *AudioSourceClient) Wait() {
	m.wg.Wait()
}

func (m *AudioSourceClient) loop(ctx context.Context) {
	defer m.wg.Done()
	for {
		msg, err := m.channel.Recv()
		if err != nil {
			continue
		}
		m.handleControlChannelMsg(ctx, msg)
	}
}

func (m *AudioSourceClient) handleControlChannelMsg(ctx context.Context, msg *audiosourcepb.ControlChannelRequest) {
	resp := &audiosourcepb.ControlChannelResponse{
		Id: msg.GetId(),
	}
	switch req := msg.GetRequest().(type) {
	case *audiosourcepb.ControlChannelRequest_Register:
		resp.Response = &audiosourcepb.ControlChannelResponse_Register{Register: &audiosourcepb.Register_Response{SourceType: m.audioType}}
	case *audiosourcepb.ControlChannelRequest_Search:
		res, err := m.handler.Search(ctx, req.Search)
		if err != nil {
			panic(err)
		}
		resp.Response = &audiosourcepb.ControlChannelResponse_Search{Search: res}
	}
	m.channel.Send(resp)
}
