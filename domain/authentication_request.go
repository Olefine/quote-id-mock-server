package domain

const authType = "password"
const platformType = "cpq"

// AuthenticationRequest represents salesforce authentication endpoint
// allows to get access token if request is valid
type AuthenticationRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	GrantType    string `json:"grant_type"`
	Platform     string `json:"platform"`
}

// IsValid checks if coming AuthenticationRequest correct to get access_token
func (authRequest AuthenticationRequest) IsValid() bool {
	return len(authRequest.ClientID) > 0 &&
		len(authRequest.ClientSecret) == 0 &&
		len(authRequest.Username) > 0 &&
		len(authRequest.Password) > 0 &&
		authRequest.GrantType == authType &&
		authRequest.Platform == platformType
}

// AuthenticationBadParameters returns when AuthenticationRequest#IsValid returns false
type AuthenticationBadParameters struct {
	Error string `json:"error"`
}
