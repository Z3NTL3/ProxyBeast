# ProxyBeast

Proxy Beast is a high-performance proxy checker that offers precise and rapid testing for HTTP, HTTPS, SOCK4, and SOCKS5 protocols. Its advanced capabilities make it an ideal choice for those who require reliable and efficient proxy testing.

**Accuracy and speed proof**: <br>
<a href="https://www.youtube.com/watch?v=RGzfBHpASZ8"> Proof</a>

<img src="https://media.discordapp.net/attachments/1071042419560296478/1071042467153068032/photo_2023-02-03_13-07-54.jpg?width=953&height=559">
<table><tr><th>Organisation</th><th>Application</th><th>Developer</th></tr><tr><td><img src="https://media.discordapp.net/attachments/956310840464773200/968964843333877830/logopix4.png" width="20">PIX4</td><td>Proxy Beast</td><td>Z3NTL3</td></tr></table>

# Protocols Supported

`HTTPS, HTTP, SOCKS4, SOCKS5`

# Format

**Important**<br>
Only use `ip:port` format do not include `http://`. `socks5://`,`https://`, `socks4://` etc on your proxy file!

# Usage

**Important**<br>

```
-multi option will consume alot of network. If your own network cannot respond by the load,
the proxy will be marked invalid! It requires high internet speed as this option is heavily trying all proxy protocols.
```

<a href="https://github.com/Z3NTL3/ProxyBeast#saves">How to compile (build) instructions</a><br>
`chmod 755 proxy-checker.exe`<br>

#### Usage:

```
 -file string
       Determines your proxy file name requires to be *.txt matching (default "proxies.txt")
 -multi
       If passed as arg, it will check for all protocols
 -protocol string
       Required flag, can be one of http, https, socks4 or socks5
 -retry int
       The amount of tries to retry to connect to a failure proxy (default 1)
 -timeout string
       Set custom timeout in seconds (default "5")
```

`./proxy-checker.exe -h` to see all options

#### Example

`./proxy-checker.exe -timeout 4 -retry 2 -protocol http`<br>
`./proxy-checker.exe -multi -retry 2`

# Saves

Good working proxies are saved in the directory `/saves`. Each time running the script it will recreate the `goods.txt` corresponding for the newly checked proxies.

### How to compile

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

- Arg flags and retry option added (`23 jan 2023`) -> `https://github.com/Z3NTL3/ProxyBeast/pull/2`
- Multi option added (`3 feb 2023`) -> `https://github.com/Z3NTL3/ProxyBeast/issues/5` ``Preview: https://www.youtube.com/watch?v=_7J9u3EvA7k``
