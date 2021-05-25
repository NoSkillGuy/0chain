package magmasc

import (
	"encoding/json"

	"0chain.net/chaincore/chain/state"
	"0chain.net/core/datastore"
	"0chain.net/core/util"
)

type (
	// Node represents interface for nodes types with magma smart contract may interact.
	Node interface {
		// GetID returns ID of Node.
		GetID() string

		// Serializable is an embedded interface.
		util.Serializable
	}

	// Consumer represents consumers node stored in block chain.
	Consumer struct {
		ID string `json:"id"`
	}

	// Provider represents providers node stored in block chain.
	Provider struct {
		ID    string `json:"id"`
		Terms Terms  `json:"terms"`
	}

	Terms struct {
		Price int64 `json:"price"`
		QoS   QoS   `json:"qos"`
	}

	QoS struct {
		DownloadMBPS int64 `json:"download_mbps"`
		UploadMBPS   int64 `json:"upload_mbps"`
	}
)

var (
	// Ensure Consumer implements Node.
	_ Node = (*Consumer)(nil)

	// Ensure Provider implements Node.
	_ Node = (*Provider)(nil)
)

// nodeKey returns a specific key for Node interacting with magma smart contract.
// scKey is an ID of magma smart contract and nodeID is and ID of Node.
//
// Should be used while inserting, removing or getting Node in state.StateContextI
func nodeKey(scKey, nodeID string) datastore.Key {
	return scKey + nodeID
}

// GetID returns Consumer.ID.
func (c *Consumer) GetID() string {
	return c.ID
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

// GetID returns Provider.ID.
func (p *Provider) GetID() string {
	return p.ID
}

// Encode implements util.Serializable interface.
func (p *Provider) Encode() []byte {
	buff, _ := json.Marshal(p)
	return buff
}

// Decode implements util.Serializable interface.
func (p *Provider) Decode(input []byte) error {
	err := json.Unmarshal(input, p)
	if err != nil {
		return err
	}
	return nil
}

// Nodes represent sorted in alphabetic order by ID nodes.
type Nodes struct {
	Nodes sortedNodes
}

// Decode implements util.Serializable interface.
func (cs *Nodes) Decode(input []byte) error {
	err := json.Unmarshal(input, cs)
	if err != nil {
		return err
	}
	return nil
}

// Encode implements util.Serializable interface.
func (cs *Nodes) Encode() []byte {
	buff, _ := json.Marshal(cs)
	return buff
}

// containsNode looks for provided Node in provided Nodes and state.StateContextI.
// If Node will be found it returns true, else false.
func containsNode(scKey string, node Node, consumers *Nodes, balances state.StateContextI) bool {
	for _, c := range consumers.Nodes {
		if c.GetID() == node.GetID() {
			return true
		}
	}

	_, err := balances.GetTrieNode(nodeKey(scKey, node.GetID()))
	if err == nil {
		return true
	}

	return false
}
