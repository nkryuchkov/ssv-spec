package ssv

import (
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/types"
	ssz "github.com/ferranbt/fastssz"
	"github.com/pkg/errors"
)

func (b *BaseRunner) signBeaconObject(
	runner Runner,
	obj ssz.HashRoot,
	slot spec.Slot,
	domainType spec.DomainType,
) (*types.PartialSignatureMessage, error) {
	epoch := runner.GetBaseRunner().BeaconNetwork.EstimatedEpochAtSlot(slot)
	domain, err := runner.GetBeaconNode().DomainData(epoch, domainType)
	if err != nil {
		return nil, errors.Wrap(err, "could not get beacon domain")
	}

	sig, r, err := runner.GetSigner().SignBeaconObject(obj, domain, runner.GetBaseRunner().Share.SharePubKey, domainType)
	if err != nil {
		return nil, errors.Wrap(err, "could not sign beacon object")
	}

	return &types.PartialSignatureMessage{
		PartialSignature: sig,
		SigningRoot:      r,
		Signer:           runner.GetBaseRunner().Share.OperatorID,
	}, nil
}

func (b *BaseRunner) signPostConsensusMsg(runner Runner, msg *types.PartialSignatureMessages) (*types.SignedPartialSignatureMessage, error) {
	signature, err := runner.GetSigner().SignRoot(msg, types.PartialSignatureType, b.Share.SharePubKey)
	if err != nil {
		return nil, errors.Wrap(err, "could not sign PartialSignatureMessage for PostConsensusContainer")
	}

	return &types.SignedPartialSignatureMessage{
		Message:   *msg,
		Signature: signature,
		Signer:    b.Share.OperatorID,
	}, nil
}

func (b *BaseRunner) validatePartialSigMsgForSlot(
	signedMsg *types.SignedPartialSignatureMessage,
	slot spec.Slot,
) error {
	if err := signedMsg.Validate(); err != nil {
		return errors.Wrap(err, "SignedPartialSignatureMessage invalid")
	}

	if signedMsg.Message.Slot != slot {
		return errors.New("invalid partial sig slot")
	}

	if err := signedMsg.GetSignature().VerifyByOperators(signedMsg, b.Share.DomainType, types.PartialSignatureType, b.Share.Committee); err != nil {
		return errors.Wrap(err, "failed to verify PartialSignature")
	}

	for _, msg := range signedMsg.Message.Messages {
		if err := b.verifyBeaconPartialSignature(msg); err != nil {
			return errors.Wrap(err, "could not verify Beacon partial Signature")
		}
	}

	return nil
}

func (b *BaseRunner) verifyBeaconPartialSignature(msg *types.PartialSignatureMessage) error {
	signer := msg.Signer
	signature := msg.PartialSignature
	root := msg.SigningRoot

	for _, n := range b.Share.Committee {
		if n.GetID() == signer {
			pk := new(types.BLSPublicKey).Uncompress(n.GetPublicKey())
			if pk == nil {
				return errors.New("could not deserialized pk")
			}

			sig := new(types.BLSSignature).Uncompress(signature)
			if sig == nil {
				return errors.New("could not deserialized Signature")
			}

			if !sig.Verify(false, pk, false, root[:], types.CipherSuite) {
				return errors.New("wrong signature")
			}
			return nil
		}
	}
	return errors.New("unknown signer")
}
