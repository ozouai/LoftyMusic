package audiosourceserver

import (
	"context"
	"fmt"
	"sync"

	"github.com/ozouai/loftymusic/audiosource/audiosourcepb"
)

type router struct {
	audioRouter map[string][]*subRouter
}

type subRouter struct {
	slots     []chan *audiosourcepb.ControlChannelResponse
	freeSlots chan int64
	wg        sync.WaitGroup
	server    audiosourcepb.AudioSource_ControlChannelServer
	slotLock  sync.RWMutex
}

func (m *router) Add(audioType string, server audiosourcepb.AudioSource_ControlChannelServer) {
	sub := &subRouter{
		slots:     make([]chan *audiosourcepb.ControlChannelResponse, 32),
		freeSlots: make(chan int64, 32+1),
	}
	m.audioRouter[audioType] = append(m.audioRouter[audioType], sub)
	sub.wg.Add(1)
	go sub.loop()
	sub.wg.Wait()
}

func (m *subRouter) loop() {
	defer m.wg.Done()
	for {
		msg, err := m.server.Recv()
		if err != nil {
			panic(err)
		}
		m.processMsg(context.TODO(), msg)
	}
}

func (m *subRouter) processMsg(ctx context.Context, msg *audiosourcepb.ControlChannelResponse) {
	m.slotLock.RLock()
	defer m.slotLock.RUnlock()
	m.slots[msg.GetId()] <- msg
}

func (m *router) send(ctx context.Context, audioType string, req *audiosourcepb.ControlChannelRequest) (*audiosourcepb.ControlChannelResponse, error) {
	if len(m.audioRouter[audioType]) == 0 {
		return nil, fmt.Errorf("error no audio sources for type")
	}
	res, err := m.audioRouter[audioType][0].send(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	return res, nil
}

func (m *subRouter) send(ctx context.Context, req *audiosourcepb.ControlChannelRequest) (*audiosourcepb.ControlChannelResponse, error) {
	slot := <-m.freeSlots
	defer func() {
		m.freeSlots <- slot
	}()
	req.Id = slot
	responseChan := make(chan *audiosourcepb.ControlChannelResponse, 1)
	m.slots[slot] = responseChan
	err := m.server.Send(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request to audio source: %w", err)
	}
	resp := <-responseChan
	close(responseChan)
	return resp, nil
}
