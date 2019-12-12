package tss

import (
	"encoding/json"
	"fmt"
	"github.com/keep-network/keep-tecdsa/pkg/net"
)

// MessageRouting holds the information required to route a message. It determines
// a type of routing based on the receiver. If receiver is `nil` it will assume
// to broadcast the message. If receiver is provided it will send direct unicast
// message to the receiver. It holds the message itself as well.
type MessageRouting struct {
	ReceiverID []byte
	Message    net.TaggedMarshaler
}

// TSSMessage is a network message used to transport messages in TSS protocol
// execution.
type TSSMessage struct {
	SenderID    []byte
	Payload     []byte
	IsBroadcast bool
}

// Type returns a string type of the `TSSMessage` so that it conforms to
// `net.Message` interface.
func (m *TSSMessage) Type() string {
	return fmt.Sprintf("%T", m)
}

// Marshal converts this message to a byte array suitable for network communication.
func (m *TSSMessage) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

// Unmarshal converts a byte array produced by Marshal to a message.
func (m *TSSMessage) Unmarshal(bytes []byte) error {
	var message TSSMessage
	if err := json.Unmarshal(bytes, &message); err != nil {
		return err
	}

	m.Payload = message.Payload
	m.SenderID = message.SenderID
	m.IsBroadcast = message.IsBroadcast

	return nil
}
