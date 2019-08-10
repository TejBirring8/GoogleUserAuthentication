/*
    Google OAuth 2.0 User Authentication Library for Mobile and Desktop Application
    Author: Tej Birring (tej.birring@gmail.com)

    Overview: https://developers.google.com/identity/protocols/OAuth2,
              https://developers.google.com/identity/protocols/OpenIDConnect
    Obtain clientId and clientSecret from: https://console.developers.google.com/apis/credentials


 */

package gauthentication

import (
  "crypto/tls"
  . "logging"
  "net/http"
)

var (
  // http over TLS (https) client
  client *http.Client
  // Google OAuth 2.0 API
  gClientId string
  gClientSecret string
  /* variables to sync operation of this module */
  chnlTerminateSignal chan bool
  /* callbacks */
  userVerificationInstructionCallback UserVerificationInstructionCallback
)




/* initializes this module */
func Initialize(clientId string, clientSecret string, callback UserVerificationInstructionCallback) {
  /* initialize logging if not already -
  this will be initialized as if in debug mode i.e.
  everything to stdout */
  if IsLoggingInitialized() == false {
    InitializeLogs("stdout", "stdout", "stdout", "stdout")
  }
  /* set clientId and clientSecret */
  gClientId = clientId
  gClientSecret = clientSecret
  /* create a HTTP client struct. for HTTP/TLS transport channel - used to talk to Google server */
  client = &http.Client{
    Transport: &http.Transport{
      TLSClientConfig: &tls.Config{
        InsecureSkipVerify:       false,
        //VerifyPeerCertificate:  <TODO: implement callback for SSL pinning of server cert>
      }, //&tls.Config
    }, //&http.Transport
  }
  /* create a channel to communicate if/when to terminate the Authenticate() process */
  chnlTerminateSignal = make(chan bool)
  /* set callback, called when user is to be informed of the link + code to use */
  userVerificationInstructionCallback = callback
}




/* starts the authentication process */
func Authenticate() (UserIdentity, error) {
  /* obtain public key of the server */
  mapGoogleServerCerts, err := obtainGoogleServerCerts()
  if err != nil {
    return UserIdentity{}, err
  }
  LogTrace.Println(mapGoogleServerCerts)

  /* make an authentication request */
  responseGoogleDeviceEndpoint, err := obtainUsercode()
  if err != nil {
    return UserIdentity{}, err
  }
  LogTrace.Println(responseGoogleDeviceEndpoint);
  /* execute user-defined callback funct w/ purpose of notifying user of URL and code */
  userVerificationInstructionCallback(responseGoogleDeviceEndpoint)

  /* request for token until it is obtained or timeout reached */
  responseGoogleTokenEndpoint, err := obtainToken(responseGoogleDeviceEndpoint)
  if err != nil {
    return UserIdentity{}, err
  }
  LogTrace.Println(responseGoogleTokenEndpoint)

  /* decode obtained user identity info */
  userIdentity, err := decodeToken(responseGoogleTokenEndpoint, mapGoogleServerCerts)
  if err != nil {
    return UserIdentity{}, err
  }
  LogTrace.Println(userIdentity)

  /* return successful */
  return userIdentity, nil
}




/* terminates the authentication process */
func Terminate() {
  chnlTerminateSignal <- true
}
