package key_manager

import "crypto/ecdsa"

/**
 * KeyManager
 *
 * This module maintains the key pair of the node used for signing for the node.
 * This node also performs the signing of messages for the node.
 * It is also used for validating that peer nodes correctly signed messages.
 */
type KeyManager interface {
	/**
	 * CreateNodeKeys
	 *
	 * Generates a new private / public key pair and stashes it in local object.
	 * Also, stashes it in the node key store.
	 */
	CreateNodeKeys() error

	/**
	 * LoadNodeKeys
	 *
	 * Loads the keys from the node key store.
	 */
	LoadNodeKeys(privKey string) error

	/**
	 * GetPublicKey
	 *
	 * Returns the public key stored in object.
	 */
	GetPublicKey() (string, error)

	/**
	 * GetPrivateKey
	 *
	 * Returns the private key as ecdsa.
	 */
	GetPrivateKey() *ecdsa.PrivateKey

	/**
	 * SignMessage
	 *
	 * Uses the private key to sign the given message and returns the signed message
	 */
	SignMessage(original []byte) ([]byte, error)

	/**
	 * SignMessage
	 *
	 * Uses the private key to sign the given message and returns the signed message
	 * returns base64 signature
	 */
	SignMessageObj(original interface{}) *string

	/**
	 * DecodeMessage
	 *
	 * Uses the provided key to decode the given message. Returns decoded.
	 * Note. This does not use the local node's keys but the provided key.
	 */
	VerifySignatureObj(original interface{}, key string) bool

	/**
	 * SignTransaction
	 * Decodes the transaction, and return signature with
	 * our key
	 */
	SignTransaction(encoded string) (string, error)

	/**
	 * VerifyTX
	 * Decode the transaction envelope, and check the signature
	 */
	VerifyTransaction(encoded string) (bool, error)
}

func NewKeyManager(network string) (KeyManager, error) {
	return &QuantaKeyManager{network: network}, nil
}

func NewEthKeyManager() (KeyManager, error) {
	return &EthereumKeyManager{}, nil
}

func NewGrapheneKeyManager(chain string) (KeyManager, error) {
	return &QuantaKeyGraphene{chain: chain}, nil
}
func NewBitCoinKeyManager() (KeyManager, error) {
	return &BitcoinKeyManager{}, nil
}
