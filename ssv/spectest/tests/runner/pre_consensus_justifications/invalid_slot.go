package pre_consensus_justifications

import (
	"github.com/attestantio/go-eth2-client/spec"

	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// InvalidSlot tests an invalid pre-consensus justification message slot
func InvalidSlot() tests.SpecTest {
	ks := testingutils.Testing4SharesSet()

	signedMsg := func(obj *types.ConsensusData, id []byte) *qbft.SignedMessage {
		fullData, _ := obj.Encode()
		root, _ := qbft.HashDataRoot(fullData)
		msg := &qbft.Message{
			MsgType:    qbft.ProposalMsgType,
			Height:     1,
			Round:      qbft.FirstRound,
			Identifier: id,
			Root:       root,
		}
		signed := testingutils.SignQBFTMsg(ks.Shares[1], 1, msg)
		signed.FullData = fullData

		return signed
	}

	msgF := func(obj *types.ConsensusData, id []byte) *qbft.SignedMessage {
		// change slot
		if len(obj.PreConsensusJustifications) > 0 {
			obj.PreConsensusJustifications[0].Message.Slot = testingutils.TestingDutySlot2
		}
		return signedMsg(obj, id)
	}

	msgFV := func(obj *types.ConsensusData, id []byte) *qbft.SignedMessage {
		// change slot
		if len(obj.PreConsensusJustifications) > 0 {
			obj.PreConsensusJustifications[0].Message.Slot = testingutils.TestingInvalidDutySlotV(obj.Version)
		}
		return signedMsg(obj, id)
	}

	expectedErr := "failed processing consensus message: invalid pre-consensus justification: invalid partial sig slot"

	return &tests.MultiMsgProcessingSpecTest{
		Name: "pre consensus invalid slot",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:   "sync committee aggregator selection proof",
				Runner: decideFirstHeight(testingutils.SyncCommitteeContributionRunner(ks)),
				Duty:   &testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommitteeContribution(msgF(testingutils.TestContributionProofWithJustificationsConsensusData(ks), testingutils.SyncCommitteeContributionMsgID), nil),
				},
				PostDutyRunnerStateRoot: "2e0d0c26372ecd5d3ff786bd28581b624fb2f30c84821b23bb15a61faf377d3e",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
				ExpectedError: expectedErr,
			},
			{
				Name:   "aggregator selection proof",
				Runner: decideFirstHeight(testingutils.AggregatorRunner(ks)),
				Duty:   &testingutils.TestingAggregatorDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAggregator(msgF(testingutils.TestSelectionProofWithJustificationsConsensusData(ks), testingutils.AggregatorMsgID), nil),
				},
				PostDutyRunnerStateRoot: "c2433300980af6fae62151b9fc6bda67842e019205b5a31660692ab8e99ddbb2",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
				ExpectedError: expectedErr,
			},
			{
				Name:   "randao",
				Runner: decideFirstHeight(testingutils.ProposerRunner(ks)),
				Duty:   testingutils.TestingProposerDutyV(spec.DataVersionBellatrix),
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgProposer(msgFV(testingutils.TestProposerWithJustificationsConsensusDataV(ks, spec.DataVersionBellatrix), testingutils.ProposerMsgID), nil),
				},
				PostDutyRunnerStateRoot: "3544989d49ef3071258fd72b2befa80366afafcb788914a17f6226ad957eb980",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, spec.DataVersionBellatrix), // broadcasts when starting a new duty
				},
				ExpectedError: expectedErr,
			},
			{
				Name:   "randao (blinded block)",
				Runner: decideFirstHeight(testingutils.ProposerBlindedBlockRunner(ks)),
				Duty:   testingutils.TestingProposerDutyV(spec.DataVersionBellatrix),
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgProposer(msgF(testingutils.TestProposerBlindedWithJustificationsConsensusDataV(ks, spec.DataVersionBellatrix), testingutils.ProposerMsgID), nil),
				},
				PostDutyRunnerStateRoot: "35fce5d0dc47c393efcdaf98168a6b00e931695c73f6cc47f90d7d6b706ca089",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsgV(ks.Shares[1], 1, spec.DataVersionBellatrix), // broadcasts when starting a new duty
				},
				ExpectedError: expectedErr,
			},
			{

				Name:   "attester",
				Runner: decideFirstHeight(testingutils.AttesterRunner(ks)),
				Duty:   &testingutils.TestingAttesterDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAttester(msgF(testingutils.TestAttesterConsensusData, testingutils.AttesterMsgID), nil),
				},
				PostDutyRunnerStateRoot: "863e4147a71bcffc479bce197ef78d675271aad4702fdae08162fe3147340205",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
			},
			{
				Name:   "sync committee",
				Runner: decideFirstHeight(testingutils.SyncCommitteeRunner(ks)),
				Duty:   &testingutils.TestingSyncCommitteeDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommittee(msgF(testingutils.TestSyncCommitteeConsensusData, testingutils.SyncCommitteeMsgID), nil),
				},
				PostDutyRunnerStateRoot: "840de1bdaa22e3bf373cdaa0f3a56871170c0ba04fd637b3640ce4ff348b3c46",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
			},
		},
	}
}
