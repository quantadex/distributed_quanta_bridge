package consensus

import (
	"github.com/op/go-logging"
	"github.com/quantadex/quanta_book/consensus"
	"time"
	"github.com/bluele/gcache"
	"math"
)

var logger = logging.MustGetLogger("cosi")

type Phase int

const (
	phasePrepare Phase = iota
	phaseCommit
)

/**
 * CosiMessage is the message used for consensus
 * retains original msg for tracking purpose
 * string is used for base64, should work for general use case
 */
type CosiMessage struct {
	Msg    string // base64 transaction without signature
	Signed string // base64 transaction with signature
	Phase  Phase
	Initial bool
}

type FinalSignature struct {
	Msg []string
}

type Cosi struct {
	node    consensus.Node
	leader  bool
	timeout time.Duration
	gc      gcache.Cache

	// public
	CosiMsgChan  chan *CosiMessage
	FinalSigChan chan FinalSignature
	Verify       func(msg string) error
	SignMsg      func(msg string) (string, error)
	Persist      func(msg string, repair bool) error
}

func NewProtocol(node consensus.Node, leader bool, timeout time.Duration) *Cosi {
	cache := gcache.New(20).
		LRU().
		Build()

	return &Cosi{node: node,
		CosiMsgChan:  make(chan *CosiMessage, 1),
		FinalSigChan: make(chan FinalSignature),
		leader:       leader, timeout: timeout,
		gc:	cache,}
}

func (c *Cosi) Node() consensus.Node {
	return c.node
}

func (c *Cosi) IsLeader() bool {
	return c.leader
}

func (c *Cosi) StartNewRound(msg string) {
	logger.Infof("Start new round leader=%v nodes=%d threshold=%d", c.leader, len(c.node.Members()), Threshold(len(c.node.Members())))

	// broadcast prep
	if !c.leader {
		logger.Error("Only leader can propose signing")
		return
	}

	go func(msg string) {

		// tell all the nodes to validate my message
		req := &CosiMessage{msg, "", phasePrepare, true}
		c.node.BroadcastMessage("Protocol.Cosi", req)

		commitCount := 1 // count ourself as one
		totalNodes := len(c.node.Members())

		// we expect replies if they are willing to sign the message
		running := true
		for commitCount < totalNodes && running {
			select {
			case msgReceived := <-c.CosiMsgChan:
				if msgReceived.Phase == phasePrepare && msgReceived.Msg == msg {
					commitCount += 1
					logger.Infof("Got commitment %d/%d", commitCount, totalNodes)
				}

			case <-time.After(c.timeout):
				logger.Errorf("Commitment - Timeout should not happen received %d/%d", commitCount, totalNodes)
				running = false
			}
		}

		if commitCount < Threshold(totalNodes) {
			logger.Errorf("Bailing out of leader prep round, not enough nodes")
			c.FinalSigChan <- FinalSignature{nil}
			return
		}

		logger.Infof("Got total of %d commitments, moving forward", commitCount)

		// we got enough messages, let's move to the commit phase
		req = &CosiMessage{msg, "", phaseCommit, true}
		c.node.BroadcastMessage("Protocol.Cosi", req)

		// we expect replies with the signatures
		signCount := 1
		finalSigs := make([]string, 0)

		// we sign it first
		signed, err := c.SignMsg(msg)
		if err != nil {
			logger.Error("Error signing message: ", err.Error())
			c.FinalSigChan <- FinalSignature{nil}
			return
		}
		finalSigs = append(finalSigs, signed)

		running = true
		for signCount < totalNodes && running {
			select {
			case msgReceived := <-c.CosiMsgChan:
				if msgReceived.Phase == phaseCommit && msgReceived.Msg == msg {
					signCount += 1
					finalSigs = append(finalSigs, msgReceived.Signed)
					logger.Infof("Got signature %d/%d", signCount, totalNodes)
				}

			case <-time.After(c.timeout):
				logger.Errorf("Timeout should not happen received sig %d/%d", signCount, totalNodes)
				running = false
			}
		}

		if signCount < Threshold(totalNodes) {
			logger.Error("Bailing out of leader commit round, not enough nodes")
			c.FinalSigChan <- FinalSignature{nil}
			return
		}

		err = c.Persist(msg, false)
		if err != nil {
			logger.Error("Error persisting message: ", err.Error())
		}
		c.FinalSigChan <- FinalSignature{finalSigs}
	}(msg)
}

