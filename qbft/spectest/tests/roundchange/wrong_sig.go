package roundchange

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// WrongSig tests a signed round change msg with wrong signature
func WrongSig() *tests.MsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()
	pre := testingutils.BaseInstance(testingutils.TestingProposer(ks, qbft.FirstHeight, qbft.FirstRound))
	pre.State.Round = 2

	msgs := []*qbft.SignedMessage{
		testingutils.SignQBFTMsg(ks.Shares[1], types.OperatorID(2), &qbft.Message{
			MsgType:    qbft.RoundChangeMsgType,
			Height:     qbft.FirstHeight,
			Round:      2,
			Identifier: []byte{1, 2, 3, 4},
			Data:       testingutils.RoundChangeDataBytes(nil, qbft.NoRound),
		}),
	}

	return &tests.MsgProcessingSpecTest{
		Name:           "round change invalid sig",
		Pre:            pre,
		PostRoot:       "4aafcc4aa9e2435579c85aa26e659fe650aefb8becb5738d32dd9286f7ff27c3",
		InputMessages:  msgs,
		OutputMessages: []*qbft.SignedMessage{},
		ExpectedError:  "round change msg invalid: round change msg signature invalid: failed to verify signature",
	}
}
