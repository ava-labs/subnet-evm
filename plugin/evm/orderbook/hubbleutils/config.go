package hubbleutils

var (
	ChainId           int64
	VerifyingContract string
	HState            *HubbleState
)

func SetChainIdAndVerifyingSignedOrdersContract(chainId int64, verifyingContract string) {
	ChainId = chainId
	VerifyingContract = verifyingContract
}

func SetHubbleState(hState *HubbleState) {
	HState = hState
}
