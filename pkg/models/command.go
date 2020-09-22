package models

// Command is an operatiion to Finite State Machine
type Command struct {
	Operation string
	Key       []byte
	Value     []byte
}
