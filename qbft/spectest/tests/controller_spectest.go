package tests

import (
	"bytes"
	"encoding/hex"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

type RunInstanceData struct {
	InputValue         []byte
	InputMessages      []*qbft.SignedMessage
	DecidedVal         []byte
	DecidedCnt         uint
	SavedDecided       *qbft.SignedMessage
	BroadcastedDecided *qbft.SignedMessage
	ControllerPostRoot string
	ExpectedTimerState *testingutils.TimerState
}

type ControllerSpecTest struct {
	Name            string
	RunInstanceData []*RunInstanceData
	OutputMessages  []*qbft.SignedMessage
	ExpectedError   string
}

func (test *ControllerSpecTest) Run(t *testing.T) {
	ks := testingutils.Testing4SharesSet()
	identifier := types.NewMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	config := testingutils.TestingConfig(ks)
	contr := testingutils.NewTestingQBFTController(
		identifier[:],
		testingutils.TestingShare(ks, testingutils.TestingProposer(ks, qbft.FirstHeight, qbft.FirstRound)),
		config,
	)

	var lastErr error
	for _, runData := range test.RunInstanceData {
		err := contr.StartNewInstance(runData.InputValue)
		if err != nil {
			lastErr = err
		} else if runData.ExpectedTimerState != nil {
			if timer, ok := config.GetTimer().(*testingutils.TestQBFTTimer); ok {
				require.Equal(t, runData.ExpectedTimerState.Timeouts, timer.State.Timeouts)
				require.Equal(t, runData.ExpectedTimerState.Round, timer.State.Round)
			}
		}

		decidedCnt := 0
		for _, msg := range runData.InputMessages {
			decided, err := contr.ProcessMsg(msg)
			if err != nil {
				lastErr = err
			}
			if decided != nil {
				decidedCnt++

				data, _ := decided.Message.GetCommitData()
				require.EqualValues(t, runData.DecidedVal, data.Data)
			}
		}

		require.EqualValues(t, runData.DecidedCnt, decidedCnt)

		if runData.SavedDecided != nil {
			// test saved to storage
			decided, err := config.GetStorage().GetHighestDecided(identifier[:])
			require.NoError(t, err)
			require.NotNil(t, decided)
			r1, err := decided.GetRoot()
			require.NoError(t, err)

			r2, err := runData.SavedDecided.GetRoot()
			require.NoError(t, err)

			require.EqualValues(t, r2, r1)
			require.EqualValues(t, runData.SavedDecided.Signers, decided.Signers)
			require.EqualValues(t, runData.SavedDecided.Signature, decided.Signature)
		}
		if runData.BroadcastedDecided != nil {
			// test broadcasted
			broadcastedMsgs := config.GetNetwork().(*testingutils.TestingNetwork).BroadcastedMsgs
			require.Greater(t, len(broadcastedMsgs), 0)
			found := false
			for _, msg := range broadcastedMsgs {
				if !bytes.Equal(identifier[:], msg.MsgID[:]) {
					continue
				}

				msg1 := &qbft.SignedMessage{}
				require.NoError(t, msg1.Decode(msg.Data))
				r1, err := msg1.GetRoot()
				require.NoError(t, err)

				r2, err := runData.BroadcastedDecided.GetRoot()
				require.NoError(t, err)

				if bytes.Equal(r1, r2) &&
					reflect.DeepEqual(runData.BroadcastedDecided.Signers, msg1.Signers) &&
					reflect.DeepEqual(runData.BroadcastedDecided.Signature, msg1.Signature) {
					require.False(t, found)
					found = true
				}
			}
			require.True(t, found)
		}

		r, err := contr.GetRoot()
		require.NoError(t, err)
		require.EqualValues(t, runData.ControllerPostRoot, hex.EncodeToString(r))
	}

	if len(test.ExpectedError) != 0 {
		require.EqualError(t, lastErr, test.ExpectedError)
	} else {
		require.NoError(t, lastErr)
	}
}

func (test *ControllerSpecTest) TestName() string {
	return "qbft controller " + test.Name
}
