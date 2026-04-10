package accountschecker

type AccountData struct {
	Address  string  `json:"address"`
	Nonce    uint64  `json:"nonce"`
	Balance  string  `json:"balance"`
	RootHash *string `json:"rootHash"`
}

// IsValid checks if the account data is valid according to the specifications.
func (a *AccountData) IsValid() bool {
	// The response is valid if one ot the following fields are not default: rootHash, nonce, balance.
	if a.Nonce != 0 {
		return true
	}
	if a.Balance != "" && a.Balance != "0" {
		return true
	}
	if a.RootHash != nil && *a.RootHash != "" {
		return true
	}
	return false
}
