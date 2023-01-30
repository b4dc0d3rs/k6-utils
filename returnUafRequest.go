package k6utils

const lifetimeMillis int64 = 300000

// UafStatusCode An enum to define FIDO UafStatusCode
type UafStatusCode int

const (
	OK       UafStatusCode = 1200
	ACCEPTED               = 1202
)

func (c UafStatusCode) code() int {
	return int(c)
}

func (c UafStatusCode) fromCode() UafStatusCode {
	switch c {
	case OK:
		return OK
	case ACCEPTED:
		return ACCEPTED
	default:
		return 0
	}
}

// Operation An enum used to define FIDO operation which a user requests
type Operation int

const (
	Reg = iota
	Auth
	Dereg
)

func (o Operation) fromOperation() string {
	switch o {
	case Reg:
		return "Reg"
	case Auth:
		return "Auth"
	case Dereg:
		return "Dereg"
	default:
		return ""
	}
}

type ReturnUafRequest struct {
	StatusCode     UafStatusCode `json:"statusCode"`
	UafRequest     string        `json:"uafRequest"`
	Op             Operation     `json:"op"`
	LifetimeMillis int64         `json:"lifetimeMillis"`
}

func NewFidoRegistrationReturnUafRequest(uafRequest string) *ReturnUafRequest {
	return &ReturnUafRequest{
		StatusCode:     OK,
		UafRequest:     uafRequest,
		Op:             Reg,
		LifetimeMillis: lifetimeMillis,
	}
}

func NewFidoAuthenticationReturnUafRequest(uafRequest string) *ReturnUafRequest {
	return &ReturnUafRequest{
		StatusCode:     OK,
		UafRequest:     uafRequest,
		Op:             Auth,
		LifetimeMillis: lifetimeMillis,
	}
}
