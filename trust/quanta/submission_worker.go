package quanta

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/quantadex/distributed_quanta_bridge/common/kv_store"
	"github.com/quantadex/distributed_quanta_bridge/trust/db"
	"github.com/stellar/go/clients/horizon"
	"net/http"
	"time"
)

type SubmitWorkerImpl struct {
	horizonClient *horizon.Client
	QuantaClientOptions
}

func ErrorString(err error, showStackTrace ...bool) string {
	var errorString string
	herr, isHorizonError := errors.Cause(err).(*horizon.Error)

	if isHorizonError {
		errorString += fmt.Sprintf("%v: %v", herr.Problem.Status, herr.Problem.Title)

		resultCodes, err := herr.ResultCodes()
		if err == nil {
			errorString += fmt.Sprintf(" (%v)", resultCodes)
		}
	} else {
		errorString = fmt.Sprintf("%v", err)
	}

	if len(showStackTrace) > 0 {
		if isHorizonError {
			errorString += fmt.Sprintf("\nDetail: %s\nType: %s\n", herr.Problem.Detail, herr.Problem.Type)
		}
		errorString += fmt.Sprintf("\nStack trace:\n%+v\n", err)
	}

	return errorString
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
				s.Logger.Error("could not submit transaction " + ErrorString(err, false))
			} else {
				s.Logger.Infof("Successful tx submission %s,remove %s", res.Hash, k)
				err = db.ChangeSubmitState(s.Db, v.Tx, db.SUBMIT_SUCCESS, "")
				if err != nil {
					s.Logger.Error("Error removing key=" + v.Tx)
				}
			}

		}
	}
}

func (s *SubmitWorkerImpl) AttachQueue(kv kv_store.KVStore) error {
	s.horizonClient = &horizon.Client{
		//URL:  s.HorizonUrl,
		HTTP: http.DefaultClient,
	}

	return nil
}
