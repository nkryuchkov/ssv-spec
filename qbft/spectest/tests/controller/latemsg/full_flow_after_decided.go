package latemsg

import (
	"github.com/herumi/bls-eth-go-binary/bls"

	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// FullFlowAfterDecided tests a decided msg for round 1 followed by a full proposal, prepare, commit for round 2
func FullFlowAfterDecided() *tests.ControllerSpecTest {
	identifier := types.NewMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	ks := testingutils.Testing4SharesSet()

	rcMsgs := []*qbft.SignedMessage{
		testingutils.SignQBFTMsg(ks.Shares[1], types.OperatorID(1), &qbft.Message{
			MsgType:    qbft.RoundChangeMsgType,
			Height:     qbft.FirstHeight,
			Round:      2,
			Identifier: identifier[:],
			Data:       testingutils.RoundChangeDataBytes(nil, qbft.NoRound),
		}),
		testingutils.SignQBFTMsg(ks.Shares[2], types.OperatorID(2), &qbft.Message{
			MsgType:    qbft.RoundChangeMsgType,
			Height:     qbft.FirstHeight,
			Round:      2,
			Identifier: identifier[:],
			Data:       testingutils.RoundChangeDataBytes(nil, qbft.NoRound),
		}),
		testingutils.SignQBFTMsg(ks.Shares[3], types.OperatorID(3), &qbft.Message{
			MsgType:    qbft.RoundChangeMsgType,
			Height:     qbft.FirstHeight,
			Round:      2,
			Identifier: identifier[:],
			Data:       testingutils.RoundChangeDataBytes(nil, qbft.NoRound),
		}),
	}

	msgs := []*qbft.SignedMessage{
		testingutils.MultiSignQBFTMsg(
			[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
			[]types.OperatorID{1, 2, 3},
			&qbft.Message{
				MsgType:    qbft.CommitMsgType,
				Height:     qbft.FirstHeight,
				Round:      qbft.FirstRound,
				Identifier: identifier[:],
				Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
			}),

		testingutils.SignQBFTMsg(
			ks.Shares[testingutils.TestingProposer(ks, qbft.FirstHeight, 2)],
			testingutils.TestingProposer(ks, qbft.FirstHeight, 2),
			&qbft.Message{
				MsgType:    qbft.ProposalMsgType,
				Height:     qbft.FirstHeight,
				Round:      2,
				Identifier: identifier[:],
				Data:       testingutils.ProposalDataBytes([]byte{1, 2, 3, 4}, rcMsgs, nil),
			},
		),

		testingutils.SignQBFTMsg(ks.Shares[1], types.OperatorID(1), &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      2,
			Identifier: identifier[:],
			Data:       testingutils.PrepareDataBytes([]byte{1, 2, 3, 4}),
		}),
		testingutils.SignQBFTMsg(ks.Shares[2], types.OperatorID(2), &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      2,
			Identifier: identifier[:],
			Data:       testingutils.PrepareDataBytes([]byte{1, 2, 3, 4}),
		}),
		testingutils.SignQBFTMsg(ks.Shares[3], types.OperatorID(3), &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      2,
			Identifier: identifier[:],
			Data:       testingutils.PrepareDataBytes([]byte{1, 2, 3, 4}),
		}),

		testingutils.SignQBFTMsg(ks.Shares[1], types.OperatorID(1), &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      2,
			Identifier: identifier[:],
			Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
		}),
		testingutils.SignQBFTMsg(ks.Shares[2], types.OperatorID(2), &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      2,
			Identifier: identifier[:],
			Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
		}),
		testingutils.SignQBFTMsg(ks.Shares[4], types.OperatorID(4), &qbft.Message{
			MsgType:    qbft.CommitMsgType,
			Height:     qbft.FirstHeight,
			Round:      2,
			Identifier: identifier[:],
			Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
		}),
	}

	return &tests.ControllerSpecTest{
		Name: "full flow after decided",
		RunInstanceData: []*tests.RunInstanceData{
			{
				InputValue:    []byte{1, 2, 3, 4},
				InputMessages: msgs,
				SavedDecided: testingutils.MultiSignQBFTMsg(
					[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
					[]types.OperatorID{1, 2, 3},
					&qbft.Message{
						MsgType:    qbft.CommitMsgType,
						Height:     qbft.FirstHeight,
						Round:      qbft.FirstRound,
						Identifier: identifier[:],
						Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
					}),
				DecidedVal:         []byte{1, 2, 3, 4},
				DecidedCnt:         1,
				ControllerPostRoot: "dd2d0705ee3fc49644aa4cc264a095376c9c36802663e6e8e778851c2871bbce",
			},
		},
	}
}
