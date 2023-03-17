package types

import (
	"github.com/herumi/bls-eth-go-binary/bls"
	blst "github.com/supranational/blst/bindings/go"
	"math/big"
)

var (
	curveOrder = new(big.Int)
)

// InitBLS initializes BLS
func InitBLS() {
	_ = bls.Init(bls.BLS12_381)
	_ = bls.SetETHmode(bls.EthModeDraft07)

	curveOrder, _ = curveOrder.SetString(bls.GetCurveOrder(), 10)
}

type BLSPublicKey = blst.P1Affine
type BLSSecretKey = blst.SecretKey
type BLSSignature = blst.P2Affine
type BLSAggregateSignature = blst.P2Aggregate
type BLSAggregatePublicKey = blst.P1Aggregate

var CipherSuite = []byte("BLS_SIG_BLS12381G2_XMD:SHA-256_SSWU_RO_POP_")