func (c *Cosi) Start() {
	if !c.leader {
		go c.dispatchFollower()
	}
}

type ConsensusState struct {
	message string
	prepareCount int
	commitCount int
	phase Phase
	persisted bool
	verified bool
}

func (c *Cosi) dispatchFollower() {
	for {
		select {
		case msg := <-c.CosiMsgChan:
			totalNodes := len(c.node.Members())
			var stateObj ConsensusState

			// if someone ask me to prepare, validate, and move to commit
			if msg.Phase == phasePrepare {
				state, err := c.gc.Get(msg.Msg)

				if err != nil || msg.Initial {
					logger.Infof("Initialize cosi (follower) state.")
					state = ConsensusState{
						message: msg.Msg,
						prepareCount: 1, // count 1 from leader
						commitCount: 0, // count 1 from commit
						phase: phasePrepare,
						persisted: false,
					}
				}
				stateObj = state.(ConsensusState)

				// if we saw it before, just count it
				stateObj.prepareCount += 1

				// if initial message, keep counting
				if err == nil && !msg.Initial {
					// my peer also validated it
					c.gc.Set(msg.Msg, stateObj)
					logger.Infof("%s Follower got prepare %d/%d", c.Node().Address(), stateObj.prepareCount, totalNodes)
					continue
				}

				// verify if this is first time we see it
				err = c.Verify(msg.Msg)
				stateObj.verified = err == nil
				c.gc.Set(msg.Msg, stateObj)

				if err != nil {
					logger.Error("%s Unable to validate transactions: %v", c.Node().Address(), err.Error())
					continue
				}

				msg2 := CosiMessage{ Msg: msg.Msg, Signed: msg.Signed, Phase: msg.Phase, Initial: false}
				c.node.BroadcastMessage("Protocol.Cosi", &msg2)
			}

			// the leader is asking me to commit, only sign if enough of my peers will do it
			if msg.Phase == phaseCommit {
				state, err := c.gc.Get(msg.Msg)

				if err != nil {
					logger.Error("Could not find message in cache")
					continue
				}

				stateObj = state.(ConsensusState)

				//logger.Infof("%s Total prepare count %d", c.Node().Address(), stateObj.commitCount)
				if stateObj.prepareCount < Threshold(totalNodes) {
					logger.Errorf("Bailing out of follower commit round, not enough nodes %d/%d", stateObj.prepareCount, totalNodes)
					continue
				}

				stateObj.commitCount += 1
				logger.Infof("%s Follower got phaseCommit %d/%d verified=%v", c.Node().Address(), stateObj.commitCount, totalNodes, stateObj.verified)

				// we have enough commitments from peer, so let's sign it ourself
				if stateObj.phase == phasePrepare {
					stateObj.phase = phaseCommit
					if stateObj.verified {
						stateObj.commitCount += 1
						signed, err := c.SignMsg(msg.Msg)
						if err != nil {
							logger.Error("Error signing message: ", err.Error())
						}
						msg.Signed = signed
						c.node.BroadcastMessage("Protocol.Cosi", msg)
					}
				}
				c.gc.Set(msg.Msg, stateObj)
			}

			logger.Infof("%s Total commit count %d", c.Node().Address(), stateObj.commitCount)
			// we have enough commitment from peers, let's save flag in our database
			if stateObj.commitCount >= Threshold(totalNodes) {
				if stateObj.persisted == false {
					err := c.Persist(msg.Msg, !stateObj.verified)
					if err != nil {
						logger.Error("Error persisting message", err.Error())
					}
					stateObj.persisted = true
					c.gc.Set(msg.Msg, stateObj)
				}
			}
		}
	}
}

// FaultThreshold computes the number of faults that byzcoinx tolerates.
func FaultThreshold(n int) int {
	return int(math.Ceil(float64(n - 1) / 3))
}

// Threshold computes the number of nodes needed for successful operation.
func Threshold(n int) int {
	return n - FaultThreshold(n)
}
