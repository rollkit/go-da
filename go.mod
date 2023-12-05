module github.com/rollkit/go-da

go 1.21.1

require (
	github.com/celestiaorg/celestia-app v1.6.0
	github.com/gogo/protobuf v1.3.3
	github.com/stretchr/testify v1.8.4
	google.golang.org/grpc v1.58.3
)

require (
	github.com/celestiaorg/merkletree v0.0.0-20210714075610-a84dc3ddbbe4 // indirect
	github.com/celestiaorg/rsmt2d v0.11.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/klauspost/cpuid/v2 v2.1.1 // indirect
	github.com/klauspost/reedsolomon v1.11.8 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/tendermint/tendermint v0.34.28 // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/sync v0.3.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20231009173412-8bfb1ae86b6c // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1

replace github.com/tendermint/tendermint => github.com/celestiaorg/celestia-core v1.29.0-tm-v0.34.29
