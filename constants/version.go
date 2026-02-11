package constants

// Version is the current version of the application.
// It can be overridden at build time using ldflags:
//
//	go build -ldflags "-X genius/constants.Version=1.0.0"
var Version = "0.0.10"
