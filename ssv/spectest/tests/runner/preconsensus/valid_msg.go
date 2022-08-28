package preconsensus

import (
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// ValidMessage tests a valid SignedPartialSignatureMessage with multi PartialSignatureMessages
func ValidMessage() *tests.MultiMsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()
	return &tests.MultiMsgProcessingSpecTest{
		Name: "pre consensus valid msg",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name:   "sync committee aggregator selection proof",
				Runner: testingutils.SyncCommitteeContributionRunner(ks),
				Duty:   testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommitteeContribution(nil, testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1)),
				},
				PostDutyRunnerStateRoot: "0b955fd092aa319ac5c9e3a9fa0e405d9eebe81c19714c3ccc56945652777ede",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusContributionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:   "aggregator selection proof",
				Runner: testingutils.AggregatorRunner(ks),
				Duty:   testingutils.TestingAggregatorDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAggregator(nil, testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1)),
				},
				PostDutyRunnerStateRoot: "07875a38ed47e6f04c5f0820f4f684dc56e8686622f1906fb3577f245b57984f",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusSelectionProofMsg(ks.Shares[1], ks.Shares[1], 1, 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:   "randao",
				Runner: testingutils.ProposerRunner(ks),
				Duty:   testingutils.TestingProposerDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgProposer(nil, testingutils.PreConsensusRandaoDifferentSignerMsg(ks.Shares[1], ks.Shares[1], 1, 1)),
				},
				PostDutyRunnerStateRoot: "3f55ac00cd2b200bf7ad8784a38f60bc9f1389d110724a35bfed64aa29cde9c6",
				OutputMessages: []*ssv.SignedPartialSignatureMessage{
					testingutils.PreConsensusRandaoMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
			},
			{
				Name:   "attester",
				Runner: testingutils.AttesterRunner(ks),
				Duty:   testingutils.TestingAttesterDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAttester(nil, testingutils.PreConsensusFailedMsg(ks.Shares[1], 1)),
				},
				PostDutyRunnerStateRoot: "e80945f69c9142e391550f605d18a1e39f423d5f865ef865421f07ec25c8b803",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				ExpectedError:           "no pre consensus sigs required for attester role",
			},
			{
				Name:   "sync committee",
				Runner: testingutils.SyncCommitteeRunner(ks),
				Duty:   testingutils.TestingSyncCommitteeDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommittee(nil, testingutils.PreConsensusFailedMsg(ks.Shares[1], 1)),
				},
				PostDutyRunnerStateRoot: "5ff863a832e766731418768336ccb6580bb3dd33c2a0436d0196ea6525a0749a",
				OutputMessages:          []*ssv.SignedPartialSignatureMessage{},
				ExpectedError:           "no pre consensus sigs required for sync committee role",
			},
		},
	}
}
