package manifest

import (
	"encoding/json"
	"errors"
	"strconv"
)

/**
 * TrustNode
 *
 * Describes a trust node. The IP and port on which the node is listening and it's public key.
 */
type TrustNode struct {
	IP           string
	Port         string
	ExternalPort string
	PubKey       string
	State        string
	ChainAddress map[string]string
}

/**
 * Manifest
 *
 * The manifest is a struct that describes a full trust of nodes.
 */
type Manifest struct {
	N                int          // Total number of nodes in trust
	Q                int          // Number nodes needed to spend trust
	Nodes            []*TrustNode // The nodes in the trust. The key is the nodeID which is the order they were added.
	ContractAddress  string       // The address of the coin contract
	ContractCallSite string       // The address where nodes post to the contract
	QuantaAddress    string       // The quanta-trust address
}

/**
 * CreateNewManifest
 *
 * Creates a new Manifest struct with N and Q set and an empty
 * trust node list
 */
func CreateNewManifest(quorumNodes int) *Manifest {
	return &Manifest{Q: quorumNodes, Nodes: []*TrustNode{}}
}

/**
 *  CreateNewManifestFromJSON
 *
 *  Takes a JSON ([]byte) array which was created from a fully formed Manifest.
 *  Returns a fully formed manifest.
 *
 */
func CreateManifestFromJSON(data []byte) (*Manifest, error) {
	m := Manifest{}
	err := json.Unmarshal(data, &m)
	return &m, err
}

/**
 * GetJSON
 *
 * Return the JSON ([]byte) representation of this manifest object.
 * This is the inverse of CreateNewManifestFromJSON
 */
func (m *Manifest) GetJSON() ([]byte, error) {
	return json.Marshal(&m)
}

/**
 * AddNode
 *
 * Creates a new TrustNode object and inserts it into the manifest.
 * This should create an error if the manifest is already completed.
 * The nodeID of this new node is the order in which it was added. (e.g very first node is 0, last is N-1)
 *
 */
func (m *Manifest) AddNode(ip string, port string, externalPort string, pubKey string, chainAddress map[string]string) error {
	if m.ManifestComplete() {
		return errors.New("Manifest already completed")
	}

	nodeId, err := m.FindNode(ip, port, pubKey)
	if err != nil {
		m.Nodes = append(m.Nodes, &TrustNode{ip, port, externalPort,pubKey, "ADDED", chainAddress})
		m.N++
		return nil
	}

	return errors.New("Node already exist on " + strconv.Itoa(nodeId))
}

/**
 * Manifest Complete
 *
 * Returns true if the Mnaifest has exactly N nodes. False otherwise.
 */
func (m *Manifest) ManifestComplete() bool {
	//for running two nodes
	return m.N >= m.Q

	//return m.N >= m.Q
}

/**
 * FindNode
 *
 * If a node with the given IP/port exists return it's nodeID. Otherwise return error
 */
func (m *Manifest) FindNode(ip string, port string, pubKey string) (nodeID int, err error) {
	for k, v := range m.Nodes {
		if v.IP == ip && v.Port == port && v.PubKey == pubKey {
			return k, nil
		}
	}
	return 0, errors.New("cannot find node")
}

/**
 * Update state
 */
func (m *Manifest) UpdateState(nodeKey string, state string) error {
	for _, v := range m.Nodes {
		if v.PubKey == nodeKey {
			v.State = state
			return nil
		}
	}

	return errors.New("Node not found for pubkey=" + nodeKey)
}
