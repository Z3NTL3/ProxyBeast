# ProxyBeast

Proxy Beast is a high-performance proxy checker that offers precise and rapid testing for HTTP, HTTPS, SOCK4, and SOCKS5 protocols. Its advanced capabilities make it an ideal choice for those who require reliable and efficient proxy testing.

**ProxyBeast Showcase**: <br>
`Click the image`
[![IMAGE_ALT](https://i.imgur.com/FB7xMyA.png)](https://www.youtube.com/watch?v=riaf6DqaIsI)

# Protocols Supported

- **SOCKS4**
- **SOCKS5**
- **HTTP**
- **HTTPS**

# Format
**Important**<br>
Only use ``ip:port`` format do not include ``http://``. ``socks5://``,``https://``, ``socks4://`` etc on your proxy file!


# Usage

First you will need to compile (for help see How to compile section)

 #### Usage:

 ```
  -file string
        Determine the name of the file containing proxies, must be in .txt format (default "proxies.txt")
  -protocol string
        Specify a required protocol, options include: http, https, socks4, and socks5
  -retries int
        Specify the number of attempts to reconnect to a failed proxy (default 1)
  -threads int
        Choose the number of threads to use for checking proxies, default is the number of CPU cores available (default 12)
  -timeout int
        Set the timeout in seconds using a custom value (default 5)
```
`./proxy-checker -h` to see all options

#### Example

- Check socks4 proxy with 10 second timeout and 3 retries<br>
`./proxy-checker --file proxies.txt -protocol socks4 -retries 3 -timeout 10`<br>

- Check socks4 proxy with 10 second timeout, 3 retries and with 12 threads<br>
`./proxy-checker --file proxies.txt -protocol socks4 -retries 3 -timeout 10 -threads 12`<br>

# Saves

Good working proxies are saved in the directory `/saves`. Each time running the script it will recreate the `goods.txt` corresponding for the newly checked proxies.

### How to compile

run `go build` to compile the tool from source

Now your installation is done, just run the executable and there u go

### How to install Go

Install Go `minimum Go version: 1.19`

Navigate to `https://go.dev/dl/` install the one you need compabitle with your OS.<br>

See `https://go.dev/doc/install` for installation instructions
