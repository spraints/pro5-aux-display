package pro5state

import (
  "container/list"
)

type State struct {
  Listeners list.List
  DisplayLayouts string
  LastSlide string
}

func New() *State {
  return new(State)
}

// Listen for new websocket clients (channels that receive a single payload).
func (s *State) ListenForClients() (chan<- (chan<- string)) {
  acceptor := make(chan (chan<- string))
  go func() {
    for listener := range acceptor {
      s.Listeners.PushBack(listener)
      sendToListener(s, listener, s.DisplayLayouts)
      sendToListener(s, listener, s.LastSlide)
    }
  }()
  return acceptor
}

// Implement pro5stage.Client
func (s *State) SendMessage(name string, payload string) {
  switch name {
  case "DisplayLayouts":
    s.DisplayLayouts = payload
  case "StageDisplayData":
    s.LastSlide = payload
  }
  sendToListeners(s, payload)
}

func sendToListeners(s *State, payload string) (err error) {
  for e := s.Listeners.Front(); e != nil; e = e.Next() {
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
