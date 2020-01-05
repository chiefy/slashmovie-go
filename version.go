package main

import "fmt"

// Name is the application name
const Name = "slashmovie"

// Version is the application Version
var Version string = "dev"

// Date is the built date and time
var Date string

// Commit is the commit in which the package is based
var Commit string

// GetVersion returns the version as a string
func GetVersion() string {
	return fmt.Sprintf("%s Version: %s Commit: %s Date: %s", Name, Version, Commit, Date)
}
