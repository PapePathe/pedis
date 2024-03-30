package node

type Config struct {
	Bootstrap      bool
	JoinAddr       string
	MembershipAddr string
	RaftAddr       string
	ServerAddr     string
	ServerId       string
	StartJoinAddrs []string
}
