package magmasc

import (
	"encoding/json"
	"strings"

	"0chain.net/chaincore/chain/state"
	"0chain.net/core/datastore"
)

type (
	// Consumer represents consumers node stored in block chain.
	Consumer struct {
		ID string `json:"id"`
	}

	// Provider represents providers node stored in block chain.
	Provider struct {
		ID    string `json:"id"`
		Terms Terms  `json:"terms"`
	}

	// Terms represents information of Provider's services.
	Terms struct {
		Price int64 `json:"price"` // per MB
		QoS   QoS   `json:"qos"`
	}

	// QoS represents a Quality of Service and contains uploading and downloading speed represented in megabytes per second.
	QoS struct {
		DownloadMBPS int64 `json:"download_mbps"`
		UploadMBPS   int64 `json:"upload_mbps"`
	}
)

const (
	// providerType is a type of Provider's node.
	providerType = "provider"

	// consumerType is a type of Consumer's node.
	consumerType = "consumer"
)

// nodeKey returns a specific key for Node interacting with magma smart contract.
// scKey is an ID of magma smart contract and nodeID is and ID of Node.
//
// Should be used while inserting, removing or getting Node in state.StateContextI
func nodeKey(scKey, nodeID, nodeType string) datastore.Key {
	return strings.Join([]string{scKey, nodeType, nodeID}, ":")
}

// GetID returns Consumer.ID.
func (c *Consumer) GetID() string {
	return c.ID
}

// GetType returns Consumer's type.
func (c *Consumer) GetType() string {
	return consumerType
}

// Encode implements util.Serializable interface.
func (c *Consumer) Encode() []byte {
	buff, _ := json.Marshal(c)
	return buff
}

// Decode implements util.Serializable interface.
func (c *Consumer) Decode(input []byte) error {
	return json.Unmarshal(input, c)
}

// GetID returns Provider.ID.
func (p *Provider) GetID() string {
	return p.ID
}

// GetType returns Provider's type.
func (p *Provider) GetType() string {
	return providerType
}

// Encode implements util.Serializable interface.
func (p *Provider) Encode() []byte {
	buff, _ := json.Marshal(p)
	return buff
}

// Decode implements util.Serializable interface.
func (p *Provider) Decode(input []byte) error {
	return json.Unmarshal(input, p)
}

// containsNode looks for provided Consumer in provided Consumers and state.StateContextI.
// If Consumer will be found it returns true, else false.
func containsConsumer(scKey string, consumer Consumer, consumers *Consumers, balances state.StateContextI) bool {
	for _, c := range consumers.Nodes {
		if c.GetID() == consumer.GetID() {
			return true
		}
	}

	_, err := balances.GetTrieNode(nodeKey(scKey, consumer.GetID(), consumer.GetType()))
	if err == nil {
		return true
	}

	return false
}

// containsNode looks for provided Provider in provided Providers and state.StateContextI.
// If Provider will be found it returns true, else false.
func containsProvider(scKey string, provider Provider, providers *Providers, balances state.StateContextI) bool {
	for _, c := range providers.Nodes {
		if c.GetID() == provider.GetID() {
			return true
		}
	}

	_, err := balances.GetTrieNode(nodeKey(scKey, provider.GetID(), provider.GetType()))
	if err == nil {
		return true
	}

	return false
}

type (
	// Consumers represents sorted Consumer nodes, used to inserting, removing or getting from
	// state.StateContextI with AllConsumersKey.
	Consumers struct {
		Nodes sortedConsumers
	}

	// Providers represents sorted Provider nodes, used to inserting, removing or getting from
	// state.StateContextI with AllProvidersKey.
	Providers struct {
		Nodes sortedProviders
	}
)

// Encode implements util.Serializable interface.
func (s *Consumers) Encode() []byte {
	blob, _ := json.Marshal(s)
	return blob
}

// Decode implements util.Serializable interface.
func (s *Consumers) Decode(blob []byte) error {
	return json.Unmarshal(blob, s)
}

// Encode implements util.Serializable interface.
func (s *Providers) Encode() []byte {
	blob, _ := json.Marshal(s)
	return blob
}

// Decode implements util.Serializable interface.
func (s *Providers) Decode(blob []byte) error {
	return json.Unmarshal(blob, s)
}

type (
	// Acknowledgment contains the necessary data obtained when the consumer accepts the provider Terms.
	//
	// Acknowledgment stores in the state of the blockchain as a result of performing the acceptTerms
	// MagmaSmartContract function.
	Acknowledgment struct {
		ID            string `json:"id"`
		ConsumerID    string `json:"consumer_id"`
		ProviderID    string `json:"provider_id"`
		AccessPointID string `json:"access_point_id"`
		SessionID     string `json:"session_id"`
	}
)

// acknowledgmentKey returns a specific key for Acknowledgment.
// scKey is an ID of magma smart contract and id is and ID of Acknowledgment.
//
// Should be used while inserting, removing or getting Node in state.StateContextI
func acknowledgmentKey(ssKey, id string) datastore.Key {
	return ssKey + id
}

// isValid checks Acknowledgment.ProviderID, Acknowledgment.AccessPointID, Acknowledgment.SessionID for emptiness.
// If any field is empty isValid return false, else true.
func (a *Acknowledgment) isValid() bool {
	return a.ProviderID != "" && a.AccessPointID != "" && a.SessionID != ""
}

// Decode implements util.Serializable interface.
func (a *Acknowledgment) Decode(input []byte) error {
	err := json.Unmarshal(input, a)
	if err != nil {
		return err
	}
	return nil
}

// Encode implements util.Serializable interface.
func (a *Acknowledgment) Encode() []byte {
	buff, _ := json.Marshal(a)
	return buff
}
