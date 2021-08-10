package model

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2/txn"
	"strings"
	"time"
)

const (
	CollectionNameTxn = "sync_txn"
)

var (
	session *mgo.Session
	_conf   *DataBaseConf
)

type DataBaseConf struct {
	Addrs    string
	User     string
	Passwd   string `json:"-"`
	Database string
	ChainId  string
}

func getChainId() string {
	return _conf.ChainId
}

func Init(conf *DataBaseConf) {
	_conf = conf
	addrs := strings.Split(conf.Addrs, ",")
	dialInfo := &mgo.DialInfo{
		Addrs:     addrs,
		Database:  conf.Database,
		Username:  conf.User,
		Password:  conf.Passwd,
		Direct:    true,
		Timeout:   time.Second * 10,
		PoolLimit: 4096, // Session.SetPoolLimit
	}

	var err error
	session, err = mgo.DialWithInfo(dialInfo)
	if err != nil {

	}
	session.SetMode(mgo.Strong, true)
}

func Close() {
	session.Close()
}

func getSession() *mgo.Session {
	// max session num is 4096
	return session.Clone()
}

func _getTxnName() string {
	if _conf.ChainId == "" {
		return CollectionNameTxn
	}
	return fmt.Sprintf("sync_%v_txn", _conf.ChainId)
}

//mgo transaction method
//detail to see: https://godoc.org/gopkg.in/mgo.v2/txn
func Txn(ops []txn.Op) error {
	session := getSession()
	defer session.Close()

	c := session.DB(_conf.Database).C(_getTxnName())
	runner := txn.NewRunner(c)

	txObjectId := bson.NewObjectId()
	err := runner.Run(ops, txObjectId, nil)
	if err != nil {
		if err == txn.ErrAborted {
			err = runner.Resume(txObjectId)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}

func InsertOps(collectionName string, data interface{}) txn.Op {
	return txn.Op{
		C:      collectionName,
		Id:     bson.NewObjectId(),
		Insert: data,
	}
}

func UpdateOps(collectionName string, Id interface{}, data interface{}) txn.Op {
	return txn.Op{
		C:  collectionName,
		Id: Id,
		Update: bson.M{
			"$set": data,
		},
	}
}
