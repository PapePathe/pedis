package commands

type CommandHandler interface {
	// Runs the request and write response to client
	Handle(IClientRequest)
	// Checks that calling user has the permissions required by the command
	Authorize(IClientRequest) error
	// Returns the list of permissions required to run the command
	Permissions(IClientRequest) []string
	// Returns true if the command is going to persist data
	Persistent(IClientRequest) bool
}
