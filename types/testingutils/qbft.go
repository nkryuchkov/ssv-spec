package testingutils

import (
	"bytes"

	"github.com/pkg/errors"

	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/bloxapp/ssv-spec/types"
)

var TestingConfig = func(keySet *TestKeySet) *qbft.Config {
	return &qbft.Config{
		Signer:    NewTestingKeyManager(),
		SigningPK: keySet.Shares[1].GetPublicKey().Serialize(),
		Domain:    types.PrimusTestnet,
		ValueCheckF: func(data []byte) error {
			if bytes.Equal(data, TestingInvalidValueCheck) {
				return errors.New("invalid value")
			}

			// as a base validation we do not accept nil values
			if len(data) == 0 {
				return errors.New("invalid value")
			}
			return nil
		},
		ProposerF: qbft.RoundRobinProposer,
		Storage:   NewTestingStorage(),
		Network:   NewTestingNetwork(),
		Timer:     NewTestingTimer(),
	}
}

var TestingInvalidValueCheck = []byte{1, 1, 1, 1}

var TestingShare = func(keysSet *TestKeySet, operatorID types.OperatorID) *types.Share {
	return &types.Share{
		OperatorID:      operatorID,
		ValidatorPubKey: keysSet.ValidatorPK.Serialize(),
		SharePubKey:     keysSet.Shares[operatorID].GetPublicKey().Serialize(),
		DomainType:      types.PrimusTestnet,
		Quorum:          keysSet.Threshold,
		PartialQuorum:   keysSet.PartialThreshold,
		Committee:       keysSet.Committee(),
	}
}

func TestingProposer(keySet *TestKeySet, height qbft.Height, round qbft.Round) types.OperatorID {
	return TestingConfig(keySet).ProposerF(&qbft.State{Share: &types.Share{Committee: keySet.Committee()}, Height: height}, round)
}

var BaseInstance = func(proposerID types.OperatorID) *qbft.Instance {
	return baseInstance(TestingShare(Testing4SharesSet(), proposerID), Testing4SharesSet(), []byte{1, 2, 3, 4})
}

var SevenOperatorsInstance = func(proposerID types.OperatorID) *qbft.Instance {
	return baseInstance(TestingShare(Testing7SharesSet(), proposerID), Testing7SharesSet(), []byte{1, 2, 3, 4})
}

var TenOperatorsInstance = func(proposerID types.OperatorID) *qbft.Instance {
	return baseInstance(TestingShare(Testing10SharesSet(), proposerID), Testing10SharesSet(), []byte{1, 2, 3, 4})
}

var ThirteenOperatorsInstance = func(proposerID types.OperatorID) *qbft.Instance {
	return baseInstance(TestingShare(Testing13SharesSet(), proposerID), Testing13SharesSet(), []byte{1, 2, 3, 4})
}

var baseInstance = func(share *types.Share, keySet *TestKeySet, identifier []byte) *qbft.Instance {
	ret := qbft.NewInstance(TestingConfig(keySet), share, identifier, qbft.FirstHeight)
	ret.StartValue = []byte{1, 2, 3, 4}
	return ret
}

func NewTestingQBFTController(
	identifier []byte,
	share *types.Share,
	config qbft.IConfig,
) *qbft.Controller {
	return qbft.NewController(
		identifier,
		share,
		types.PrimusTestnet,
		config,
	)
}
