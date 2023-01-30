package k6utils

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type FinalChallengeParams struct {
	AppID     string `json:"appID"`
	Challenge string `json:"challenge"`
	FacetID   string `json:"facetID"`
}

type FidoResponseEntry struct {
	Header         Header                       `json:"header"`
	Assertions     []AuthenticatorSignAssertion `json:"assertions"`
	Base64FcParams string                       `json:"fcParams"`
}

type FidoRegistrationResponseBuilder interface {
	Build(aaid string, overriddenSignature string, signatureSignData string,
		privKey string, pubKey string) (*SendUafResponse, error)
}

type FidoRegistrationResponse struct {
	facetId          string
	returnUafRequest ReturnUafRequest
}

func (b *FidoRegistrationResponse) Build(aaid string, overriddenSignature string, signatureSignData string,
	privKey string, pubKey string) (*SendUafResponse, error) {

	var regRequestEntries []RegRequestEntry
	err := json.Unmarshal([]byte(b.returnUafRequest.UafRequest), &regRequestEntries)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling uafRequest")
	}

	regRequestEntry := regRequestEntries[0]

	finalChallengeParams := FinalChallengeParams{
		AppID:     regRequestEntry.Header.AppID,
		Challenge: regRequestEntry.Challenge,
		FacetID:   b.facetId,
	}

	base64FcByte, _ := json.Marshal(finalChallengeParams)
	base64FcString := base64.URLEncoding.EncodeToString(base64FcByte)
	finalChallengeParamsHash := sha256.Sum256([]byte(base64FcString))

	fidoRegistrationAssertion, _ := NewFidoRegistrationSignedAssertions(aaid, signatureSignData, pubKey, privKey, overriddenSignature, finalChallengeParamsHash[:])
	assertions := []AuthenticatorSignAssertion{*fidoRegistrationAssertion}

	regResponseEntry := FidoResponseEntry{
		Header:         regRequestEntry.Header,
		Assertions:     assertions,
		Base64FcParams: base64FcString,
	}

	regResponseEntries := []FidoResponseEntry{regResponseEntry}

	responseJson, err := json.Marshal(regResponseEntries)
	if err != nil {
		return nil, fmt.Errorf("Error marshalling registration response entries")
	}

	sendUafResponse := &SendUafResponse{
		UafResponse: string(responseJson),
		Context:     "",
	}

	return sendUafResponse, nil
}
