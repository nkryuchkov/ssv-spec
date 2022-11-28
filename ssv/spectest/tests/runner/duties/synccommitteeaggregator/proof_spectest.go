package synccommitteeaggregator

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

type SyncCommitteeAggregatorProofSpecTest struct {
	Name                    string
	Messages                []*types.SSVMessage
	PostDutyRunnerStateRoot string
	ProofRootsMap           map[string]bool // if true then root returned from beacon node will be an aggregator
	ExpectedError           string
}

func (test *SyncCommitteeAggregatorProofSpecTest) TestName() string {
	return test.Name
}

func (test *SyncCommitteeAggregatorProofSpecTest) Run(t *testing.T) {
	ks := testingutils.Testing4SharesSet()
	share := testingutils.TestingShare(ks, testingutils.TestingProposer(ks, qbft.FirstHeight, qbft.FirstRound))
	v := testingutils.BaseValidator(keySetForShare(share))
	r := v.DutyRunners[types.BNRoleSyncCommitteeContribution]
	r.GetBeaconNode().(*testingutils.TestingBeaconNode).SetSyncCommitteeAggregatorRootHexes(test.ProofRootsMap)
	v.Beacon = r.GetBeaconNode()

	lastErr := v.StartDuty(testingutils.TestingSyncCommitteeContributionDuty)
	for _, msg := range test.Messages {
		err := v.ProcessMessage(msg)
		if err != nil {
			lastErr = err
		}
	}

	if len(test.ExpectedError) != 0 {
		require.EqualError(t, lastErr, test.ExpectedError)
	} else {
		require.NoError(t, lastErr)
	}

	// post root
	postRoot, err := r.GetBaseRunner().State.GetRoot()
	require.NoError(t, err)
	require.EqualValues(t, test.PostDutyRunnerStateRoot, hex.EncodeToString(postRoot))
}

func keySetForShare(share *types.Share) *testingutils.TestKeySet {
	if share.Quorum == 5 {
		return testingutils.Testing7SharesSet()
	}
	if share.Quorum == 7 {
		return testingutils.Testing10SharesSet()
	}
	if share.Quorum == 9 {
		return testingutils.Testing13SharesSet()
	}
	return testingutils.Testing4SharesSet()
}
