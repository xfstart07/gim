// Author: xufei
// Date: 2019-09-29

package command

import "gim/internal/ciface"

type shutdownCommand struct {
	userClient ciface.UserClient
}

func (s *shutdownCommand) Process(msg string) {
	s.userClient.Shutdown()
}

func NewShutDownCommand(userClient ciface.UserClient) *shutdownCommand {
	return &shutdownCommand{
		userClient: userClient,
	}
}
