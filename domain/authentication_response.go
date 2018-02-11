package domain

const onedaySecods = 86400
const definedTokenType = "bearer"
const definedScopeType = "MajorityStrategies2"

// AuthenticationResponse represents response for authenticated user
type AuthenticationResponse struct {
	Token             string      `json:"access_token"`
	ExpiresIn         int         `json:"expires_in"`
	TokenType         string      `json:"token_type"`
	Scope             string      `json:"scope"`
	RefreshToken      interface{} `json:"refresh_token"`
	RefreshExpiresInt interface{} `json:"refresh_expires_in"`
	DownloadToken     interface{} `json:"download_token"`
}

// NewAuthenticationResponse returns AuthenticationResponse, behide the scence
// generates random token using TokenRepository
func NewAuthenticationResponse(token string) *AuthenticationResponse {
	return &AuthenticationResponse{
		Token:             token,
		ExpiresIn:         onedaySecods,
		TokenType:         definedTokenType,
		Scope:             definedScopeType,
		RefreshToken:      nil,
		RefreshExpiresInt: nil,
		DownloadToken:     nil,
	}
}
