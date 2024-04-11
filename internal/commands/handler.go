package commands

type CommandHandler interface {
	// Runs the request and write response to client
	Handle(ClientRequest)
	// Checks that calling user has the permissions required by the command
	Authorize(ClientRequest) error
	// Returns the list of permissions required to run the command
	Permissions() []string
	// Returns true if the command is going to persist data
	Persistent() bool
}
