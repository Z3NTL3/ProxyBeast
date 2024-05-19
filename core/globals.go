/*

     ProxyBeast GUI

The ultimate proxy checker
       by @z3ntl3

    [proxy.pix4.dev]

License: GNU
Note: If you've liked ProxyBeast, please consider starring our Github repository.

[App core]
*/

package core

import "time"

type CWD = string

var (
     RootDir CWD
     MX *Controller = &Controller{}
     
     DefaultTimeout time.Duration = time.Second * 8
     DefaultPoolSize uint32 = 10_000
     DefaultJudge string = "https://pool.proxyspace.pro/judge.php"
     DefaultConfigFile string = "config.json"
)
