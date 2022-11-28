package roundchange

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// QuorumNotPrepared tests a round change quorum for non-prepared state
func QuorumNotPrepared() *tests.MsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()
	pre := testingutils.BaseInstance(testingutils.TestingProposer(ks, qbft.FirstHeight, 2))
	pre.State.Round = 2

	msgs := []*qbft.SignedMessage{
		testingutils.SignQBFTMsg(ks.Shares[1], types.OperatorID(1), &qbft.Message{
			MsgType:    qbft.RoundChangeMsgType,
			Height:     qbft.FirstHeight,
			Round:      2,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.RoundChangeDataBytes(nil, qbft.NoRound),
		}),
		testingutils.SignQBFTMsg(ks.Shares[2], types.OperatorID(2), &qbft.Message{
			MsgType:    qbft.RoundChangeMsgType,
			Height:     qbft.FirstHeight,
			Round:      2,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.RoundChangeDataBytes(nil, qbft.NoRound),
		}),
		testingutils.SignQBFTMsg(ks.Shares[3], types.OperatorID(3), &qbft.Message{
			MsgType:    qbft.RoundChangeMsgType,
			Height:     qbft.FirstHeight,
			Round:      2,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.RoundChangeDataBytes(nil, qbft.NoRound),
		}),
	}

	return &tests.MsgProcessingSpecTest{
		Name:          "round change not prepared",
		Pre:           pre,
		PostRoot:      "692615005a5c2ee1a292cea199abb03a9b03c42c9b9dec67ace830718fba6181",
		InputMessages: msgs,
		OutputMessages: []*qbft.SignedMessage{
			testingutils.SignQBFTMsg(
				ks.Shares[testingutils.TestingProposer(ks, qbft.FirstHeight, 2)],
				testingutils.TestingProposer(ks, qbft.FirstHeight, 2),
				&qbft.Message{
					MsgType:    qbft.ProposalMsgType,
					Height:     qbft.FirstHeight,
					Round:      2,
					Identifier: []byte{1, 2, 3, 4},
					Data:       testingutils.ProposalDataBytes([]byte{1, 2, 3, 4}, msgs, nil),
				},
			),
		},
	}
}
