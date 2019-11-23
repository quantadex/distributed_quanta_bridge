package consensus

import (
	"fmt"
	"github.com/go-errors/errors"
	"github.com/quantadex/quanta_book/common"
	"github.com/quantadex/quanta_book/consensus/tests"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

func GetProtocolByAddress(servers []*Cosi, address string) *Cosi {
	for _, p := range servers {
		if p.Node().Address() == address {
			return p
		}
	}
	return nil
}

func StartNServer(n int, verify []error) ([]*Cosi, []*tests.TestNode) {
	servers := []*Cosi{}
	nodes := []*tests.TestNode{}

	for i := 0; i < n; i++ {
		node := tests.NewTestNode("node:" + strconv.Itoa(i))
		node.Start()
		protocol := NewProtocol(node, i == 0, time.Duration(1*time.Second))
		node.SendMessage = func(sender, address string, call string, req interface{}, res interface{}) error {
			switch call {
			case "Protocol.Cosi":
				//fmt.Println("send to " + address)
				if p := GetProtocolByAddress(servers, address); p != nil {
					p.CosiMsgChan <- req.(*CosiMessage)
				}
			}
			return nil
		}

		func(i int) {
			protocol.Verify = func(msg string) error {
				return verify[i]
			}

			protocol.Persist = func(msg string, repair bool) error {
				fmt.Printf("Persist message node=%d repair=%v\n", i, repair)
				return nil
			}

			protocol.SignMsg = func(msg string) (string, error) {
				return msg, nil
			}

			protocol.Start()
		}(i)

		servers = append(servers, protocol)
		nodes = append(nodes, node)
	}

	for _, node := range nodes {
		node.GetAllNodes = func() []*tests.TestNode {
			return nodes
		}
	}

	return servers, nodes
}

// TEST 1:  1 leader, 1 follower
func TestConsensus2Nodes(t *testing.T) {
	common.InitLogger()

	protocols, _ := StartNServer(2, []error{nil, nil, nil})
	protocols[0].StartNewRound("tx message")
	finalMsg := <-protocols[0].FinalSigChan
	assert.Equal(t, 2, len(finalMsg.Msg), "Expect to have 2 signatures")
}

// TEST 2:  1 leader, 2 follower
// TEST 3:  1 leader, 1 follower (1 fail)  2/3
func TestConsensus3Nodes(t *testing.T) {
	common.InitLogger()
	protocols, _ := StartNServer(3, []error{nil, nil, nil})
	protocols[0].StartNewRound("tx message")
	finalMsg := <-protocols[0].FinalSigChan

	assert.Equal(t, 3, len(finalMsg.Msg), "Expect to have 3 signatures")
	println("Start of #1")

	// stop 1 node, should fail
	//stopping 2 nodes as threshold is 2
	protocols[1].Node().Stop()
	protocols[2].Node().Stop()
	protocols[0].StartNewRound("tx message 2")
	finalMsg = <-protocols[0].FinalSigChan
	assert.Equal(t, 0, len(finalMsg.Msg), "Expect to have 2 signatures")

	println("Start of #2")
	protocols[0].StartNewRound("tx message 2")
	finalMsg = <-protocols[0].FinalSigChan

	println("Start of #3")
	protocols[1].Node().Start()
	protocols[2].Node().Start()
	protocols[0].StartNewRound("tx message 2")
	finalMsg = <-protocols[0].FinalSigChan
	assert.Equal(t, 3, len(finalMsg.Msg), "Expect to have 2 signatures")
	protocols[0].Node().Stop()
	protocols[1].Node().Stop()
	protocols[2].Node().Stop()
}

// TEST 4:  1 leader, 4 nodes, 1 failure

//commenting to test other tests
//func TestConsensus4Nodes(t *testing.T) {
//	common.InitLogger()
//	protocols, _ := StartNServer(4, []error{nil, nil, nil, nil})
//	protocols[1].Node().Stop()
//
//	protocols[0].StartNewRound("tx message")
//	finalMsg := <-protocols[0].FinalSigChan
//
//	assert.Equal(t, 3, len(finalMsg.Msg), "Expect to have 3 signatures %v", finalMsg.Msg)
//}

//func TestConsensus3NodesDelayed(t *testing.T) {
//	common.InitLogger()
//	protocols,testNodes := StartNServer(3)
//	testNodes[0]
//	protocols[0].StartNewRound("tx message")
//	finalMsg := <- protocols[0].FinalSigChan
//
//	assert.Equal(t, 3, len(finalMsg.Msg), "Expect to have 3 signatures")
//
//}

func TestConsensus3Nodes1Fail(t *testing.T) {
	common.InitLogger()
	protocols, _ := StartNServer(3, []error{nil, errors.New("verify failure"), nil})
	protocols[0].StartNewRound("tx message")
	finalMsg := <-protocols[0].FinalSigChan

	assert.Equal(t, 2, len(finalMsg.Msg), "Expect to have 3 signatures")
}
