package verifier

const (
	ReasonParamsFormatIncorrect = "Params: format is incorrect"
	ReasonParamsChainIdEmpty    = "Params: chainId not found"
	ReasonParamsChainIdError    = "Params: chainId is error"

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

	ReasonIbcDestPortNotMatch        = "IBC: dest port not match"
	ReasonIbcDestChanNotMatch        = "IBC: dest channel not match"
	ReasonIbcClassNotMatch           = "IBC: ibc class not match"
	ReasonIbcOriginalClassIdNotMatch = "IBC: original class id not match"

	ReasonRaceUnexpectedFlowPath  = "Race: race flow unexpected"
	ReasonRaceFirstLastSenderNotMatch = "Race: first and last sender not match"
	ReasonRaceDataUnachievable    = "Race: data is unachievable"
	ReasonRaceStartTooEarly = "Race: you start too early"
)