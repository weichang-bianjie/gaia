package gaia

import (
	"encoding/hex"
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gaia/v4/app/model"
	abci "github.com/tendermint/tendermint/abci/types"
	tentypes "github.com/tendermint/tendermint/types"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2/txn"
	stdlog "log"
)

type MongoDbService struct {
	//todo start_height,end_height
	header  tentypes.Header
	txCache []model.DeliverTx // the cache that write tx out to
}

func (fss *MongoDbService) ListenBeginBlock(ctx sdk.Context, req abci.RequestBeginBlock, res abci.ResponseBeginBlock) {
	var err error
	fss.header, err = tentypes.HeaderFromProto(&req.Header)
	if err != nil {
		ctx.Logger().Error(err.Error())
		return
	}
	////add latest block
	//latest := model.LastestBlock{
	//	Height:   fss.header.Height,
	//	Time:     fss.header.Time.Unix(),
	//	Hash:     fss.header.Hash().String(),
	//	Proposer: fss.header.ProposerAddress.String(),
	//}
	// NOTE: this could either be done synchronously or asynchronously
	// create a new file with the req info according to naming schema
	// write req to file
	// write all state changes cached for this stage to file
	// reset cache
	// write res to file
	// close file
}

func (fss *MongoDbService) ListenEndBlock(ctx sdk.Context, req abci.RequestEndBlock, res abci.ResponseEndBlock) {
	// NOTE: this could either be done synchronously or asynchronously
	// create a new file with the req info according to naming schema
	// write req to file
	// write all state changes cached for this stage to file
	// reset cache
	// write res to file
	// close file
	defer func() {
		fss.txCache = make([]model.DeliverTx, 0)
		fss.header = tentypes.Header{}
	}()
	var ops []txn.Op

	if len(fss.txCache) > 0 {
		txs := model.Txs{
			Id:  bson.NewObjectId(),
			Txs: fss.txCache,
		}
		ops = append(ops, model.InsertOps(txs.Name(), txs))
	}
	//todo set hook tx start_height and end_height
	block := model.Block{
		Id:       bson.NewObjectId(),
		Height:   fss.header.Height,
		Time:     fss.header.Time.Unix(),
		Hash:     fss.header.Hash().String(),
		Txn:      int64(len(fss.txCache)),
		Proposer: fss.header.ProposerAddress.String(),
	}
	ops = append(ops, model.InsertOps(block.Name(), block))
	if err := model.Txn(ops); err != nil {
		stdlog.Println("Failed to save block and tx", "height", fss.header.Height, "err", err)
	}

}

func (fss *MongoDbService) ListenDeliverTx(ctx sdk.Context, req abci.RequestDeliverTx, res abci.ResponseDeliverTx) {
	// NOTE: this could either be done synchronously or asynchronously
	// create a new file with the req info according to naming schema
	// NOTE: if the tx failed, handle accordingly
	// write req to file
	// write all state changes cached for this stage to file
	// reset cache
	// write res to file
	// close file

	//todo set hook tx start_height and end_height

	txResult, err := json.Marshal(res)
	if err == nil {
		fss.txCache = append(fss.txCache, model.DeliverTx{
			Tx:       hex.EncodeToString(req.Tx),
			TxResult: string(txResult),
		})
	}

}

// NewMongoDbService creates a new MongoDbService for the provided writeDir, (optional) filePrefix, and storeKeys
func NewMongoDbService(conf model.DataBaseConf) baseapp.Hook {
	model.Init(&conf)
	return &MongoDbService{
		txCache: make([]model.DeliverTx, 0),
	}
}
