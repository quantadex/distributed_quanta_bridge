package peer_contact

type PeerMsgRequest struct {
	Body        PeerMessage	`json: "body"`
	Signature   string      `json: "signature"`
}