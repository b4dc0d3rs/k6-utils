package k6utils

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type FidoAuthenticationResponseBuilder interface {
	Build(aaid string, overriddenSignature string, signatureSignData string,
		privKey string, pubKey string) (*SendUafResponse, error)
}

type FidoAuthenticationResponse struct {
	facetId                string
	returnUafRequest       ReturnUafRequest
	isKeyRotationSupported bool
	username               string
}

func (b *FidoAuthenticationResponse) Build(aaid string, overriddenSignature string, signatureSignData string,
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

	fidoAuthenticationSignedAssertions, _ := NewFidoAuthenticationSignedAssertions(aaid, pubKey, privKey, overriddenSignature, signatureSignData, finalChallengeParamsHash[:])
	assertions := []AuthenticatorSignAssertion{*fidoAuthenticationSignedAssertions}

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

	context := make(map[string]interface{})
	context["username"] = b.username

	contextJson, err := json.Marshal(context)
	if err != nil {
		// handle error
	}
	contextString := string(contextJson)

	sendUafResponse := &SendUafResponse{
		UafResponse: string(responseJson),
		Context:     contextString,
	}

	return sendUafResponse, nil
}
