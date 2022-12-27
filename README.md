# ProxyBeast

Proxy Beast is a extreme precisely accurate and fast proxy checker supporting HTTP, HTTPS, SOCK4, SOCKS5.

<table><tr><th>Organisation</th><th>Application</th><th>Developer</th></tr><tr><td><img src="https://media.discordapp.net/attachments/956310840464773200/968964843333877830/logopix4.png" width="20">PIX4</td><td>Proxy Beast</td><td>Z3NTL3</td></tr></table>

# Proof

**Accuracy and speed proof**: <br>
<a href="https://www.youtube.com/watch?v=wuYibCOgoCA&feature=youtu.be"> Proof</a>

# Protocols Supported

`HTTPS, HTTP, SOCKS4, SOCKS5`

# Usage

<a href="https://github.com/Z3NTL3/ProxyBeast#saves">How to compile (build) instructions</a><br>
`./proxy-checker.exe <protocol> <file.txt> <timeoutSecond>`

Important here is that `file.txt` has always to be `*.txt`!<br>
`<protocol>` can only be one of `https,socks4, socks5, http`<br>
`<timeoutSecond>` is the timeout in seconds<br>

#### Example

`./proxy-checker.exe socks5 proxies.txt 10`

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

Navigate to `https://go.dev/dl/` install the one you need compabitle with your OS.

```
# Installation
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.19.2.linux-amd64.tar.gz

// you may need to run the above command as sudo, if you do change alse the $HOME/.profile of root with source $HOME/.profile

export PATH=$PATH:/usr/local/go/bin
source $HOME/.profile

go version // if you get output it works. Do not forget to follow the last 2 steps on differents users on your machine
```
