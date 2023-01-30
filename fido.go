package k6utils

import (
	"encoding/json"
	"fmt"
)

type SendUafResponse struct {
	UafResponse string `json:"uafResponse"`
	Context     string `json:"context"`
}

func (k6utils *K6Utils) GenerateRegistrationResponse(aaid string, uafRequest string,
	trustedFacetId string, overriddenSignature string, signatureSignData string,
	privKey string, pubKey string) (string, error) {
	fidoRegistrationUafRequest := NewFidoRegistrationReturnUafRequest(uafRequest)

	fidoRegistrationResponse := FidoRegistrationResponse{
		facetId:          trustedFacetId,
		returnUafRequest: *fidoRegistrationUafRequest,
	}

	sendUafResponse, _ := fidoRegistrationResponse.Build(aaid, overriddenSignature, signatureSignData, privKey, pubKey)

	fidoRegistrationResponseString, err := json.Marshal(sendUafResponse)
	if err != nil {
		fmt.Println("error:", err)
		return "", err
	}

	return string(fidoRegistrationResponseString), nil
}

func (k6utils *K6Utils) GenerateAuthenticationResponse(aaid string, uafRequest string,
	trustedFacetId string, overriddenSignature string, signatureSignData string,
	privKey string, pubKey string, username string) (string, error) {
	fidoAuthenticationUafRequest := NewFidoAuthenticationReturnUafRequest(uafRequest)

	fidoAuthenticationResponse := FidoAuthenticationResponse{
		facetId:          trustedFacetId,
		returnUafRequest: *fidoAuthenticationUafRequest,
		username:         username,
	}

	sendUafResponse, _ := fidoAuthenticationResponse.Build(aaid, overriddenSignature, signatureSignData, privKey, pubKey)

	fidoRegistrationResponseString, err := json.Marshal(sendUafResponse)
	if err != nil {
		fmt.Println("error:", err)
		return "", err
	}

	return string(fidoRegistrationResponseString), nil
}
