package stellar

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/trustwallet/blockatlas"
	"github.com/trustwallet/blockatlas/coin"
	"github.com/trustwallet/blockatlas/util"
	"net/http"
	"strconv"
	"time"
)

type Platform struct {
	client Client
	CoinIndex uint
}

func (p *Platform) Init() error {
	handle := coin.Coins[p.CoinIndex].Handle
	p.client.API = viper.GetString(fmt.Sprintf("%s.api", handle))
	p.client.HTTP = &http.Client{
		Timeout: 2 * time.Second,
	}
	return nil
}

func (p *Platform) Coin() coin.Coin {
	return coin.Coins[p.CoinIndex]
}

func (p *Platform) GetTxsByAddress(address string) (blockatlas.TxPage, error) {
	payments, err := p.client.GetTxsOfAddress(address)
	if err != nil {
		return nil, err
	}

	var txs []blockatlas.Tx
	for _, payment := range payments {
		tx, ok := Normalize(&payment, p.CoinIndex)
		if !ok {
			continue
		}
		txs = append(txs, tx)
	}

	return txs, nil
}

// Normalize converts a Stellar-based transaction into the generic model
func Normalize(payment *Payment, nativeCoinIndex uint) (tx blockatlas.Tx, ok bool) {
	switch payment.Type {
	case "payment":
		if payment.AssetType != "native" {
			return tx, false
		}
	case "create_account":
		break
	default:
		return tx, false
	}
	id, err := strconv.ParseUint(payment.ID, 10, 64)
	if err != nil {
		return tx, false
	}
	date, err := time.Parse("2006-01-02T15:04:05Z", payment.CreatedAt)
	if err != nil {
		return tx, false
	}
	var value, from, to string
	if payment.Amount != "" {
		value, err = util.DecimalToSatoshis(payment.Amount)
		from = payment.From
		to = payment.To
	} else if payment.StartingBalance != "" {
		value, err = util.DecimalToSatoshis(payment.StartingBalance)
		from = payment.Funder
		to = payment.Account
	} else {
		return tx, false
	}
	if err != nil {
		return tx, false
	}
	return blockatlas.Tx{
		ID:   payment.TransactionHash,
		Coin: nativeCoinIndex,
		From: from,
		To:   to,
		// https://www.stellar.org/developers/guides/concepts/fees.html
		// Fee fixed at 100 stroops
		Fee:   "100",
		Date:  date.Unix(),
		Block: id,
		Meta:  blockatlas.Transfer{
			Value: blockatlas.Amount(value),
		},
	}, true
}
