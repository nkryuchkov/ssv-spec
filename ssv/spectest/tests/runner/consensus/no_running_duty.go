package consensus

import (
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/ssv"
	"github.com/bloxapp/ssv-spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv-spec/types/testingutils"
)

// NoRunningDuty tests a valid proposal msg before duty starts
func NoRunningDuty() *tests.MultiMsgProcessingSpecTest {
	ks := testingutils.Testing4SharesSet()
	startInstance := func(r ssv.Runner, value []byte) ssv.Runner {
		r.GetBaseRunner().QBFTController.StoredInstances = append(r.GetBaseRunner().QBFTController.StoredInstances, qbft.NewInstance(
			r.GetBaseRunner().QBFTController.GetConfig(),
			r.GetBaseRunner().QBFTController.Share,
			r.GetBaseRunner().QBFTController.Identifier,
			qbft.FirstHeight))

		return r
	}

	return &tests.MultiMsgProcessingSpecTest{
		Name: "consensus no running duty",
		Tests: []*tests.MsgProcessingSpecTest{
			{
				Name: "sync committee contribution",
				Runner: startInstance(
					testingutils.SyncCommitteeContributionRunner(ks),
					testingutils.TestSyncCommitteeContributionConsensusDataByts,
				),
				Duty: &testingutils.TestingSyncCommitteeContributionDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommitteeContribution(
						testingutils.TestingProposalMessageWithIdentifierAndFullData(
							ks.Shares[1], types.OperatorID(1), testingutils.SyncCommitteeContributionMsgID,
							testingutils.TestSyncCommitteeContributionConsensusDataByts,
						),
						nil),
				},
				PostDutyRunnerStateRoot: "40ad0587a87848675408886eda4f180caece4406851f873f0887b4800ac75429",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
				ExpectedError:           "failed processing consensus message: no running duty",
			},
			{
				Name: "sync committee",
				Runner: startInstance(
					testingutils.SyncCommitteeRunner(ks),
					testingutils.TestSyncCommitteeConsensusDataByts,
				),
				Duty: &testingutils.TestingSyncCommitteeDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgSyncCommittee(
						testingutils.TestingProposalMessageWithIdentifierAndFullData(
							ks.Shares[1], types.OperatorID(1), testingutils.SyncCommitteeMsgID,
							testingutils.TestSyncCommitteeConsensusDataByts,
						),
						nil),
				},
				PostDutyRunnerStateRoot: "91b27703f146f5e2c057d4753afaf653e35a0296146afe1d7dc64e18cd27cd2e",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
				ExpectedError:           "failed processing consensus message: no running duty",
			},
			{
				Name: "aggregator",
				Runner: startInstance(
					testingutils.AggregatorRunner(ks),
					testingutils.TestAggregatorConsensusDataByts,
				),
				Duty: &testingutils.TestingAggregatorDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAggregator(
						testingutils.TestingProposalMessageWithIdentifierAndFullData(
							ks.Shares[1], types.OperatorID(1), testingutils.AggregatorMsgID,
							testingutils.TestAggregatorConsensusDataByts,
						),
						nil),
				},
				PostDutyRunnerStateRoot: "8b1322320222b9894052c6c17d621535a37ece0d5ae7f324096a23c4c50a3b02",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
				ExpectedError:           "failed processing consensus message: no running duty",
			},
			{
				Name: "proposer",
				Runner: startInstance(
					testingutils.ProposerRunner(ks),
					testingutils.TestProposerConsensusDataByts,
				),
				Duty: &testingutils.TestingProposerDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgProposer(
						testingutils.TestingProposalMessageWithIdentifierAndFullData(
							ks.Shares[1], types.OperatorID(1), testingutils.ProposerMsgID,
							testingutils.TestProposerConsensusDataByts,
						),
						nil),
				},
				PostDutyRunnerStateRoot: "47dc9617050ef7348be9d47bfcd950d49e59d59e66d25502df480a4e656415e4",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
				ExpectedError:           "failed processing consensus message: no running duty",
			},
			{
				Name: "proposer (blinded block)",
				Runner: startInstance(
					testingutils.ProposerBlindedBlockRunner(ks),
					testingutils.TestProposerBlindedBlockConsensusDataByts,
				),
				Duty: &testingutils.TestingProposerDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgProposer(
						testingutils.TestingProposalMessageWithIdentifierAndFullData(
							ks.Shares[1], types.OperatorID(1), testingutils.ProposerMsgID,
							testingutils.TestProposerBlindedBlockConsensusDataByts,
						),
						nil),
				},
				PostDutyRunnerStateRoot: "fad9325ccafb559588cbb9d51a2fb669e16dcc2d6ffa8688b5a5d54d8bac9a50",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
				ExpectedError:           "failed processing consensus message: no running duty",
			},
			{
				Name: "attester",
				Runner: startInstance(
					testingutils.AttesterRunner(ks),
					testingutils.TestAttesterConsensusDataByts,
				),
				Duty: &testingutils.TestingAttesterDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgAttester(
						testingutils.TestingProposalMessageWithIdentifierAndFullData(
							ks.Shares[1], types.OperatorID(1), testingutils.AttesterMsgID,
							testingutils.TestAttesterConsensusDataByts,
						),
						nil),
				},
				PostDutyRunnerStateRoot: "7508cb14dd0d5833d218141b74b6fd8c512a54e99143c7f91d754ea3f08f1bd6",
				OutputMessages:          []*types.SignedPartialSignatureMessage{},
				DontStartDuty:           true,
				ExpectedError:           "failed processing consensus message: no running duty",
			},
			{
				Name:   "validator registration",
				Runner: testingutils.ValidatorRegistrationRunner(ks),
				Duty:   &testingutils.TestingValidatorRegistrationDuty,
				Messages: []*types.SSVMessage{
					testingutils.SSVMsgValidatorRegistration(
						testingutils.TestingProposalMessageWithIdentifierAndFullData(
							ks.Shares[1], types.OperatorID(1), testingutils.ValidatorRegistrationMsgID,
							testingutils.TestAttesterConsensusDataByts,
						),
						nil),
				},
				PostDutyRunnerStateRoot: "c18dd94eb6d9e94f93746320cf7eaba8ceb4ede53b81d865f0a78d310e1adde0",
				OutputMessages: []*types.SignedPartialSignatureMessage{
					testingutils.PreConsensusValidatorRegistrationMsg(ks.Shares[1], 1), // broadcasts when starting a new duty
				},
				ExpectedError: "no consensus phase for validator registration",
			},
		},
	}
}
