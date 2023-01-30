package k6utils

type Transaction struct {
	ContentType string `json:"contentType"`
	Content     string `json:"content"`
}

type MatchCriteria struct {
	UserVerification         int64    `json:"userVerification"`
	AuthenticationAlgorithms []int    `json:"authenticationAlgorithms"`
	AssertionSchemes         []string `json:"assertionSchemes"`
}

type Policy struct {
	Accepted [][]MatchCriteria `json:"accepted"`
}

type Upv struct {
	Major int `json:"major"`
	Minor int `json:"minor"`
}

type Header struct {
	Upv   Upv    `json:"upv"`
	Op    string `json:"op"`
	AppID string `json:"appID"`
}

type RegRequestEntry struct {
	Header      Header        `json:"header"`
	Challenge   string        `json:"challenge"`
	Username    string        `json:"username"`
	Policy      Policy        `json:"policy"`
	Transaction []Transaction `json:"transaction"`
}
