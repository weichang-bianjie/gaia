package model

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
)

const (
	CollectionNameTx = "sync_tx"
)

type (
	DeliverTx struct {
		Tx       string `bson:"tx" json:"tx"`
		TxResult string `bson:"tx_result" json:"tx_result"`
	}
	Txs struct {
		Id  bson.ObjectId `bson:"_id"`
		Txs []DeliverTx   `bson:"txs" json:"txs"`
	}
)

func (d Txs) Name() string {
	if getChainId() == "" {
		return CollectionNameTx
	}
	return fmt.Sprintf("sync_%v_tx", getChainId())
}
