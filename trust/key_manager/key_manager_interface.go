package key_manager

import "github.com/quantadex/distributed_quanta_bridge/trust/coin"

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
    LoadNodeKeys() error

    /**
     * GetPublicKey
     *
     * Returns the public key stored in object.
     */
    GetPublicKey() (string, error)

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
    SignMessageObj(original interface{}) (*string)

    /**
     * DecodeMessage
     *
     * Uses the provided key to decode the given message. Returns decoded.
     * Note. This does not use the local node's keys but the provided key.
     */
    VerifySignatureObj(original interface{}, key string) bool

    /**
     * SignTX
     * Decodes the transaction envelope, and adds our signature
     */
     SignTransaction(base64 string) (string, error)

    /**
     * VerifyTX
     * Decode the transaction envelope, and check the signature
     */
    VerifyTransaction(base64 string) (bool, error)

    /**
     * DecodeMessage
     * Converts the base64 tx back to Deposit
     */
     DecodeTransaction(base64 string) (*coin.Deposit, error)
}

func NewKeyManager() (KeyManager, error) {
    return &QuantaKeyManager{}, nil
}
