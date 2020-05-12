package types

import (
	"bytes"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	hmCommon "github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

// verify interface at compile time
var _ sdk.Msg = &MsgUnjail{}

// MsgUnjail - struct for unjailing jailed validator
type MsgUnjail struct {
	From        types.HeimdallAddress `json:"from"`
	ID          hmTypes.ValidatorID   `json:"id"`
	TxHash      types.HeimdallHash    `json:"tx_hash"`
	LogIndex    uint64                `json:"log_index"`
	BlockNumber uint64                `json:"block_number"`
}

func NewMsgUnjail(from types.HeimdallAddress, id uint64, txHash types.HeimdallHash, logIndex uint64, blockNumber uint64) MsgUnjail {
	return MsgUnjail{
		From:        from,
		ID:          hmTypes.NewValidatorID(id),
		TxHash:      txHash,
		LogIndex:    logIndex,
		BlockNumber: blockNumber,
	}
}

//nolint
func (msg MsgUnjail) Route() string { return RouterKey }
func (msg MsgUnjail) Type() string  { return "unjail" }
func (msg MsgUnjail) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{hmTypes.HeimdallAddressToAccAddress(msg.From)}
}

// GetSignBytes gets the bytes for the message signer to sign on
func (msg MsgUnjail) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// ValidateBasic validity check for the AnteHandler
func (msg MsgUnjail) ValidateBasic() sdk.Error {
	if msg.ID <= 0 {
		return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid validator ID %v", msg.ID)
	}

	if msg.From.Empty() {
		return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid proposer %v", msg.From.String())
	}
	return nil
}

// Tick Msg

// TickMsg - struct for unjailing jailed validator
type MsgTick struct {
	Proposer         types.HeimdallAddress `json:"proposer"`
	SlashingInfoHash types.HeimdallHash    `json:"slashInfoHash"`
}

func NewMsgTick(proposer types.HeimdallAddress, slashingInfoHash types.HeimdallHash) MsgTick {
	return MsgTick{
		Proposer:         proposer,
		SlashingInfoHash: slashingInfoHash,
	}
}

// Type returns message type
func (msg MsgTick) Type() string {
	return "tick"
}

func (msg MsgTick) Route() string {
	return RouterKey
}

// GetSigners returns address of the signer
func (msg MsgTick) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{types.HeimdallAddressToAccAddress(msg.Proposer)}
}

func (msg MsgTick) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgTick) ValidateBasic() sdk.Error {
	if bytes.Equal(msg.SlashingInfoHash.Bytes(), helper.ZeroHash.Bytes()) {
		return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid slashing info hash %v", msg.SlashingInfoHash.String())
	}

	if msg.Proposer.Empty() {
		return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid proposer %v", msg.Proposer.String())
	}
	return nil
}

// GetSideSignBytes returns side sign bytes
func (msg MsgTick) GetSideSignBytes() []byte {
	// keccak256(abi.encoded(proposer, startBlock, endBlock, rootHash, accountRootHash, bor chain id))
	/* 	return appendBytes32(
	   		msg.Proposer.Bytes(),
	   		msg.SlashingInfoHash.Bytes(),
	   	)
	*/
	return nil
}

//
// Msg Tick Ack
//

var _ sdk.Msg = &MsgTickAck{}

type MsgTickAck struct {
	From          types.HeimdallAddress `json:"from"`
	SlashedAmount *big.Int              `json:"slashed_amount"`
	TxHash        types.HeimdallHash    `json:"tx_hash"`
	LogIndex      uint64                `json:"log_index"`
	BlockNumber   uint64                `json:"block_number"`
}

func NewMsgTickAck(from types.HeimdallAddress, slashedAmount *big.Int, txHash types.HeimdallHash, logIndex uint64, blockNumber uint64) MsgTickAck {
	return MsgTickAck{
		From:          from,
		SlashedAmount: slashedAmount,
		TxHash:        txHash,
		BlockNumber:   blockNumber,
		LogIndex:      logIndex,
	}
}

// Type returns message type
func (msg MsgTickAck) Type() string {
	return "tick-ack"
}

func (msg MsgTickAck) Route() string {
	return RouterKey
}

// GetSigners returns address of the signer
func (msg MsgTickAck) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{types.HeimdallAddressToAccAddress(msg.From)}
}

func (msg MsgTickAck) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgTickAck) ValidateBasic() sdk.Error {
	if msg.From.Empty() {
		return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid from %v", msg.From.String())
	}
	return nil
}

// GetTxHash Returns tx hash
func (msg MsgTickAck) GetTxHash() types.HeimdallHash {
	return msg.TxHash
}

// GetLogIndex Returns log index
func (msg MsgTickAck) GetLogIndex() uint64 {
	return msg.LogIndex
}

// GetSideSignBytes returns side sign bytes
func (msg MsgTickAck) GetSideSignBytes() []byte {
	return nil
}