# ProxyBeast
 Proxy Beast is a high-performance proxy checker that offers precise and rapid testing for HTTP, HTTPS, SOCK4, and SOCKS5 protocols. Its advanced capabilities make it an ideal choice for those who require reliable and efficient proxy testing.

# Update
![image](https://github.com/Z3NTL3/ProxyBeast/assets/48758770/a926b578-6caa-44bc-a446-cb73235a9eab)

This tool is getting a very huge revamp that would bring back its legacy. ProxyBeast started two years ago as a side-project and it strived to be a high-performance and reliable proxy checker.
Since then, and now, my knowledge in computer science has been increased magnificently and for that reason I'd like to revise, which will make ProxyBeast great again.

#### Whats changing?
ProxyBeast is getting a revamp to be fined tuned in terms of speed, reliability, accuracy and efficiency. With additional features such as latency, anominity reporting and standardized multi protocol checking and the ability of providing complete PROXY URIs.
The core is completely changing to be very efficient and work in a parallel model, that is, lightweight and event-driven.
##### Native protocol checking
We built libraries that are native and have zero-dependencies, which will be ported with ProxyBeast. Meaning that ProxyBeast will be more low-level and for that reason it can improve accuracy as there is more deep-control.

**You can read more at:**
> https://proxy.pix4.dev
<br>
<br>


```
 .-,--.                   ,-,---.             .
  '|__/ ,-. ,-. . , . .    '|___/ ,-. ,-. ,-. |-
  ,|    |   | |  X  | |    ,|   \ |-' ,-| `-. |
  `'    '   `-' ' ` `-|   `-^---' `-' `-^ `-' `'
                     /|
                    `-'

        Tool by: @z3ntl3
        Studios: https://pix4.dev 

Usage:
   [flags]

Flags:
      --file string       Determines your proxy file name requires to be *.txt matching (default "proxies.txt")
  -h, --help              help for this command
      --multi             If passed as arg, it will check for all protocols, will tear down the accuracy
      --protocol string   The proxy protocol to check against (default "http")
      --retry int         The amount of tries to retry to connect to a failure proxy (default 2)
      --rotating          If passed, it ill use backbone check mechanism
      --timeout int       Sets custom timeout in seconds (default 5)
```

[Proof](https://www.youtube.com/shorts/TaDn6wtKqSk)
<table><tr><th>Organisation</th><th>Application</th><th>Developer</th></tr><tr><td><img src="https://media.discordapp.net/attachments/956310840464773200/968964843333877830/logopix4.png" width="20">PIX4</td><td>Proxy Beast</td><td>Z3NTL3</td></tr></table>

# Protocols Supported
`HTTPS, HTTP, SOCKS4, SOCKS5`

> **Update note:**<br>
> In the previous version there were many bugs and socks4/socks5 check was dirty,
with this new version we have enhanced these things and now it is precise and accurate in checking all protocols!

# Format

**Important**<br>
Only use `ip:port` format do not include `http://`. `socks5://`,`https://`, `socks4://` etc on your proxy file!

# Usage
<a href="[https://github.com/Z3NTL3/ProxyBeast#saves](https://github.com/Z3NTL3/ProxyBeast?tab=readme-ov-file#how-to-compile)">How to compile (build) instructions</a><br>
`chmod 755 proxy-checker.exe`<br>

#### Usage:

```
 .-,--.                   ,-,---.             .
  '|__/ ,-. ,-. . , . .    '|___/ ,-. ,-. ,-. |-
  ,|    |   | |  X  | |    ,|   \ |-' ,-| `-. |
  `'    '   `-' ' ` `-|   `-^---' `-' `-^ `-' `'
                     /|
                    `-'

        Tool by: @z3ntl3
        Studios: https://pix4.dev 

Usage:
   [flags]

Flags:
      --file string       Determines your proxy file name requires to be *.txt matching (default "proxies.txt")
  -h, --help              help for this command
      --multi             If passed as arg, it will check for all protocols, will tear down the accuracy
      --protocol string   The proxy protocol to check against (default "http")
      --retry int         The amount of tries to retry to connect to a failure proxy (default 2)
      --rotating          If passed, it ill use backbone check mechanism
      --timeout int       Sets custom timeout in seconds (default 5)
```

`<bin> -h` to see all options

#### Example

`./proxy-checker.exe --timeout 15 --retry 2 --protocol socks4`<br>

# Saves

Good working proxies are saved in the directory `/saves`. Each time running the script it will recreate the `goods.txt` corresponding for the recently checked proxies.

# How to compile

One-time run:
`go run .`

After that run one-time:
`go build` to compile everything so you can have an executable file.

Now your installation is done, just run the executable and there u go

### How to install Go

Install Go `minimum Go version: 1.19`

Navigate to `https://go.dev/dl/` install the one you need compabitle with your OS.<br>
`If you are on Windows you do not need to follow the instructions bellow.`

```
# Installation
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.19.2.linux-amd64.tar.gz

// you may need to run the above command as sudo, if you do change alse the $HOME/.profile of root with source $HOME/.profile

export PATH=$PATH:/usr/local/go/bin
source $HOME/.profile

go version // if you get output it works. Do not forget to follow the last 2 steps on differents users on your machine
```

# Update Log
- Small improvements
- Enhanced and fixed many bugs. SOCKS version 4 and 5 have been fixed and the check is now very precise and accurate (28-10-2023)

- Arg flags and retry option added (`23 jan 2023`) -> `https://github.com/Z3NTL3/ProxyBeast/pull/2`
- Multi option added (`3 feb 2023`) -> `https://github.com/Z3NTL3/ProxyBeast/issues/5` ``Preview: https://www.youtube.com/watch?v=_7J9u3EvA7k``
