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

import "os"

type CWD = string

var (
     RootDir CWD
     FD = map[string]*os.File{
          InputFile: nil,
          SaveFile: nil,
     }
     MX *Controller = &Controller{}
)

const JUDGE = "https://pool.proxyspace.pro/judge.php"