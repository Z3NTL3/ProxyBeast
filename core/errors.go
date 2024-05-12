/*
     ProxyBeast GUI

The ultimate proxy checker
       by @z3ntl3

    [proxy.pix4.dev]

License: GNU
Note: Please do give us a star on Github, if you like ProxyBeast

[App core]
*/

package core

import "errors"

var (
	ErrPropsInvalid = errors.New("[props] of invalid type")
	ErrPropsOPInvalid = errors.New("[props] operation is invalid")
     ErrFileNotProvided = errors.New("Save and/or input file not provided")
     ErrFDNil = errors.New("Please provide input proxy file, and save location.")
     ErrInvalidProxyURI = errors.New("FILE:FORMAT -> Invalid proxy URI was found")
     ErrOngoingCheck = errors.New("There's already an ongoing check")
     ErrNoProxiesFound = errors.New("There are no proxies detected in your file")
     ErrBadProxy = errors.New("Bad proxy")
     ErrTimeoutString = errors.New("Timeout can only be like a duration string like 10s/5s/1000ms etc")
)