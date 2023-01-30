package k6utils

import (
	"encoding/base64"
	"fmt"
	"math/rand"
)

const tagAaid int16 = 0x2E0B
const tagAssertionInfo int16 = 0x2E0E
const tagFinalChallengeHash int16 = 0x2E0A
const tagKeyID int16 = 0x2E09
const tagCounters int16 = 0x2E0D
const tagPubKey int16 = 0x2E0C
const tagUafv1Krd int16 = 0x3E03
const tagSignature int16 = 0x2E06
const tagAttestationBasicSurrogate int16 = 0x3E08
const tagUafv1RegAssertion int16 = 0x3E01
const tagAuthenticatorNonce int16 = 0x2E0F
const tagUafv1AuthAssertion int16 = 0x3E02
const tagUafv1SignedData int16 = 0x3E04
const tagTransactionContextHash = 0x2E10

// this might cause an issue
const algSignSecp256R1EcdsaSha256Raw = 0x0001
const algKeyEccX962Raw = 0x0100

const authModeUserVerified int8 = 1
const authModeTxnContentVerified = 2

const uafv1tlv string = "UAFV1TLV"

type AuthenticatorSignAssertion struct {
	Assertion       string `json:"assertion"`
	AssertionScheme string `json:"assertionScheme"`
}

type SignedAssertionsBuilder interface {
	Build(string, string) (*AuthenticatorSignAssertion, error)
}

type FidoRegistrationSignedAssertionsBuilder struct {
	assertions *Assertions
}
type FidoAuthenticationSignedAssertionsBuilder struct {
	assertions *Assertions
}

type Assertions struct {
	aaid                     string
	authenticatorVersion     int16
	signCounter              int32
	registrationCounter      int32
	publicKey                []byte
	authenticationMode       int8
	signatureAlgAndEncoding  int16
	publicKeyAlgAndEncoding  int16
	finalChallengeParamsHash []byte
	transactionTextHash      []byte
	signatureSignedData      string
	overriddenSignature      string
}

func NewFidoRegistrationSignedAssertions(aaid string, signatureSignedData string,
	pubKeyBas64Encoded string, privKeyBase64Encoded string, overriddenSignature string, finalChallengeParamsHash []byte) (*AuthenticatorSignAssertion, error) {

	const authenticatorVersion int16 = 1
	const signCounter int32 = 0
	const registrationCounter int32 = 0

	pubKeyBytes, _ := base64.StdEncoding.DecodeString(pubKeyBas64Encoded)

	assertion := &Assertions{
		aaid:                     aaid,
		publicKey:                pubKeyBytes,
		authenticatorVersion:     authenticatorVersion,
		signCounter:              signCounter,
		registrationCounter:      registrationCounter,
		signatureAlgAndEncoding:  algSignSecp256R1EcdsaSha256Raw,
		overriddenSignature:      overriddenSignature,
		signatureSignedData:      signatureSignedData,
		finalChallengeParamsHash: finalChallengeParamsHash,
	}

	fidoRegistrationSignedAssertionsBuilder := FidoRegistrationSignedAssertionsBuilder{assertions: assertion}

	return fidoRegistrationSignedAssertionsBuilder.Build(privKeyBase64Encoded, pubKeyBas64Encoded)
}

func NewFidoAuthenticationSignedAssertions(aaid string, pubKeyBas64Encoded string, privKeyBase64Encoded string, overriddenSignature string, signatureSignedData string, finalChallengeParamsHash []byte) (*AuthenticatorSignAssertion, error) {

	const authenticatorVersion int16 = 1
	const signCounter int32 = 0
	const registrationCounter int32 = 0

	pubKeyBytes, _ := base64.StdEncoding.DecodeString(pubKeyBas64Encoded)

	assertion := &Assertions{
		aaid:                     aaid,
		publicKey:                pubKeyBytes,
		authenticatorVersion:     authenticatorVersion,
		signCounter:              signCounter,
		registrationCounter:      registrationCounter,
		signatureAlgAndEncoding:  algSignSecp256R1EcdsaSha256Raw,
		overriddenSignature:      overriddenSignature,
		signatureSignedData:      signatureSignedData,
		authenticationMode:       authModeUserVerified,
		finalChallengeParamsHash: finalChallengeParamsHash,
	}

	fidoAuthenticationAssertionsBuilder := FidoAuthenticationSignedAssertionsBuilder{assertions: assertion}

	return fidoAuthenticationAssertionsBuilder.Build(privKeyBase64Encoded, pubKeyBas64Encoded)
}

