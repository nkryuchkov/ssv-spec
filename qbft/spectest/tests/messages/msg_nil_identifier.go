package messages

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// MsgNilIdentifier tests Message with Identifier == nil
func MsgNilIdentifier() *tests.MsgSpecTest {
	ks := testingutils.Testing4SharesSet()
	msg := testingutils.SignQBFTMsg(ks.Shares[1], types.OperatorID(1), &qbft.Message{
		MsgType:    qbft.CommitMsgType,
		Height:     qbft.FirstHeight,
		Round:      qbft.FirstRound,
		Identifier: nil,
		Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
	})

	return &tests.MsgSpecTest{
		Name: "msg identifier nil",
		Messages: []*qbft.SignedMessage{
			msg,
		},
		ExpectedError: "message identifier is invalid",
	}
}
