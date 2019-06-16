package swapservice

import (
	"fmt"
	"log"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/bank"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	coinKeeper bank.Keeper

	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context

	cdc *codec.Codec // The wire codec for binary encoding/decoding.
}

// NewKeeper creates new instances of the swapservice Keeper
func NewKeeper(coinKeeper bank.Keeper, storeKey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		coinKeeper: coinKeeper,
		storeKey:   storeKey,
		cdc:        cdc,
	}
}

// Gets the entire AccStruct metadata struct for a acc ID
func (k Keeper) GetAccStruct(ctx sdk.Context, accID string) AccStruct {
	if !strings.HasPrefix(accID, "acc-") {
		accID = fmt.Sprintf("acc-%s", accID)
	}
	store := ctx.KVStore(k.storeKey)
	if !store.Has([]byte(accID)) {
		return NewAccStruct()
	}
	bz := store.Get([]byte(accID))
	var accstruct AccStruct
	k.cdc.MustUnmarshalBinaryBare(bz, &accstruct)
	return accstruct
}

// Sets the entire AccStruct metadata struct for a acc ID
func (k Keeper) SetAccStruct(ctx sdk.Context, accID string, accstruct AccStruct) {
	if !strings.HasPrefix(accID, "acc-") {
		accID = fmt.Sprintf("acc-%s", accID)
	}
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(accID), k.cdc.MustMarshalBinaryBare(accstruct))
}

// SetAccData - sets the value string that a acc ID resolves to
func (k Keeper) SetAccData(ctx sdk.Context, accID string, name, ticker, amount string) {
	if !strings.HasPrefix(accID, "acc-") {
		accID = fmt.Sprintf("acc-%s", accID)
	}
	accstruct := k.GetAccStruct(ctx, accID)
	found := false
	ticker = strings.ToUpper(ticker)
	for i, record := range accstruct.Holdings {
		if record.Ticker == ticker {
			accstruct.Holdings[i].Amount = amount
			found = true
			break
		}
	}
	if !found {
		record := Holding{
			Ticker: ticker,
			Amount: amount,
		}
		accstruct.Holdings = append(accstruct.Holdings, record)
	}
	k.SetAccStruct(ctx, accID, accstruct)
}

func (k Keeper) GetAccData(ctx sdk.Context, accID, ticker string) string {
	if !strings.HasPrefix(accID, "acc-") {
		accID = fmt.Sprintf("acc-%s", accID)
	}
	accstruct := k.GetAccStruct(ctx, accID)
	ticker = strings.ToUpper(ticker)
	for _, record := range accstruct.Holdings {
		if record.Ticker == ticker {
			return record.Amount
		}
	}
	return ""
}

// Gets the entire StakeStruct metadata struct for a stake ID
func (k Keeper) GetStakeStruct(ctx sdk.Context, stakeID string) StakeStruct {
	if !strings.HasPrefix(stakeID, "stake-") {
		stakeID = fmt.Sprintf("stake-%s", stakeID)
	}
	store := ctx.KVStore(k.storeKey)
	if !store.Has([]byte(stakeID)) {
		return NewStakeStruct()
	}
	bz := store.Get([]byte(stakeID))
	var stakestruct StakeStruct
	k.cdc.MustUnmarshalBinaryBare(bz, &stakestruct)
	return stakestruct
}

// Get stake data for a specific user
func (k Keeper) GetStakeData(ctx sdk.Context, stakeID, name string) AccStake {
	if !strings.HasPrefix(stakeID, "stake-") {
		stakeID = fmt.Sprintf("stake-%s", stakeID)
	}
	stakestruct := k.GetStakeStruct(ctx, stakeID)
	for _, record := range stakestruct.Stakes {
		if record.Name == name {
			return record
		}
	}
	return AccStake{
		Name:  name,
		Atom:  "0",
		Token: "0",
	}
}

// Sets the entire StakeStruct metadata struct for a stake ID
func (k Keeper) SetStakeStruct(ctx sdk.Context, stakeID string, stakestruct StakeStruct) {
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(stakeID), k.cdc.MustMarshalBinaryBare(stakestruct))
}

// SetStakeData - sets the value string that a stake ID resolves to
func (k Keeper) SetStakeData(ctx sdk.Context, stakeID string, name, atom, token string) {
	parts := strings.Split(stakeID, "-")
	stakestruct := k.GetStakeStruct(ctx, stakeID)
	stakestruct.Ticker = parts[1]
	found := false
	for i, record := range stakestruct.Stakes {
		if record.Name == name {
			stakestruct.Stakes[i].Atom = atom
			stakestruct.Stakes[i].Token = token
			found = true
			break
		}
	}
	if !found {
		record := AccStake{
			Name:  name,
			Atom:  atom,
			Token: token,
		}
		stakestruct.Stakes = append(stakestruct.Stakes, record)
	}
	log.Printf("Saving struct: %s", stakeID)
	log.Printf("Struct: %+v", stakestruct)
	k.SetStakeStruct(ctx, stakeID, stakestruct)
}

// Gets the entire PoolStruct metadata struct for a pool ID
func (k Keeper) GetPoolStruct(ctx sdk.Context, poolID string) PoolStruct {
	if !strings.HasPrefix(poolID, "pool-") {
		poolID = fmt.Sprintf("pool-%s", poolID)
	}
	store := ctx.KVStore(k.storeKey)
	if !store.Has([]byte(poolID)) {
		return NewPoolStruct()
	}
	bz := store.Get([]byte(poolID))
	var poolstruct PoolStruct
	k.cdc.MustUnmarshalBinaryBare(bz, &poolstruct)
	if poolstruct.BalanceAtom == "" {
		poolstruct.BalanceAtom = "0"
	}
	if poolstruct.BalanceToken == "" {
		poolstruct.BalanceToken = "0"
	}
	return poolstruct
}

// Sets the entire PoolStruct metadata struct for a pool ID
func (k Keeper) SetPoolStruct(ctx sdk.Context, poolID string, poolstruct PoolStruct) {
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(poolID), k.cdc.MustMarshalBinaryBare(poolstruct))
}

// GetPool - gets the balances of a pool. Specifying ticker dictates which
// balance is return in 0 vs 1 spot.
func (k Keeper) GetPoolData(ctx sdk.Context, poolID, ticker string) (string, string) {
	poolstruct := k.GetPoolStruct(ctx, poolID)
	if strings.ToUpper(ticker) == "ATOM" {
		return poolstruct.BalanceAtom, poolstruct.BalanceToken
	}
	return poolstruct.BalanceToken, poolstruct.BalanceAtom
}

// SetPoolData - sets the value string that a pool ID resolves to
func (k Keeper) SetPoolData(ctx sdk.Context, poolID string, tokenName, ticker, balanceAtom, balanceToken string) {
	poolstruct := k.GetPoolStruct(ctx, poolID)
	poolstruct.TokenName = tokenName
	poolstruct.Ticker = strings.ToUpper(ticker)
	poolstruct.BalanceAtom = balanceAtom
	poolstruct.BalanceToken = balanceToken
	k.SetPoolStruct(ctx, poolID, poolstruct)
}

// SetBalances - sets the current balances of a pool
func (k Keeper) SetBalances(ctx sdk.Context, poolID, atom, token string) {
	poolstruct := k.GetPoolStruct(ctx, poolID)
	poolstruct.BalanceAtom = atom
	poolstruct.BalanceToken = token
	k.SetPoolStruct(ctx, poolID, poolstruct)
}

// Get an iterator over all pool IDs in which the keys are the pool IDs and the values are the poolstruct
func (k Keeper) GetDatasIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, nil)
}
