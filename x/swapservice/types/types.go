package types

import (
	"fmt"
	"strings"
)

// PoolStruct is a struct that contains all the metadata of a pooldata
type PoolStruct struct {
	BalanceAtom  string `json:"balance_atom"`
	BalanceToken string `json:"balance_token"`
	Ticker       string `json:"ticker"`
	TokenName    string `json:"token_name"`
}

// Returns a new PoolStruct
func NewPoolStruct() PoolStruct {
	return PoolStruct{}
}

// implement fmt.Stringer
func (w PoolStruct) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Token: %s (%s) 
Atom Balance: %s
Token Balance: %s`, w.TokenName, w.Ticker, w.BalanceAtom, w.BalanceToken))
}

type Holding struct {
	Ticker string `json:"ticker"`
	Amount string `json:"amount"`
}

//AccStruct is a struct that contains balances of an account
type AccStruct struct {
	Name     string    `json:"name"`
	Holdings []Holding `json:"holdings"`
}

func NewAccStruct() AccStruct {
	return AccStruct{}
}

func (w AccStruct) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Balances for %s:
%+v
`, w.Name, w.Holdings))
}

type AccStake struct {
	Name  string `json:"name"`
	Atom  string `json:"atom"`
	Token string `json:"token"`
}

// Stake Struct is a struct that contain amount of coins stake towards a specific pool
type StakeStruct struct {
	Ticker string     `json:"ticker"`
	Stakes []AccStake `json:"stakes"`
}

func NewStakeStruct() StakeStruct {
	return StakeStruct{
		Stakes: make([]AccStake, 0),
	}
}

func (w StakeStruct) String() string {
	// TODO: Print better stakes
	return strings.TrimSpace(fmt.Sprintf("TODO: Print better stakes"))
}
