# ProxyBeast


Proxy Beast is a high-performance proxy checker that offers precise and rapid testing for HTTP, HTTPS, SOCK4, and SOCKS5 protocols. Its advanced capabilities make it an ideal choice for those who require reliable and efficient proxy testing.

> Please note that ProxyBeast will be revised very soon. Between 29-04 to 05-05 2024.
> 
> **What is going to change?**
>
> Performance, efficiency and reliability will be maximized with this upcoming patch. As it will not require the use of ``net/http`` (with third party libs), but a low-level client for quality proxy checking and reporting.
> You can expect additional features such as anonimity and latency reporting. This patch also does introduce a full optimization to efficiency, reliability and performance. The whole core of the tool will change.
> In the end it will be finalized as a enterpise-ready easy to use and modular tool for precisely checking proxy qualities.
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
