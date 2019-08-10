package gauthentication

import (
  . "logging"
)

type UserIdentity struct {
  GivenName string
  FamilyName string
  Email string
  IsEmailVerified bool
  PictureUrl string
}

type GoogleDeviceEndpointResponse struct {
    DeviceCode string `json:"device_code"`
    UserCode string `json:"user_code"`
    ExpiresIn int `json:"expires_in"`
    Interval int `json:"interval"`
    VerificationUrl string `json:"verification_url"`
}

type GoogleTokenEndpointResponse struct {
  AccessToken string `json:"access_token"`
  RefreshToken string `json:"refresh_token"`
  ExpiresIn int `json:"expires_in"`
  TokenType string `json:"token_type"`
  IdToken string  `json:"id_token"`
}

/* errors */
type ErrorTerminationSignalReceived struct {}
func (e *ErrorTerminationSignalReceived) Error() string {
  return "authentication process terminated!"
}

/* callback function definitions */
type UserVerificationInstructionCallback func(response GoogleDeviceEndpointResponse)

/* empty callbacks */
func NoUserVerificationInstructionCallback(response GoogleDeviceEndpointResponse) {
  LogTrace.Println("NoUserVerificationInstructionCallback(...) was called.")
}
