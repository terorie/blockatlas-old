package coin

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

//go:generate rm -f list.go
//go:generate go run gen.go

// Coin is the native currency of a blockchain
type Coin struct {
	ID          uint   `yaml:"id"`            // SLIP-44 ID (e.g. 242)
	Handle      string `yaml:"handle"`        // Trust Wallet handle (e.g. nimiq)
	Symbol      string `yaml:"symbol"`        // Symbol of native currency
	Title       string `yaml:"name"`          // Full name of native currency
	Decimals    uint   `yaml:"decimals"`      // Number of decimals
	BlockTime   int    `yaml:"blockTime"`     // Average time between blocks (ms)
	SampleAddr  string `yaml:"sampleAddress"` // Random address seen on chain (optional)
	SampleToken string `yaml:"sampleToken"`   // Random token seen on chain (optional)
}

func (c Coin) String() string {
	return fmt.Sprintf("[%s] %s (#%d)", c.Symbol, c.Title, c.ID)
}

func Load(fPath string) (list []Coin, err error) {
	f, err := os.Open(fPath)
	if err != nil {
		return nil, err
	}
	dec := yaml.NewDecoder(f)
	err = dec.Decode(&list)
	if err != nil {
		return nil, err
	}

	return list, nil
}
