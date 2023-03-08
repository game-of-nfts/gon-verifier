package verifier

const (
	ReasonParamsFormatIncorrect = "Params: format is incorrect"
	ReasonParamsChainIdEmpty    = "Params: chainId is empty"

	ReasonTxResultUnexpected   = "Tx: result is unexpected"
	ReasonTxResultUnachievable = "Tx: result is unachievable"
	ReasonTxResultUnsuccessful = "Tx: result is unsuccessful"
	ReasonTxMsgSenderNotMatch  = "Tx: sender not match register address"

	ReasonClassNotFound        = "Class: not found"
	ReasonClassCreatorNotMatch = "Class: creator not match register address"
	ReasonClassDataInvalid     = "Class: data is invalid"
	ReasonClassUrIEmpty        = "Class: uri is empty"

	ReasonNftNotFound          = "NFT: not found"
	ReasonNftOwnerNotMatch     = "NFT: initial owner not match register address"
	ReasonNftRecipientNotMatch = "NFT: recipient not match register address"
	ReasonNftTokenIdNotMatch   = "NFT: token id not match"
	ReasonNftUriEmpty          = "NFT: uri is empty"
	ReasonNftDataEmpty         = "NFT: data is empty"
)
