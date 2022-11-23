package processmsg

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// MsgError tests a process msg returning an error
func MsgError() *tests.ControllerSpecTest {
	identifier := types.NewMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	ks := testingutils.Testing4SharesSet()

	return &tests.ControllerSpecTest{
		Name: "process msg error",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue: []byte{1, 2, 3, 4},
				InputMessages: []*qbft.SignedMessage{
					testingutils.SignQBFTMsg(
						ks.Shares[testingutils.TestingProposer(ks, qbft.FirstHeight, 100)],
						testingutils.TestingProposer(ks, qbft.FirstHeight, 100),
						&qbft.Message{
							MsgType:    qbft.ProposalMsgType,
							Height:     qbft.FirstHeight,
							Round:      100,
							Identifier: identifier[:],
							Data:       testingutils.ProposalDataBytes([]byte{1, 2, 3, 4}, nil, nil),
						},
					),
				},
				ControllerPostRoot: "5b6ebc3aa0bfcedd466fca3fca7e1dcc0245def7d61d65aee1462436d819c7d0",
			},
		},
		ExpectedError: "could not process msg: proposal invalid: proposal not justified: change round has no quorum",
	}
}
