package global

// finite state machine
type robotState struct {
	GroupSpace map[string]string
}

var Rob = robotState{
	GroupSpace: make(map[string]string),
}