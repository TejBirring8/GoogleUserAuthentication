package gauthentication

import (
  "errors"
  "github.com/dgrijalva/jwt-go"
  . "logging"
)

var (
  _mapGoogleServerCerts map[string]string
)




func getGoogleServerCertForToken(token *jwt.Token) (interface{}, error) {
  /* attempt to fetch appropriate key for the JWT token */
  strKeyId := token.Header["kid"].(string)
  strKey := _mapGoogleServerCerts[strKeyId]
  if (strKeyId == "" || strKey == "") {
    return nil, errors.New("could not find Google server certificate for the received token!")
  }
  /* convert the key into correct format for parsing and verification */
  key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(strKey))
  if err != nil {
    return nil, err
  }
  /* return successful*/
  return key, nil
}




func decodeToken(responseGoogleTokenEndpoint GoogleTokenEndpointResponse, mapGoogleServerCerts map[string]string) (UserIdentity, error) {
  LogTrace.Println("decoding JWT token and verifying with Google server certificate...")
  _mapGoogleServerCerts = mapGoogleServerCerts
  var userId UserIdentity

  /* parse and verify the received JWT token */
  parser := jwt.Parser{UseJSONNumber:false, SkipClaimsValidation:true}
  token, err := parser.Parse(responseGoogleTokenEndpoint.IdToken, getGoogleServerCertForToken)
  if err != nil {
    return userId, err
  }

  /* parse the readable/decoded token information into TUserIdentity struct */
  if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
    userId.GivenName = claims["given_name"].(string)
    userId.FamilyName = claims["family_name"].(string)
    userId.Email = claims["email"].(string)
    userId.IsEmailVerified = claims["email_verified"].(bool)
    userId.PictureUrl = claims["picture"].(string)
  } else {
      return userId, errors.New("token not valid!")
  }

  /* finally, return error if user's email has not been verified */
  if userId.IsEmailVerified != true {
    return userId, errors.New("user's email has not been verified!")
  }

  /* return successful */
  LogTrace.Println("decoded and verified JWT token.")
  return userId, nil
}
