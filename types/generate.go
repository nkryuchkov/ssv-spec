package types

//go:generate rm -f ./operator_encoding.go
//go:generate go run .../fastssz/sszgen --path operator.go --exclude-objs OperatorID

//go:generate rm -f ./share_encoding.go
//go:generate go run .../fastssz/sszgen --path share.go --include ./operator.go,./messages.go,./signer.go

//go:generate rm -f ./messages_encoding.go
//go:generate go run .../fastssz/sszgen --path messages.go --exclude-objs ValidatorPK,MessageID,MsgType

//go:generate rm -f ./beacon_types_encoding.go
//go:generate go run .../fastssz/sszgen --path beacon_types.go --include $GOPATH/pkg/mod/github.com/attestantio/go-eth2-client@v0.15.7/spec/phase0 --exclude-objs BeaconNetwork,BeaconRole

//go:generate rm -f ./partial_sig_message_encoding.go
//go:generate go run .../fastssz/sszgen --path partial_sig_message.go --include $GOPATH/pkg/mod/github.com/attestantio/go-eth2-client@v0.15.1/spec/phase0,./signer.go,./operator.go --exclude-objs PartialSigMsgType

//go:generate rm -f ./consensus_data_encoding.go
//go:generate go run .../fastssz/sszgen --path consensus_data.go --include ./operator.go,./signer.go,./partial_sig_message.go,./beacon_types.go,$GOPATH/pkg/mod/github.com/attestantio/go-eth2-client@v0.15.1/spec/phase0,$GOPATH/pkg/mod/github.com/attestantio/go-eth2-client@v0.15.1/spec,$GOPATH/pkg/mod/github.com/attestantio/go-eth2-client@v0.15.1/spec/altair --exclude-objs Contributions,BeaconNetwork,BeaconRole
