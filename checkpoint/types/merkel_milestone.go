package types

import (
	"errors"

	"github.com/maticnetwork/heimdall/helper"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

// ValidateMilestone - Validates if milestone rootHash matches or not
func ValidateMilestone(start uint64, end uint64, rootHash hmTypes.HeimdallHash, milestoneID string, contractCaller helper.IContractCaller, milestoneLength uint64) (bool, error) {

	if start+milestoneLength-1 != end {
		return false, errors.New("Invalid milestone, difference in start and end block is not equal to sprint length")
	}

	// Check if blocks exist locally
	if !contractCaller.CheckIfBlocksExist(end) {
		return false, errors.New("blocks not found locally")
	}

	return true, nil
	// Compare RootHash
	vote, err := contractCaller.GetVoteOnRootHash(start, end, milestoneLength, rootHash.String(), milestoneID)
	if err != nil {
		return false, err
	}

	return vote, nil
}
