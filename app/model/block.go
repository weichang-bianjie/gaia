package model

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
)

const (
	CollectionNameBlock = "sync_block"
)

type (
	Block struct {
		Id       bson.ObjectId `bson:"_id"`
		Height   int64         `bson:"height" json:"height"`
		Hash     string        `bson:"hash" json:"hash"`
		Txn      int64         `bson:"txn" json:"txn"`
		Time     int64         `bson:"time" json:"time"`
		Proposer string        `bson:"proposer" json:"proposer"`
	}
)

func (d Block) Name() string {
	if getChainId() == "" {
		return CollectionNameBlock
	}
	return fmt.Sprintf("sync_%v_block", getChainId())
}
