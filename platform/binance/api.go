package binance

import (
	"fmt"
	"github.com/trustwallet/blockatlas"
	"net/http"

	"github.com/spf13/viper"
	"github.com/trustwallet/blockatlas/coin"
	"github.com/trustwallet/blockatlas/util"
)

type Platform struct {
	client Client
}

func (p *Platform) Init() error {
	p.client.BaseURL = viper.GetString("binance.api")
	p.client.HTTPClient = http.DefaultClient
	return nil
}

func (p *Platform) Coin() coin.Coin {
	return coin.Coins[coin.BNB]
}

func (p *Platform) CurrentBlockNumber() (int64, error) {
	// No native function to get height in explorer API
	// Workaround: Request list of blocks
	// and return number of the newest one
	list, err := p.client.GetBlockList(1)
	if err != nil {
		return 0, err
	}
	if len(list.BlockArray) == 0 {
		return 0, fmt.Errorf("no block descriptor found")
	}
	return list.BlockArray[0].BlockHeight, nil
}

func (p *Platform) GetBlockByNumber(num int64) (*blockatlas.Block, error) {
	srcTxs, err := p.client.GetBlockByNumber(num)
	if err != nil {
		return nil, err
	}
	// TODO: Only returns BNB transactions for now
	txs := NormalizeTxs(srcTxs.Txs, "")
	return &blockatlas.Block{
		Number: num,
		Txs:    txs,
	}, nil
}

func (p *Platform) GetTxsByAddress(address string) (blockatlas.TxPage, error) {
	// Endpoint supports queries without token query parameter
	return p.GetTokenTxsByAddress(address, "")
}

func (p *Platform) GetTokenTxsByAddress(address string, token string) (blockatlas.TxPage, error) {
	srcTxs, err := p.client.GetTxsOfAddress(address, token)
	if err != nil {
		return nil, err
	}
	return NormalizeTxs(srcTxs.Txs, token), nil
}

// NormalizeTx converts a Binance transaction into the generic model
func NormalizeTx(srcTx *Tx, token string) (tx blockatlas.Tx, ok bool) {
	value := util.DecimalExp(string(srcTx.Value), 8)
	fee := util.DecimalExp(string(srcTx.Fee), 8)

	tx = blockatlas.Tx{
		ID:    srcTx.Hash,
		Coin:  coin.BNB,
		Date:  srcTx.Timestamp / 1000,
		From:  srcTx.FromAddr,
		To:    srcTx.ToAddr,
		Fee:   blockatlas.Amount(fee),
		Block: srcTx.BlockHeight,
		Memo:  srcTx.Memo,
	}

	// Condition for native transfer (BNB)
	if srcTx.Asset == "BNB" && srcTx.Type == "TRANSFER" && token == "" {
		tx.Meta = blockatlas.Transfer{
			Value: blockatlas.Amount(value),
		}
		return tx, true
	}

	// Condition for native token transfer
	if srcTx.Asset == token && srcTx.Type == "TRANSFER" {
		tx.Meta = blockatlas.NativeTokenTransfer{
			TokenID:  srcTx.Asset,
			Symbol:   srcTx.MappedAsset,
			Value:    blockatlas.Amount(value),
			Decimals: 8,
			From:     srcTx.FromAddr,
			To:       srcTx.ToAddr,
		}

		return tx, true
	}

	return tx, false
}

// NormalizeTxs converts multiple Binance transactions
func NormalizeTxs(srcTxs []Tx, token string) (txs []blockatlas.Tx) {
	for _, srcTx := range srcTxs {
		tx, ok := NormalizeTx(&srcTx, token)
		if !ok || len(txs) >= blockatlas.TxPerPage {
			continue
		}
		txs = append(txs, tx)
	}
	return
}
