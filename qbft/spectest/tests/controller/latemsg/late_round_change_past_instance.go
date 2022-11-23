package latemsg

import (
	"github.com/herumi/bls-eth-go-binary/bls"

	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/qbft/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// LateProposalPastInstance tests process proposal msg for a previously decided instance
func LateProposalPastInstance() *tests.ControllerSpecTest {
	identifier := types.NewMsgID(testingutils.TestingValidatorPubKey[:], types.BNRoleAttester)
	ks := testingutils.Testing4SharesSet()

	allMsgs := testingutils.DecidingMsgsForHeight([]byte{1, 2, 3, 4}, identifier[:], 5, ks)
	msgPerHeight := make(map[qbft.Height][]*qbft.SignedMessage)
	msgPerHeight[qbft.FirstHeight] = allMsgs[0:7]
	msgPerHeight[1] = allMsgs[7:14]
	msgPerHeight[2] = allMsgs[14:21]
	msgPerHeight[3] = allMsgs[21:28]
	msgPerHeight[4] = allMsgs[28:35]
	msgPerHeight[5] = allMsgs[35:42]

	instanceData := func(height qbft.Height, postRoot string) *tests.RunInstanceData {
		return &tests.RunInstanceData{
			InputValue:    []byte{1, 2, 3, 4},
			InputMessages: msgPerHeight[height],
			SavedDecided: testingutils.MultiSignQBFTMsg(
				[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
				[]types.OperatorID{1, 2, 3},
				&qbft.Message{
					MsgType:    qbft.CommitMsgType,
					Height:     height,
					Round:      qbft.FirstRound,
					Identifier: identifier[:],
					Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
				}),
			BroadcastedDecided: testingutils.MultiSignQBFTMsg(
				[]*bls.SecretKey{ks.Shares[1], ks.Shares[2], ks.Shares[3]},
				[]types.OperatorID{1, 2, 3},
				&qbft.Message{
					MsgType:    qbft.CommitMsgType,
					Height:     height,
					Round:      qbft.FirstRound,
					Identifier: identifier[:],
					Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
				}),
			DecidedVal:         []byte{1, 2, 3, 4},
			DecidedCnt:         1,
			ControllerPostRoot: postRoot,
		}
	}

	return &tests.ControllerSpecTest{
		Name: "late proposal past instance",
		RunInstanceData: []*tests.RunInstanceData{
			instanceData(qbft.FirstHeight, "aa402d7487719b17dde352e2ac602ba2c7d895e615ab12cd93d816f6c4fa0967"),
			instanceData(1, "6eaf87922fcf15949d1a7e81c29528cf219dbb445dd3f323452b29bce86b8606"),
			instanceData(2, "4ab670bd4d1aedb4eaba76ea578880d7a3f5e585a4b443ff6bf27899d1d9433c"),
			instanceData(3, "8d91e5a72e416ef0b46cf7cadae265cc335467a7683698f120ae6b9d52afe3d5"),
			instanceData(4, "9fd86252efe33d7a4335ba4e17f6e957dc546cb071a6b1f147328c994f657738"),
			instanceData(5, "887b7a4af9cac0e27552768e7648ef9e8e198ec9c66c629026d0b404de9ced51"),
			{
				InputValue: []byte{1, 2, 3, 4},
				InputMessages: []*qbft.SignedMessage{
					testingutils.MultiSignQBFTMsg(
						[]*bls.SecretKey{ks.Shares[testingutils.TestingProposer(ks, 4, qbft.FirstRound)]},
						[]types.OperatorID{testingutils.TestingProposer(ks, 4, qbft.FirstRound)},
						&qbft.Message{
							MsgType:    qbft.ProposalMsgType,
							Height:     4,
							Round:      qbft.FirstRound,
							Identifier: identifier[:],
							Data:       testingutils.CommitDataBytes([]byte{1, 2, 3, 4}),
						}),
				},
				ControllerPostRoot: "0b9c323cf47c4a653a7f79f5b203a2c26d39865c59078579cb78e318e5a60849",
			},
		},
		ExpectedError: "could not process msg: proposal invalid: proposal is not valid with current state",
	}
}
