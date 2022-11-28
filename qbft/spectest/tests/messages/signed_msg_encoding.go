package messages

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// SignedMessageEncoding tests encoding SignedMessage
func SignedMessageEncoding() *tests.MsgSpecTest {
	ks := testingutils.Testing4SharesSet()
	msg := testingutils.SignQBFTMsg(ks.Shares[1], types.OperatorID(1), &qbft.Message{
		MsgType:    qbft.ProposalMsgType,
		Height:     qbft.FirstHeight,
		Round:      qbft.FirstRound,
		Identifier: []byte{1, 2, 3, 4},
		Data: testingutils.ProposalDataBytes(
			[]byte{1, 2, 3, 4},
			[]*qbft.SignedMessage{
				testingutils.SignQBFTMsg(ks.Shares[1], types.OperatorID(1), &qbft.Message{
					MsgType:    qbft.PrepareMsgType,
					Height:     qbft.FirstHeight,
					Round:      qbft.FirstRound,
					Identifier: []byte{1, 2, 3, 4},
					Data:       testingutils.PrepareDataBytes([]byte{1, 2, 3, 4}),
				}),
			},
			[]*qbft.SignedMessage{
				testingutils.SignQBFTMsg(ks.Shares[1], types.OperatorID(1), &qbft.Message{
					MsgType:    qbft.RoundChangeMsgType,
					Height:     qbft.FirstHeight,
					Round:      qbft.FirstRound,
					Identifier: []byte{1, 2, 3, 4},
					Data:       testingutils.PrepareDataBytes([]byte{1, 2, 3, 4}),
				}),
			},
		),
	})

	b, _ := msg.Encode()

	return &tests.MsgSpecTest{
		Name: "signed message encoding",
		Messages: []*qbft.SignedMessage{
			msg,
		},
		EncodedMessages: [][]byte{
			b,
		},
	}
}
