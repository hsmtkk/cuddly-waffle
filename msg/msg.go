package msg

import (
	"encoding/json"
	"fmt"
)

type Message struct {
	ID    int64 `json:"id"`
	Count int64 `json:"count"`
}

func NewMessage(id, count int64) Message {
	return Message{id, count}
}

func (m Message) ToJSON() ([]byte, error) {
	bs, err := json.Marshal(&m)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON; %w", err)
	}
	return bs, nil
}

func FromJSON(bs []byte) (Message, error) {
	m := Message{}
	if err := json.Unmarshal(bs, &m); err != nil {
		return m, fmt.Errorf("failed to unmarshal JSON; %w", err)
	}
	return m, nil
}
