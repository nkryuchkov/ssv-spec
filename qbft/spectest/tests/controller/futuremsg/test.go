package futuremsg

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

type ControllerSyncSpecTest struct {
	Name                 string
	InputMessages        []*qbft.SignedMessage
	SyncDecidedCalledCnt int
	ControllerPostRoot   string
	ExpectedError        string
}

func (test *ControllerSyncSpecTest) TestName() string {
	return "qbft controller sync " + test.Name
}

func (test *ControllerSyncSpecTest) Run(t *testing.T) {
	ks := testingutils.Testing4SharesSet()
	identifier := types.NewMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	config := testingutils.TestingConfig(ks)
	contr := testingutils.NewTestingQBFTController(
		identifier[:],
		testingutils.TestingShare(ks, testingutils.TestingProposer(ks, qbft.FirstHeight, qbft.FirstRound)),
		config,
	)

	err := contr.StartNewInstance([]byte{1, 2, 3, 4})
	if err != nil {
		t.Fatalf(err.Error())
	}

	var lastErr error
	for _, msg := range test.InputMessages {
		_, err := contr.ProcessMsg(msg)
		if err != nil {
			lastErr = err
		}
	}

	syncedDecidedCnt := config.GetNetwork().(*testingutils.TestingNetwork).SyncHighestDecidedCnt
	require.EqualValues(t, test.SyncDecidedCalledCnt, syncedDecidedCnt)

	r, err := contr.GetRoot()
	require.NoError(t, err)
	require.EqualValues(t, test.ControllerPostRoot, hex.EncodeToString(r))

	if len(test.ExpectedError) != 0 {
		require.EqualError(t, lastErr, test.ExpectedError)
	} else {
		require.NoError(t, lastErr)
	}
}
