package chain

const (
	ChainIdAbbreviationIris     = "i"
	ChainIdAbbreviationStars    = "s"
	ChainIdAbbreviationJuno     = "j"
	ChainIdAbbreviationUptick   = "u"
	ChainIdAbbreviationOmniflix = "o"

	ChainIdValueIirs     = "gon-irishub-1"
	ChainIdValueStars    = "elgafar-1"
	ChainIdValueJuno     = "uni-6"
	ChainIdValueUptick   = "uptick_7000-2"
	ChainIdValueOmniflix = "gon-flixnet-1"

	ChainGRPCIris     = "34.80.93.133:9090"
	ChainGRPCStars    = "grpc-1.elgafar-1.stargaze-apis.com:26660"
	ChainGRPCJuno     = "juno-testnet-grpc.polkachu.com:12690"
	ChainGRPCUptick   = "52.220.252.160:9090"
	ChainGRPCOmniflix = "65.21.93.56:9090"

	ChainRPCIris  = "http://34.80.93.133:26657/"
	ChainRPCStars = "https://rpc.elgafar-1.stargaze-apis.com:443/"
	//ChainRPCJuno  = "https://rpc.uni.junonetwork.io:443/"
	ChainRPCJuno     = "https://rpc.uni.juno.deuslabs.fi:443/"
	ChainRPCUptick   = "http://52.220.252.160:26657/"
	ChainRPCOmnilfix = "http://65.21.93.56:26657/"
)

type (
	Class struct {
		ID      string
		Name    string
		Schema  string
		Creator string
		Uri     string
		UriHash string
		Data    string
	}

	NFT struct {
		ID      string
		Name    string
		URI     string
		Data    string
		Owner   string
		URIHash string
	}

	Chain interface {
		GetTx(txHash, txType string) (any, error)
		GetNFT(classID, nftID string) (*NFT, error)
		HasNFT(classID, nftID string) bool
		GetClass(classID string) (*Class, error)
		HasClass(classID string) bool
		Close()
	}

	Registry struct {
		chains map[string]Chain
	}
)

func NewRegistry() *Registry {
	return &Registry{
		chains: map[string]Chain{
			ChainIdAbbreviationIris:     NewIris(),
			ChainIdAbbreviationStars:    NewStargaze(),
			ChainIdAbbreviationJuno:     NewJuno(),
			ChainIdAbbreviationUptick:   NewUptick(),
			ChainIdAbbreviationOmniflix: NewOmniflix(),
		},
	}
}

func (cr *Registry) GetChain(chainID string) Chain {
	return cr.chains[chainID]
}

func (cr *Registry) GetChains() map[string]Chain {
	return cr.chains
}
