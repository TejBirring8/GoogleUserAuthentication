package logging

import (
  "errors"
  "io"
  "io/ioutil"
  . "log"
  "os"
)

var (
  LogTrace      *Logger
  LogInfo       *Logger
  LogWarning    *Logger
  LogError      *Logger
  isInitialized bool = false
)

/* converts input string to the appropriate 'io.Writer' reference */
func parseStringArg(arg string) io.Writer {
  switch (arg) {
    case "stdout":
      return os.Stdout
    case "stderr":
      return os.Stderr
    default:
      return ioutil.Discard
  }
}

/* call this ONCE in the application to initialize the logging system */
func InitializeLogs(traceOutput string, infoOutput string, warningOutput string, errOutput string) {
  /* return error if already initialized, prevent bad use of this module */
  if isInitialized {
    panic(errors.New("logging already initialized!"))
  }
  /* initialize the logging system */
  LogTrace = New(parseStringArg(traceOutput),"", Lshortfile | Ldate | Ltime)
  LogInfo = New(parseStringArg(infoOutput),"", Lshortfile | Ldate | Ltime)
  LogWarning = New(parseStringArg(warningOutput),"", Lshortfile | Ldate | Ltime)
  LogError = New(parseStringArg(errOutput),"", Lshortfile | Ldate | Ltime)
  /* set flag */
  isInitialized = true
}

func IsLoggingInitialized() bool {
  return isInitialized;
}
