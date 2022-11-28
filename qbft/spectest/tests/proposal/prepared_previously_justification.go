package proposal

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// PreparedPreviouslyJustification tests a proposal for > 1 round, prepared previously with quorum of round change msgs justification
func PreparedPreviouslyJustification() *tests.MsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()
	pre := testingutils.BaseInstance(testingutils.TestingProposer(ks, qbft.FirstHeight, 2))

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
			ks.Shares[testingutils.TestingProposer(ks, qbft.FirstHeight, qbft.FirstRound)],
			testingutils.TestingProposer(ks, qbft.FirstHeight, qbft.FirstRound),
			&qbft.Message{
				MsgType:    qbft.ProposalMsgType,
				Height:     qbft.FirstHeight,
				Round:      qbft.FirstRound,
				Identifier: []byte{1, 2, 3, 4},
				Data:       testingutils.ProposalDataBytes([]byte{1, 2, 3, 4}, nil, nil),
			},
		),
	}
	msgs = append(msgs, prepareMsgs...)
	msgs = append(msgs, rcMsgs...)
	msgs = append(msgs,
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
	)
	return &tests.MsgProcessingSpecTest{
		Name:          "previously prepared proposal",
		Pre:           pre,
		PostRoot:      "de8ead2908813fa1b14e1b704f060ae730be38ccf3c08b37094a2b0a8b1be65f",
		InputMessages: msgs,
		OutputMessages: []*qbft.SignedMessage{
			testingutils.SignQBFTMsg(
				ks.Shares[testingutils.TestingProposer(ks, qbft.FirstHeight, qbft.FirstRound)],
				testingutils.TestingProposer(ks, qbft.FirstHeight, qbft.FirstRound),
				&qbft.Message{
					MsgType:    qbft.PrepareMsgType,
					Height:     qbft.FirstHeight,
					Round:      qbft.FirstRound,
					Identifier: []byte{1, 2, 3, 4},
					Data:       testingutils.PrepareDataBytes([]byte{1, 2, 3, 4}),
				},
			),
			testingutils.SignQBFTMsg(
				ks.Shares[testingutils.TestingProposer(ks, qbft.FirstHeight, qbft.FirstRound)],
				testingutils.TestingProposer(ks, qbft.FirstHeight, qbft.FirstRound),
				&qbft.Message{
					MsgType:    qbft.CommitMsgType,
					Height:     qbft.FirstHeight,
					Round:      qbft.FirstRound,
					Identifier: []byte{1, 2, 3, 4},
					Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
				},
			),
			testingutils.SignQBFTMsg(
				ks.Shares[testingutils.TestingProposer(ks, qbft.FirstHeight, 2)],
				testingutils.TestingProposer(ks, qbft.FirstHeight, 2),
				&qbft.Message{
					MsgType:    qbft.RoundChangeMsgType,
					Height:     qbft.FirstHeight,
					Round:      2,
					Identifier: []byte{1, 2, 3, 4},
					Data:       testingutils.RoundChangePreparedDataBytes([]byte{1, 2, 3, 4}, qbft.FirstRound, prepareMsgs),
				},
			),
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
			testingutils.SignQBFTMsg(
				ks.Shares[testingutils.TestingProposer(ks, qbft.FirstHeight, 2)],
				testingutils.TestingProposer(ks, qbft.FirstHeight, 2),
				&qbft.Message{
					MsgType:    qbft.PrepareMsgType,
					Height:     qbft.FirstHeight,
					Round:      2,
					Identifier: []byte{1, 2, 3, 4},
					Data:       testingutils.PrepareDataBytes([]byte{1, 2, 3, 4}),
				},
			),
		},
	}
}
