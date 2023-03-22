package types

const (
	// Issue Denom on Iris
	EventTypeIssueDenom   = "issue_denom"
	AttributeDenomId      = "denom_id"
	AttributeDenomCreator = "creator"
	AttributeMsgSender    = "sender"

	// Mint NFT on Iris
	EventTypeNftMint      = "mint_nft"
	AttributeKeyTokenId   = "token_id"
	AttributeKeyRecipient = "recipient"
	// AttributeDenomId = "denom_id"

	// Transfer NFT on Iris
	EventTypeNftTransfer = "transfer_nft"

	// IBC NFT Transfer
	EventTypeIbcNftTransfer    = "ibc_nft_transfer"
	EventTypeIbcSendPacket     = "send_packet"
	AttributeKeyIbcPackageData = "packet_data"
	AttributeKeySender         = "sender"
	AttributeKeyReceiver       = "receiver"
	AttributeKeyDestPort       = "packet_dst_port"
	AttributeKeyDestChan       = "packet_dst_channel"

	EventTypeWasm = "wasm"
	// AttributeKeySender = "sender"
	// AttributeKeyRecipient = "recipient" ics-721 contract addr
	// AttributeKeyTokenId = "token_id"
)
