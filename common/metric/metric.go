package metric

import (
	"encoding/json"
	"errors"
	"expvar"
	"github.com/quantadex/distributed_quanta_bridge/common/kv_store"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"github.com/quantadex/distributed_quanta_bridge/trust/control"
	"github.com/quantadex/distributed_quanta_bridge/trust/db"
	"github.com/quantadex/distributed_quanta_bridge/trust/quanta"

	//"github.com/quantadex/distributed_quanta_bridge/trust/quanta"
	"time"
)

const NORMAL = "normal"
const DEGRADED = "degraded"
const FAILURE = "failure"

type BlockchainStatus struct {
	CurrentBlock        int64
	TimeSinceLastBlock  float64 // sec
	DegradedThreshold   int64   //
	FailureThreshold    int64   //
	TotalAddresses      int
	AddressesCreated24H int
	State               string
}

type Counter struct {
	Interval int64    `json:"interval"`
	Samples  []Sample `json:"samples"`
}

type Sample struct {
	Type  string `json:"type"`
	Count int64  `json:"count"`
}

type TransactionStatus struct {
	ConsensusRetries  int64
	DegradedThreshold int64
	FailureThreshold  int64
	State             string
}

func GetBlockchainStatus(coins interface{}, store kv_store.KVStore, rDb *db.DB, degradedThreshold, failureThreshold int64) (BlockchainStatus, error) {
	res := BlockchainStatus{}
	res.DegradedThreshold = degradedThreshold
	res.FailureThreshold = failureThreshold

	var coinName string
	flag := false
	switch v := coins.(type) {
	case coin.Coin:
		coinName = v.Blockchain()
		flag = true
	case quanta.Quanta:
		coinName = control.QUANTA
	default:
		return res, errors.New("unknown type")
	}

	lastProcessed, valid := control.GetLastBlock(store, coinName)
	if valid {
		res.CurrentBlock = lastProcessed
	}

	var blockTime time.Time
	var err error
	if flag == true {
		blockTime, err = coins.(coin.Coin).GetBlockTime(lastProcessed)
		if err != nil {
			return res, err
		}
	} else {
		blockTime, err = coins.(quanta.Quanta).GetBlockTime(lastProcessed)
		if err != nil {
			return res, err
		}
	}

	state, timeDiff := GetStateAndTime(blockTime, degradedThreshold, failureThreshold)
	res.State = state
	res.TimeSinceLastBlock = timeDiff.Seconds()

	totalAddresses, err := rDb.GetAddressCountByBlockchain(coinName)
	if err != nil {

	}
	res.TotalAddresses = totalAddresses
	addressesSince24hr, err := rDb.GetAddressCountByBlockchainAndTime(coinName)
	if err != nil {

	}
	res.AddressesCreated24H = addressesSince24hr
	return res, nil
}

func GetStateAndTime(blockTime time.Time, degradedThreshold, failureThreshold int64) (string, time.Duration) {
	timeDiff := time.Since(blockTime)

	if timeDiff.Seconds() > (time.Second * time.Duration(failureThreshold)).Seconds() {
		return FAILURE, timeDiff
	}
	if timeDiff.Seconds() > (time.Second * time.Duration(degradedThreshold)).Seconds() {
		return DEGRADED, timeDiff
	}
	return NORMAL, timeDiff
}

func IncrFailuresAndDegraded(state string, totalDegraded, totalFailure *int64) {
	if state == FAILURE {
		*totalFailure = *totalFailure + 1
	} else if state == DEGRADED {
		*totalDegraded = *totalDegraded + 1
	}
}

func GetDepositOrWithdrawalStatus(txType string, degradedThreshold, failureThreshold int64) (TransactionStatus, error) {
	var v expvar.Var
	res := TransactionStatus{}

	res.DegradedThreshold = degradedThreshold
	res.FailureThreshold = failureThreshold

	if txType == db.DEPOSIT {
		v = expvar.Get("deposit_status")
	} else if txType == db.WITHDRAWAL {
		v = expvar.Get("withdrawal_status")
	} else {
		return res, errors.New("unknown transaction type")
	}
	if v != nil {
		var counter *Counter
		err := json.Unmarshal([]byte(v.String()), &counter)
		if err != nil {
			return res, err
		}
		count := counter.Samples[0].Count
		res.ConsensusRetries = count

		state := ""
		if count > failureThreshold {
			state = FAILURE
		} else if count > degradedThreshold {
			state = DEGRADED
		} else {
			state = NORMAL
		}
		res.State = state
	}
	return res, nil
}
