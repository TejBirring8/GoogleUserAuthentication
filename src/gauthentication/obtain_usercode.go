package gauthentication

import (
  "encoding/json"
  "errors"
  . "logging"
  "net/url"
)




func obtainUsercode() (GoogleDeviceEndpointResponse, error) {
  LogTrace.Println("obtaining usercode from Google Device Endpoint...")
  var response GoogleDeviceEndpointResponse
  /* send a request to Google Device Endpoint with 'client ID' and 'scope' */
  resp, err := client.PostForm(
    "https://accounts.google.com/o/oauth2/device/code",
     url.Values{"client_id": {gClientId}, "scope": {"profile email"}})
  if err != nil {
    return GoogleDeviceEndpointResponse{}, err
  }
  /* decode the response */
  dec := json.NewDecoder(resp.Body)
  err = dec.Decode(&response)
  if err != nil {
    return GoogleDeviceEndpointResponse{}, err
  }
  resp.Body.Close()
  /* error if nothing in the decoded response */
  if response.UserCode == "" {
    return GoogleDeviceEndpointResponse{}, errors.New("failed to obtain usercode from Google Device Endpoint.")
  }
  /* return succesful */
  LogTrace.Println("obtained usercode from Google Device Endpoint.")
  return response, nil
}
