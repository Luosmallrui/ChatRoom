package client

import "chatroom/pkg/email"

type Providers struct {
	EmailClient *email.Client
}