func (fra *FidoRegistrationSignedAssertionsBuilder) Build(privKeyStr string, pubKeyStr string) (*AuthenticatorSignAssertion, error) {
	// Generate assertion structure
	tlvObjectAaid := NewFidoUafTlvObject(tagAaid, []byte(fra.assertions.aaid))

	var assertionInfoContentLength = 2 + 1 + 2 + 2
	tlvObjectAssertionInfo := NewFidoUafTlvObjectWithSize(tagAssertionInfo, int16(assertionInfoContentLength))
	tlvObjectAssertionInfo.PutInt16(fra.assertions.authenticatorVersion)
	tlvObjectAssertionInfo.PutInt8(fra.assertions.authenticationMode)
	tlvObjectAssertionInfo.PutInt16(fra.assertions.signatureAlgAndEncoding)
	tlvObjectAssertionInfo.PutInt16(fra.assertions.publicKeyAlgAndEncoding)

	tlvFinalChallengeHash := NewFidoUafTlvObject(tagFinalChallengeHash, fra.assertions.finalChallengeParamsHash)

	tlvObjectKeyId := NewFidoUafTlvObject(tagKeyID, []byte("any"))

	objectCountersContentLength := 4 + 4
	tlvObjectCounters := NewFidoUafTlvObjectWithSize(tagCounters, int16(objectCountersContentLength))
	tlvObjectCounters.PutInt32(fra.assertions.signCounter)
	tlvObjectCounters.PutInt32(fra.assertions.registrationCounter)

	tlvObjectPubKey := NewFidoUafTlvObject(tagPubKey, fra.assertions.publicKey)

	tlvObjectKeyRegistrationData := NewFidoUafTlvObjectFromArray(tagUafv1Krd, []*FidoUafTlvObject{
		tlvObjectAaid,
		tlvObjectAssertionInfo,
		tlvFinalChallengeHash,
		tlvObjectKeyId,
		tlvObjectCounters,
		tlvObjectPubKey,
	}...)

	// Generate attestation basic surrogate
	// allow the ability to override the data signed in the response
	krdSignature, err := SignDataLocal(tlvObjectKeyRegistrationData.GetByteArray(), privKeyStr, pubKeyStr)

	if err != nil {
		return nil, fmt.Errorf("Error creating signature")
	}

	tlvObjectSignature := NewFidoUafTlvObject(tagSignature, krdSignature)

	tlvObjectAttestationBasicSurrogate := NewFidoUafTlvObjectFromArray(tagAttestationBasicSurrogate, tlvObjectSignature)

	tlvObjectRegistrationAssertion := NewFidoUafTlvObjectFromArray(tagUafv1RegAssertion, tlvObjectKeyRegistrationData, tlvObjectAttestationBasicSurrogate)

	assertionBytes := tlvObjectRegistrationAssertion.GetByteArray()

	encoder := base64.URLEncoding.WithPadding(base64.NoPadding)
	base64EncodedAssertion := encoder.EncodeToString(assertionBytes)

	return &AuthenticatorSignAssertion{
		Assertion:       base64EncodedAssertion,
		AssertionScheme: uafv1tlv,
	}, nil
}

func (fra *FidoAuthenticationSignedAssertionsBuilder) Build(privKeyStr string, pubKeyStr string) (*AuthenticatorSignAssertion, error) {
	// Generate assertion structure
	tlvObjectAaid := NewFidoUafTlvObject(tagAaid, []byte(fra.assertions.aaid))

	var assertionInfoContentLength = 2 + 1 + 2
	tlvObjectAssertionInfo := NewFidoUafTlvObjectWithSize(tagAssertionInfo, int16(assertionInfoContentLength))
	tlvObjectAssertionInfo.PutInt16(fra.assertions.authenticatorVersion)
	tlvObjectAssertionInfo.PutInt8(fra.assertions.authenticationMode)
	tlvObjectAssertionInfo.PutInt16(fra.assertions.signatureAlgAndEncoding)

	authenticatorNonce := make([]byte, 8)
	rand.Read(authenticatorNonce)
	tlvNonceObject := NewFidoUafTlvObject(tagAuthenticatorNonce, authenticatorNonce)

	tlvFinalChallengeHash := NewFidoUafTlvObject(tagFinalChallengeHash, fra.assertions.finalChallengeParamsHash)

	tlvTransactionContentHash := NewFidoUafTlvObject(tagTransactionContextHash, fra.assertions.transactionTextHash)

	tlvObjectKeyId := NewFidoUafTlvObject(tagKeyID, []byte("any"))

	var objectCountersContentLength = 4
	tlvObjectCounters := NewFidoUafTlvObjectWithSize(tagCounters, int16(objectCountersContentLength))
	tlvObjectCounters.PutInt32(fra.assertions.signCounter)

	tlvObjectKeyAuthenticationData := NewFidoUafTlvObjectFromArray(tagUafv1SignedData, []*FidoUafTlvObject{
		tlvObjectAaid,
		tlvObjectAssertionInfo,
		tlvNonceObject,
		tlvFinalChallengeHash,
		tlvTransactionContentHash,
		tlvObjectKeyId,
		tlvObjectCounters,
	}...)

	// Generate attestation basic surrogate
	// allow the ability to override the data signed in the response
	krdSignature, err := SignDataLocal(tlvObjectKeyAuthenticationData.GetByteArray(), privKeyStr, pubKeyStr)

	if err != nil {
		return nil, fmt.Errorf("Error creating signature")
	}

	tlvObjectSignature := NewFidoUafTlvObject(tagSignature, krdSignature)

	tlvObjectAuthenticationAssertion := NewFidoUafTlvObjectFromArray(tagUafv1AuthAssertion, tlvObjectKeyAuthenticationData, tlvObjectSignature)

	assertionBytes := tlvObjectAuthenticationAssertion.GetByteArray()

	encoder := base64.URLEncoding.WithPadding(base64.NoPadding)
	base64EncodedAssertion := encoder.EncodeToString(assertionBytes)

	return &AuthenticatorSignAssertion{
		Assertion:       base64EncodedAssertion,
		AssertionScheme: uafv1tlv,
	}, nil
}
