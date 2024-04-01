package discovery

import (
	"fmt"
	"log"
	"log/slog"
	"net"

	"github.com/hashicorp/serf/serf"
)

type Membership struct {
	Config
	handler Handler
	serf    *serf.Serf
	events  chan serf.Event
}

func NewMembership(handler Handler, config Config) (*Membership, error) {
	slog.Info("Creating new Member in cluster")
	c := &Membership{
		Config:  config,
		handler: handler,
	}

	if err := c.setupSerf(); err != nil {
		return nil, err
	}

	return c, nil
}

func (m *Membership) setupSerf() (err error) {
	addr, err := net.ResolveTCPAddr("tcp", m.BindAddr)

	if err != nil {
		return fmt.Errorf("error resolving tcp addr %v", err)
	}
	config := serf.DefaultConfig()
	config.Init()
	config.MemberlistConfig.BindAddr = addr.IP.String()
	config.MemberlistConfig.BindPort = addr.Port
	m.events = make(chan serf.Event)
	config.EventCh = m.events
	config.Tags = m.Tags
	config.NodeName = m.Config.NodeName
	m.serf, err = serf.Create(config)

	if err != nil {
		return fmt.Errorf("error creating serf node %v", err)
	}
	go m.eventHandler()
	log.Println("serf start join addrs", m.StartJoinAddrs)
	if m.StartJoinAddrs != nil {
		_, err = m.serf.Join(m.StartJoinAddrs, true)

		if err != nil {
			return fmt.Errorf("error joining serf node %v", err)
		}
	}

	return nil
}

func (m *Membership) eventHandler() {
	slog.Info("Started handling events")
	for e := range m.events {
		switch e.EventType() {

		case serf.EventMemberJoin:
			for _, member := range e.(serf.MemberEvent).Members {
				if m.isLocal(member) {
					continue
				}
				m.handleJoin(member)
			}

		case serf.EventMemberLeave:
			for _, member := range e.(serf.MemberEvent).Members {
				if m.isLocal(member) {
					return
				}
				m.handleLeave(member)
			}
		case serf.EventMemberFailed:
		}
	}
}

func (m *Membership) handleJoin(member serf.Member) {
	if err := m.handler.Join(member.Name, m.Tags["rpc_addr"]); err != nil {
		slog.Error(err.Error(), "failed to join", member)
	}
}

func (m *Membership) handleLeave(member serf.Member) {
	if err := m.handler.Leave(member.Name); err != nil {
		slog.Error(err.Error(), "failed to leave", member)
	}
}

func (m *Membership) isLocal(member serf.Member) bool {
	return m.serf.LocalMember().Name == member.Name
}

func (m *Membership) Members() []serf.Member {
	return m.serf.Members()
}

func (m *Membership) Leave() error {
	return m.serf.Leave()
}
