package milestone

import (
	"bytes"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/milestone/types"
)

// NewHandler creates new handler for handling messages for milestone module
func NewHandler(k Keeper, contractCaller helper.IContractCaller) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {

		case types.MsgMilestone:
			return handleMsgMilestone(ctx, msg, k, contractCaller)
		default:
			return sdk.ErrTxDecode("Invalid message in milestone module").Result()
		}
	}
}

// handleMsgMilestone Validates milestone transaction
func handleMsgMilestone(ctx sdk.Context, msg types.MsgMilestone, k Keeper, contractCaller helper.IContractCaller) sdk.Result {
	logger := k.Logger(ctx)
	milestoneLength := helper.MilestoneLength

	//
	//Check for the msg milestone
	//

	if ctx.BlockHeight() < helper.GetMilestoneHardForkHeight() {
		logger.Error("Network hasn't reached the", "Hard forked height", helper.GetMilestoneHardForkHeight())
		return common.ErrInvalidMsg(k.Codespace(), "Network hasn't reached the milestone hard forked height").Result()
	}

	if msg.StartBlock+milestoneLength-1 != msg.EndBlock {
		logger.Error("Milestone's length doesn't match the  milestone length set in configuration",
			"StartBlock", msg.StartBlock,
			"EndBlock", msg.EndBlock,
			"Milestone Length", milestoneLength,
		)

		return common.ErrMilestoneInvalid(k.Codespace()).Result()
	}

	// fetch last milestone from store
	if lastMilestone, err := k.GetLastMilestone(ctx); err == nil {
		// make sure new milestone is in continuity
		if lastMilestone.EndBlock+1 != msg.StartBlock {
			logger.Error("Milestone not in continuity ",
				"lastMilestoneEndBlock", lastMilestone.EndBlock,
				"receivedMsgStartBlock", msg.StartBlock,
			)

			return common.ErrMilestoneNotInContinuity(k.Codespace()).Result()
		}

	} else if err != nil && msg.StartBlock != helper.GetMilestoneBorBlockHeight() {
		logger.Error("First milestone to start from block %d", helper.GetMilestoneBorBlockHeight(), "milestone start block", msg.StartBlock, "error", err)
		return common.ErrNoMilestoneFound(k.Codespace()).Result()

	}

	//
	// Validate proposer
	//

	// Check proposer in message
	validatorSet := k.sk.GetValidatorSet(ctx)
	if validatorSet.Proposer == nil {
		logger.Error("No proposer in validator set", "msgProposer", msg.Proposer.String())
		return common.ErrInvalidMsg(k.Codespace(), "No proposer in stored validator set").Result()
	}

	if !bytes.Equal(msg.Proposer.Bytes(), validatorSet.Proposer.Signer.Bytes()) {
		logger.Error(
			"Invalid proposer in msg",
			"proposer", validatorSet.Proposer.Signer.String(),
			"msgProposer", msg.Proposer.String(),
		)

		return common.ErrInvalidMsg(k.Codespace(), "Invalid proposer in msg").Result()
	}

	// Emit event for milestone
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeMilestone,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyProposer, msg.Proposer.String()),
			sdk.NewAttribute(types.AttributeKeyStartBlock, strconv.FormatUint(msg.StartBlock, 10)),
			sdk.NewAttribute(types.AttributeKeyEndBlock, strconv.FormatUint(msg.EndBlock, 10)),
			sdk.NewAttribute(types.AttributeKeyRootHash, msg.RootHash.String()),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}
