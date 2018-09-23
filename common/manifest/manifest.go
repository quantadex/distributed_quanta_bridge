package manifest

/**
 * TrustNode
 *
 * Describes a trust node. The IP and port on which the node is listening and it's public key.
 */
type TrustNode struct {
    IP string
    Port string
    PubKey string
}

/**
 * Manifest
 *
 * The manifest is a struct that describes a full trust of nodes.
 */
type Manifest struct {
    N int // Total number of nodes in trust
    Q int // Number nodes needed to spend trust
    Nodes map[int]TrustNode // The nodes in the trust. The key is the nodeID which is the order they were added.
    ContractAddress string // The address of the coin contract
    ContractCallSite string // The address where nodes post to the contract
    QuantaAddress string // The quanta-trust address
}

/**
 * CreateNewManifest
 *
 * Creates a new Manifest struct with N and Q set and an empty
 * trust node list
 */
func CreateNewManifest(totalNodes string, quorumNodes string) (*Manifest, error) {
    return nil, nil
}

/**
 *  CreateNewManifestFromJSON
 *
 *  Takes a JSON ([]byte) array which was created from a fully formed Manifest.
 *  Returns a fully formed manifest.
 *
 */
func CreateManifestFromJSON(data []byte) (*Manifest, error) {
    return nil, nil
}

/**
 * GetJSON
 *
 * Return the JSON ([]byte) representation of this manifest object.
 * This is the inverse of CreateNewManifestFromJSON
 */
func (m *Manifest) GetJSON() ([]byte, error) {
    return nil, nil
}

/**
 * AddNode
 *
 * Creates a new TrustNode object and inserts it into the manifest.
 * This should create an error if the manifest is already completed.
 * The nodeID of this new node is the order in which it was added. (e.g very first node is 0, last is N-1)
 * 
 */
func (m *Manifest) AddNode(ip string, port string, pubKey string) error {
    return nil
}

/**
 * Manifest Complete
 *
 * Returns true if the Mnaifest has exactly N nodes. False otherwise.
 */
func (m *Manifest) ManifestComplete() bool {
    return true
}

/**
 * FindNode
 *
 * If a node with the given IP/port exists return it's nodeID. Otherwise return error
 */
func (m *Manifest) FindNode(ip string, port string) (nodeID int, err error) {
    return 0, nil
}
