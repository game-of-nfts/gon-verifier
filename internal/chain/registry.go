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

	ChainRPCIris     = "http://34.80.93.133:26657/"
	ChainRPCStars    = "https://rpc.elgafar-1.stargaze-apis.com:443/"
	ChainRPCJuno     = "https://rpc.uni.junonetwork.io:443/"
	ChainRPCUptick   = "http://52.220.252.160:26657/"
	ChainRPCOmnilfix = "http://65.21.93.56:26657/"
)

type (
	TxResultHttp struct {
		Jsonrpc string `json:"jsonrpc"`
		ID      int    `json:"id"`
		Result  struct {
			Hash     string `json:"hash"`
			Height   string `json:"height"`
			Index    int    `json:"index"`
			TxResult struct {
				Code      int    `json:"code"`
				Data      string `json:"data"`
				Log       string `json:"log"`
				Info      string `json:"info"`
				GasWanted string `json:"gas_wanted"`
				GasUsed   string `json:"gas_used"`
				Events    []struct {
					Type       string `json:"type"`
					Attributes []struct {
						Key   string `json:"key"`
						Value string `json:"value"`
						Index bool   `json:"index"`
					} `json:"attributes"`
				} `json:"events"`
				Codespace string `json:"codespace"`
			} `json:"tx_result"`
			Tx    string `json:"tx"`
			Proof struct {
				RootHash string `json:"root_hash"`
				Data     string `json:"data"`
				Proof    struct {
					Total    string `json:"total"`
					Index    string `json:"index"`
					LeafHash string `json:"leaf_hash"`
					Aunts    []any  `json:"aunts"`
				} `json:"proof"`
			} `json:"proof"`
		} `json:"result"`
	}

	TxResult struct {
		Sender string
	}

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
		GetTx(txHash string) (*TxResult, error)
		GetNFT(classID, nftID string) (*NFT, error)
		HasNFT(classID, nftID string) bool
		GetClass(classID string) (*Class, error)
		HasClass(classID string) bool
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
