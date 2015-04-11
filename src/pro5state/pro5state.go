package pro5state

import (
  "container/list"

  "pro5"
)

type State struct {
  listeners list.List
  displayLayouts string
  lastSlide string
}

func New() *State {
  return new(State)
}

// Listen for new websocket clients (channels that receive a single payload).
func (s *State) ListenForClients() (chan<- (chan<- string)) {
  acceptor := make(chan (chan<- string))
  go func() {
    for listener := range acceptor {
      s.listeners.PushBack(listener)
      sendToListener(s, listener, s.displayLayouts)
      sendToListener(s, listener, s.lastSlide)
    }
  }()
  return acceptor
}

// Listen for new messages from pro5.
func (s *State) ListenForMessages() (chan<- pro5.StageMessage) {
  messages := make(chan pro5.StageMessage)
  go func() {
    for message := range messages {
      switch message.Name {
      case "DisplayLayouts":
        s.displayLayouts = message.Payload
      case "StageDisplayData":
        s.lastSlide = message.Payload
      }
      sendToListeners(s, message.Payload)
    }
  }()
  return messages
}

func sendToListeners(s *State, payload string) (err error) {
  for e := s.listeners.Front(); e != nil; e = e.Next() {
    listener, ok := e.Value.(chan<- string)
    if ok {
      err = sendToListener(s, listener, payload)
      if err != nil {
        return
      }
    }
  }
  return
}

func sendToListener(s *State, listener chan<- string, payload string) error {
  if len(payload) > 0 {
    listener <- payload
  }
  return nil
}
