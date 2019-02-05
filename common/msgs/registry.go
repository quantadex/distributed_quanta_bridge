package msgs

type NodeInfo struct {
	NodeIp       string            `json: "node_ip"`
	NodePort     string            `json: "node_port"`
	NodeKey      string            `json: "node_key"`
	ChainAddress map[string]string `json: "chain_address"`
}

type RegisterReq struct {
	Body      NodeInfo `json: "body"`
	Signature string   `json: "signature"`
}

type PingBody struct {
	NodeKey string `json: "node_key"`
	Status  string `json: "status"`
}

type PingReq struct {
	Body      PingBody `json: "body"`
	Signature string   `json: "signature"`
}
