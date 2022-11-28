package proposal

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// PreparedPreviouslyDuplicateRCQuorum tests a proposal for > 1 round, prepared previously with quorum of round change but 2 are duplicates
func PreparedPreviouslyDuplicateRCQuorum() *tests.MsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()
	pre := testingutils.BaseInstance(testingutils.TestingProposer(ks, qbft.FirstHeight, qbft.FirstRound))
	pre.State.Round = 2

	prepareMsgs := []*qbft.SignedMessage{
		testingutils.SignQBFTMsg(ks.Shares[1], types.OperatorID(1), &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.PrepareDataBytes([]byte{1, 2, 3, 4}),
		}),
		testingutils.SignQBFTMsg(ks.Shares[2], types.OperatorID(2), &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.PrepareDataBytes([]byte{1, 2, 3, 4}),
		}),
		testingutils.SignQBFTMsg(ks.Shares[3], types.OperatorID(3), &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.PrepareDataBytes([]byte{1, 2, 3, 4}),
		}),
	}
	rcMsgs := []*qbft.SignedMessage{
		testingutils.SignQBFTMsg(ks.Shares[1], types.OperatorID(1), &qbft.Message{
			MsgType:    qbft.RoundChangeMsgType,
			Height:     qbft.FirstHeight,
			Round:      2,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.RoundChangePreparedDataBytes([]byte{1, 2, 3, 4}, qbft.FirstRound, prepareMsgs),
		}),
		testingutils.SignQBFTMsg(ks.Shares[2], types.OperatorID(2), &qbft.Message{
			MsgType:    qbft.RoundChangeMsgType,
			Height:     qbft.FirstHeight,
			Round:      2,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.RoundChangePreparedDataBytes([]byte{1, 2, 3, 4}, qbft.FirstRound, prepareMsgs),
		}),
		testingutils.SignQBFTMsg(ks.Shares[2], types.OperatorID(2), &qbft.Message{
			MsgType:    qbft.RoundChangeMsgType,
			Height:     qbft.FirstHeight,
			Round:      2,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.RoundChangePreparedDataBytes([]byte{1, 2, 3, 4}, qbft.FirstRound, prepareMsgs),
		}),
		testingutils.SignQBFTMsg(ks.Shares[3], types.OperatorID(3), &qbft.Message{
			MsgType:    qbft.RoundChangeMsgType,
			Height:     qbft.FirstHeight,
			Round:      2,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.RoundChangePreparedDataBytes([]byte{1, 2, 3, 4}, qbft.FirstRound, prepareMsgs),
		}),
	}

	msgs := []*qbft.SignedMessage{
		testingutils.SignQBFTMsg(
			ks.Shares[testingutils.TestingProposer(ks, qbft.FirstHeight, 2)],
			testingutils.TestingProposer(ks, qbft.FirstHeight, 2),
			&qbft.Message{
				MsgType:    qbft.ProposalMsgType,
				Height:     qbft.FirstHeight,
				Round:      2,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.ProposalDataBytes([]byte{1, 2, 3, 4}, rcMsgs, prepareMsgs),
			},
		),
	}
	return &tests.MsgProcessingSpecTest{
		Name:     "duplicate rc msg justification (prepared) quorum",
		Pre:      pre,
		PostRoot: "7555dc169779191419ccbe8833d5eac533d7a56d502f27faa819f935460b45b3",
		InputMessages: append(msgs, testingutils.SignQBFTMsg(ks.Shares[1], types.OperatorID(1), &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      2,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.PrepareDataBytes([]byte{1, 2, 3, 4}),
		})),
		OutputMessages: []*qbft.SignedMessage{
			testingutils.SignQBFTMsg(ks.Shares[1], types.OperatorID(1), &qbft.Message{
				MsgType:    qbft.PrepareMsgType,
				Height:     qbft.FirstHeight,
				Round:      2,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.PrepareDataBytes([]byte{1, 2, 3, 4}),
			}),
		},
	}
}
