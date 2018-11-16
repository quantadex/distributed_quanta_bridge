package quanta

import (
	"github.com/quantadex/distributed_quanta_bridge/common/kv_store"
	"github.com/stellar/go/clients/horizon"
	"net/http"
	"time"
	"github.com/quantadex/distributed_quanta_bridge/trust/db"
)

type SubmitWorkerImpl struct {
	horizonClient *horizon.Client
	QuantaClientOptions
}

func (s *SubmitWorkerImpl) Dispatch() {
	if s.Logger == nil {
		panic("Missing logger")
	}
	s.Logger.Infof("Submitworker started")

	for {
		//println("Wake up")
		time.Sleep(time.Second)

		data := db.QueryDepositByAge(s.Db, time.Now(), []string{db.SUBMIT_QUEUE})

		for k, v := range data {

			s.Logger.Infof("Submit TX: %s signed=%v %v", v.Tx, v.Signed, v.SubmitTx)

			res, err := s.horizonClient.SubmitTransaction(v.SubmitTx)
			if err != nil {
				err2 := err.(*horizon.Error)
				s.Logger.Error("could not submit transaction " + err2.Error() + err2.Problem.Detail)

				txError := FailedTransactionError{res.Result}
				opCodes, _ := txError.OperationResultCodes()
				txCodes, _ := txError.TransactionResultCode()
				s.Logger.Errorf("Op codes %v Tx codes %v", opCodes, txCodes)
			} else {
				s.Logger.Infof("Successful tx submission %s,remove %s", res.Hash, k)
				err = db.ChangeSubmitState(s.Db, v.Tx, db.SUBMIT_SUCCESS)
				if err != nil {
					s.Logger.Error("Error removing key=" + v.Tx)
				}
			}

		}
	}
}

func (s *SubmitWorkerImpl) AttachQueue(kv kv_store.KVStore) error {
	s.horizonClient = &horizon.Client{
		URL:  s.HorizonUrl,
		HTTP: http.DefaultClient,
	}

	return nil
}
