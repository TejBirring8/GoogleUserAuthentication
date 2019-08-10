package gauthentication

import (
  "encoding/json"
  "errors"
  . "logging"
  "net/url"
  "time"
)




func obtainToken(r GoogleDeviceEndpointResponse) (GoogleTokenEndpointResponse, error) {
  LogTrace.Println("obtaining token from Google Token Endpoint...")

  var response GoogleTokenEndpointResponse
  chnlResponse := make(chan GoogleTokenEndpointResponse)   //channel to rcv response struct
  chnlErr := make(chan error)                               //channel to rcv errors


  /* create a concurrent ### goroutine ### to poll the endpoint,
     and wait on its completion! */
  go func() {
    elapsed_time := 0
    /* create a ticket which requests a token from the Google Token Endpoint
    periodically at the set interval of time - until either the token is
    received or the expiry time has been reached */
    ticker := time.NewTicker(time.Second * time.Duration(r.Interval))
    for t:= range ticker.C {
      LogTrace.Println("Tick at", t)

      /* terminate operation if 'terminate operation' signal is received */
      select {
        case term:= <-chnlTerminateSignal:
          if term == true {
            LogTrace.Println("'terminate operation' signal received.")
            ticker.Stop()
            chnlErr <- &ErrorTerminationSignalReceived{}
          }
        default:
      }

      /* request the token from the endpoint */
      resp, err := client.PostForm("https://www.googleapis.com/oauth2/v4/token",
                      url.Values{"client_id": {gClientId},
                      "client_secret": {gClientSecret},
                      "code": {r.DeviceCode},
                      "grant_type": {"http://oauth.net/grant_type/device/1.0"}})
      if err != nil {
        chnlErr <- err
      }

      /* parse the response received from the endpoint into struct */
      dec := json.NewDecoder(resp.Body)
      response := GoogleTokenEndpointResponse{}
      err = dec.Decode(&response)
      if err != nil {
        chnlErr <- err
      }
      resp.Body.Close()

      /* increment the elapsed time by the interval */
      elapsed_time += r.Interval

      /* return error if expiry time has elapsed */
      if elapsed_time >= r.ExpiresIn {
        ticker.Stop()
        chnlErr <- errors.New("authentication process timed out!")
      }

      /* return successful */
      if (response.IdToken != "") {
        ticker.Stop()
        chnlResponse <- response
      }
    }
  }() /* end of ### goroutine ### */


  /* uncomment to TEST 'terminate operation' signal handling in goroutine above */
  /*
  time.Sleep(time.Duration(15) * time.Second)
  Terminate()
  */

  /* respond to whichever is received, an error, or a token */
  select {
    case err := <- chnlErr:
      return GoogleTokenEndpointResponse{}, err
    case response = <- chnlResponse:
      if response.IdToken == "" {
        return GoogleTokenEndpointResponse{}, errors.New("failed to obtain a token from Google Token Endpoint.")
      }
  }

  /* return successful */
  LogTrace.Println("obtained token from Google Token Endpoint.")
  return response, nil
}
