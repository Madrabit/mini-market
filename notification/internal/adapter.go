package internal

import "log"

type StubSender struct {
}

func (s *StubSender) Send(to, message string) error {
	log.Printf("[stub notify] to=%s, message=%s", to, message)
	return nil
}
