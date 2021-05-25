package magmasc

import (
	"encoding/json"

	"0chain.net/core/datastore"
)

// Consumer represents consumer stored in block chain.
type Consumer struct {
	ID      string `json:"id"`
	BaseURL string `json:"url"`
}

// GetKey returns concatenated smart contract id and Consumer.ID.
// Should be used while inserting Consumer in state.StateContextI.
func (c *Consumer) GetKey(scKey string) datastore.Key {
	return scKey + c.ID
}

// Encode implements util.Serializable interface.
func (c *Consumer) Encode() []byte {
	buff, _ := json.Marshal(c)
	return buff
}

// Decode implements util.Serializable interface.
func (c *Consumer) Decode(input []byte) error {
	err := json.Unmarshal(input, c)
	if err != nil {
		return err
	}
	return nil
}

// Consumers represent sorted in alphabetic order by Consumer.ID consumers.
type Consumers struct {
	Nodes sortedConsumers
}

// Decode implements util.Serializable interface.
func (cs *Consumers) Decode(input []byte) error {
	err := json.Unmarshal(input, cs)
	if err != nil {
		return err
	}
	return nil
}

// Encode implements util.Serializable interface.
func (cs *Consumers) Encode() []byte {
	buff, _ := json.Marshal(cs)
	return buff
}
