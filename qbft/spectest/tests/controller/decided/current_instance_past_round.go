package decided

import (
	"github.com/herumi/bls-eth-go-binary/bls"

	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// CurrentInstancePastRound tests a decided msg received for current running instance for a past round
func CurrentInstancePastRound() *tests.ControllerSpecTest {
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
		testingutils.SignQBFTMsg(
			testingutils.Testing4SharesSet().Shares[testingutils.TestingProposer(ks, qbft.FirstHeight, qbft.FirstRound)],
			testingutils.TestingProposer(ks, qbft.FirstHeight, qbft.FirstRound),
			&qbft.Message{
				MsgType:    qbft.ProposalMsgType,
				Height:     qbft.FirstHeight,
				Round:      qbft.FirstRound,
				Identifier: identifier[:],
				Data:       testingutils.ProposalDataBytes([]byte{1, 2, 3, 4}, nil, nil),
			},
		),

		testingutils.SignQBFTMsg(ks.Shares[1], types.OperatorID(1), &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: identifier[:],
			Data:       testingutils.PrepareDataBytes([]byte{1, 2, 3, 4}),
		}),
		testingutils.SignQBFTMsg(ks.Shares[2], types.OperatorID(2), &qbft.Message{
			MsgType:    qbft.PrepareMsgType,
			Height:     qbft.FirstHeight,
			Round:      qbft.FirstRound,
			Identifier: identifier[:],
			Data:       testingutils.PrepareDataBytes([]byte{1, 2, 3, 4}),
		}),
	}
	msgs = append(msgs, rcMsgs...)
	msgs = append(msgs, []*qbft.SignedMessage{
		testingutils.SignQBFTMsg(
			testingutils.Testing4SharesSet().Shares[testingutils.TestingProposer(ks, qbft.FirstHeight, 2)],
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
	}...)

	return &tests.ControllerSpecTest{
		Name: "decide current instance past round",
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
				ControllerPostRoot: "ad6be16e4b01f66fde56f7d04f733dfb5b131ac8f216168215f35fb109c28820",
			},
		},
	}
}
