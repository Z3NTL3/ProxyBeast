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

type CWD = string

var (
     RootDir CWD
     MX *Controller = &Controller{}
)

const JUDGE = "pool.proxyspace.pro"