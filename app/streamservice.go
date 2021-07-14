package gaia

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"os"
	"strings"
)

// FileStreamingService is a concrete implementation of StreamingService that writes state changes out to a file
type FileStreamingService struct {
	txCache    []string        // the cache that write tx out to
	filePrefix string          // optional prefix for each of the generated files
	writeDir   string          // directory to write files into
	marshaller codec.Marshaler // marshaller used for re-marshalling the ABCI messages to write them out to the destination files
}

func (fss *FileStreamingService) ListenBeginBlock(ctx sdk.Context, req abci.RequestBeginBlock, res abci.ResponseBeginBlock) {
	if strings.Contains(fss.filePrefix, "_") {
		fss.filePrefix = fss.filePrefix[:strings.Index(fss.filePrefix, "_")]
	}
	fss.filePrefix = fmt.Sprint(fss.filePrefix, "_", req.Header.Height)
	// NOTE: this could either be done synchronously or asynchronously
	// create a new file with the req info according to naming schema
	// write req to file
	// write all state changes cached for this stage to file
	// reset cache
	// write res to file
	// close file
}

func (fss *FileStreamingService) ListenEndBlock(ctx sdk.Context, req abci.RequestEndBlock, res abci.ResponseEndBlock) {
	// NOTE: this could either be done synchronously or asynchronously
	// create a new file with the req info according to naming schema
	// write req to file
	// write all state changes cached for this stage to file
	// reset cache
	// write res to file
	// close file
	if len(fss.txCache) > 0 {
		filename := fmt.Sprint(fss.writeDir, "/", fss.filePrefix, "_txs")
		file, err := os.Create(filename)
		if err != nil {
			ctx.Logger().Error(err.Error())
			return
		}

		file.Write([]byte(strings.Join(fss.txCache, "\n")))
		file.Close()
		fss.txCache = make([]string, 0)
	}

}

func (fss *FileStreamingService) ListenDeliverTx(ctx sdk.Context, req abci.RequestDeliverTx, res abci.ResponseDeliverTx) {
	// NOTE: this could either be done synchronously or asynchronously
	// create a new file with the req info according to naming schema
	// NOTE: if the tx failed, handle accordingly
	// write req to file
	// write all state changes cached for this stage to file
	// reset cache
	// write res to file
	// close file

	type DeliverTx struct {
		Tx       string                 `json:"tx"`
		TxResult abci.ResponseDeliverTx `json:"tx_result"`
	}

	data, _ := json.Marshal(DeliverTx{
		Tx:       hex.EncodeToString(req.Tx),
		TxResult: res,
	})
	fss.txCache = append(fss.txCache, string(data))

}

// NewFileStreamingService creates a new FileStreamingService for the provided writeDir, (optional) filePrefix, and storeKeys
func NewFileStreamingService(writeDir, filePrefix string) baseapp.Hook {
	var m codec.Marshaler
	interfaceRegistry := types.NewInterfaceRegistry()
	m = codec.NewProtoCodec(interfaceRegistry)
	return &FileStreamingService{
		filePrefix: filePrefix,
		writeDir:   writeDir,
		marshaller: m,
		txCache:    make([]string, 0),
	}
}
