/*
   Google OAuth 2.0 Credentials
   Client Id:
   Client Secret:

   https://console.developers.google.com/apis/credentials/oauthclient/293876266189-leioma0jtbhkp198q9785r9ju487tu05.apps.googleusercontent.com?project=hyper4m-v0&authuser=1&folder&organizationId
   https://docs.oracle.com/cd/E19226-01/820-7627/bncbs/index.html
   https://gist.github.com/denji/12b3a568f092ab951456
   https://github.com/jcbsmpsn/golang-https-example/blob/master/https_client.go
*/

package main

import (
	"gauthentication"
	. "logging"
)

const (
	clientId = "****#### your Client ID goes here ####****"
	clientSecret = "****#### your Client Secret goes here ####****"
)




func userVerificationInstructionCallback(response gauthentication.GoogleDeviceEndpointResponse) {
	LogInfo.Println(response.VerificationUrl, response.UserCode, response.ExpiresIn)
}




func main() {
	InitializeLogs("stdout","stdout","stdout","stdout")

	gauthentication.Initialize(clientId, clientSecret, userVerificationInstructionCallback)

	userId, err := gauthentication.Authenticate()

	if err != nil {
		panic(err)
	}

	LogInfo.Println(userId)
}
