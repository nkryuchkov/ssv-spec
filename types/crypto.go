package types

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/pkg/errors"
)

// VerifyByOperators verifies signature by the provided operators
func (s Signature) VerifyByOperators(data MessageSignature, domain DomainType, sigType SignatureType, operators []*Operator) error {
	// decode sig
	sign := new(BLSSignature).Uncompress(s)
	if sign == nil {
		return errors.New("failed to deserialize signature")
	}
	// find operators
	pks := make([]*BLSPublicKey, 0)
	for _, id := range data.GetSigners() {
		found := false
		for _, n := range operators {
			if id == n.GetID() {
				pk := new(BLSPublicKey).Uncompress(n.GetPublicKey())
				if pk == nil {
					return errors.New("failed to deserialize public key")
				}

				pks = append(pks, pk)
				found = true
			}
		}
		if !found {
			return errors.New("unknown signer")
		}
	}

	// compute root
	computedRoot, err := ComputeSigningRoot(data, ComputeSignatureDomain(domain, sigType))
	if err != nil {
		return errors.Wrap(err, "could not compute signing root")
	}

	// verify
	if res := sign.FastAggregateVerify(true, pks, computedRoot, CipherSuite); !res {
		return errors.New("failed to verify signature")
	}
	return nil
}

func (s Signature) VerifyMultiPubKey(data Root, domain DomainType, sigType SignatureType, pks [][]byte) error {
	aggPK := new(BLSAggregatePublicKey)

	for _, pkByts := range pks {
		pk := new(BLSPublicKey).Uncompress(pkByts)
		if pk == nil {
			return errors.New("failed to deserialize public key")
		}

		aggPK.Add(pk, false) // group already checked
	}

	if aggPK == nil {
		return errors.New("no public keys found")
	}

	return s.Verify(data, domain, sigType, aggPK.ToAffine().Compress())
}

func (s Signature) Verify(data Root, domain DomainType, sigType SignatureType, pkByts []byte) error {
	computedRoot, err := ComputeSigningRoot(data, ComputeSignatureDomain(domain, sigType))
	if err != nil {
		return errors.Wrap(err, "could not compute signing root")
	}

	sign := new(BLSSignature).Uncompress(s)
	if sign == nil {
		return errors.New("failed to deserialize signature")
	}

	pk := new(BLSPublicKey).Uncompress(pkByts)
	if pk == nil {
		return errors.New("failed to deserialize public key")
	}

	if res := sign.Verify(false, pk, false, computedRoot, CipherSuite); !res {
		return errors.New("failed to verify signature")
	}
	return nil
}

func (s Signature) ECRecover(data Root, domain DomainType, sigType SignatureType, address common.Address) error {
	computedRoot, err := ComputeSigningRoot(data, ComputeSignatureDomain(domain, sigType))
	if err != nil {
		return errors.Wrap(err, "could not compute signing root")
	}

	recoveredUncompressedPubKey, err := crypto.Ecrecover(computedRoot, s)
	if err != nil {
		return errors.Wrap(err, "could not recover ethereum address")
	}

	pk, err := secp256k1.ParsePubKey(recoveredUncompressedPubKey)
	if err != nil {
		return errors.Wrap(err, "could not parse ecdsa pubkey")
	}

	recoveredAdd := crypto.PubkeyToAddress(*pk.ToECDSA())

	if !bytes.Equal(address[:], recoveredAdd[:]) {
		return errors.Wrap(err, "message EC recover doesn't match address")
	}
	return nil
}

func (s Signature) Aggregate(other Signature) (Signature, error) {
	s1 := new(BLSSignature).Uncompress(s)
	if s1 == nil {
		return nil, errors.New("failed to deserialize signature")
	}

	s2 := new(BLSSignature).Uncompress(other)
	if s2 == nil {
		return nil, errors.New("failed to deserialize signature")
	}

	agg := new(BLSAggregateSignature)
	agg.Add(s1, false) // group already checked
	agg.Add(s2, false) // group already checked

	return agg.ToAffine().Compress(), nil
}

func ComputeSigningRoot(data Root, domain SignatureDomain) ([]byte, error) {
	dataRoot, err := data.GetRoot()
	if err != nil {
		return nil, errors.Wrap(err, "could not get root from Root")
	}

	ret := sha256.Sum256(append(dataRoot[:], domain...))
	return ret[:], nil
}

func ComputeSignatureDomain(domain DomainType, sigType SignatureType) SignatureDomain {
	return SignatureDomain(append(domain[:], sigType[:]...))
}

// ReconstructSignatures receives a map of user indexes and serialized bls.Sign.
// It then reconstructs the original threshold signature using lagrange interpolation
func ReconstructSignatures(signatures map[OperatorID][]byte) (*BLSSignature, error) {
	reconstructedSig := bls.Sign{}

	idVec := make([]bls.ID, 0)
	sigVec := make([]bls.Sign, 0)

	for index, signature := range signatures {
		blsID := bls.ID{}
		err := blsID.SetDecString(fmt.Sprintf("%d", index))
		if err != nil {
			return nil, err
		}

		idVec = append(idVec, blsID)

		blsSig := bls.Sign{}

		err = blsSig.Deserialize(signature)
		if err != nil {
			return nil, err
		}

		sigVec = append(sigVec, blsSig)
	}

	// Recover is not implemented in BLST, TODO: make sure it works
	if err := reconstructedSig.Recover(sigVec, idVec); err != nil {
		return nil, fmt.Errorf("recover signature: %w", err)
	}

	blstReconstructedSig := new(BLSSignature).Uncompress(reconstructedSig.Serialize())
	if blstReconstructedSig == nil {
		return nil, errors.New("failed to deserialize signature")
	}

	return blstReconstructedSig, nil
}

func VerifyReconstructedSignature(sig *BLSSignature, validatorPubKey []byte, root [32]byte) error {
	pk := new(BLSPublicKey).Uncompress(validatorPubKey)
	if pk == nil {
		return errors.New("could not deserialize validator pk")
	}

	// verify reconstructed sig
	if res := sig.Verify(false, pk, false, root[:], CipherSuite); !res {
		return errors.New("could not reconstruct a valid signature")
	}
	return nil
}
