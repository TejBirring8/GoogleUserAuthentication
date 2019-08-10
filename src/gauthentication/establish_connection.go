package gauthentication

import (
  "encoding/json"
  "errors"
  . "logging"
  _ "net/url"
)




func obtainGoogleServerCerts() (map[string]string, error) {
  var mapGoogleServerCerts map[string]string
  LogTrace.Println("obtaining Google server certificates...")
  /* obtain Google server certificates (public keys) */
  resp, err := client.Get("https://www.googleapis.com/oauth2/v1/certs")
  if err != nil {
    return nil, err
  }
  /* convert response to a [string]->string map of ID->certificate */
  dec := json.NewDecoder(resp.Body)
  err = dec.Decode(&mapGoogleServerCerts)
  if err != nil {
    return nil, err
  }
  resp.Body.Close()
  /* fail if no certificates where obtained */
  if(len(mapGoogleServerCerts) < 1) {
    return nil, errors.New("Failed to obtain Google server certificates.")
  }
  //LogTrace.Println("Obtained Google server certificates:", mapGoogleServerCerts)
  /* return as succesful */
  LogTrace.Println("obtained Google server certificates.")
  return mapGoogleServerCerts, nil
}
