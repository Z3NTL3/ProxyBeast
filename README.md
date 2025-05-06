<!-- header -->

<div align="center">   
    <div>
        <img src="https://simpaix.net/img/logo.png" width=80><br>
        <span>A project by SimpaiX</span> <br><br>
         <div>
                <img alt="GitHub License" src="https://img.shields.io/github/license/z3ntl3/ProxyBeast" >
                <img alt="GitHub Repo stars" src="https://img.shields.io/github/stars/z3ntl3/ProxyBeast">
                <img alt="GitHub go.mod Go version" src="https://img.shields.io/github/go-mod/go-version/z3ntl3/ProxyBeast">
        </div>
        <a href="https://github.com/Z3NTL3/ProxyBeast/releases/latest">Download ProxyBeast for Windows</a> <br>  
    </div>
    <img src="https://z3ntl3.com/img/proxybeast.png" width="400" style="border-radius: 4px;"><br>
</div>

<!-- intro -->
> [!IMPORTANT]
> WASM (web) variant is coming when WASM has threading support. So that the unique threadpool controller we architected for ProxyBeast can be translated into WASM, otherwise it would run on one thread which is not what we seek for. 

# ProxyBeast 


ProxyBeast is a powerful, complete and free proxy checker with [zero dependency](#what-do-you-mean-with-zero-dependency)
and advanced capabilities.

> [!NOTE]
> Start using ProxyBeast. Choose between installing from an installer or build an executable from source.<br>[Get Started](#get-started)

### Features
- Lightweight
- High performance
- Event-driven
- Rich ecosystem
- Swift

- #### Capabilities
    - *Multi protocol checking*
        > Can check all protocols at once
    - *Supports proxy checking for*
        > ``SOCKS/4/5 & HTTP/HTTPS`` type proxies<br>
        > - **NOTE**<br>
        > SOCKS protocol version 4/a can be supported. Request it from Github issues tab.
    - *Powerful event-driven goroutine pools*
        > Results in efficient and reliable architecture
    - *Lightweight app*
        > Minimizing overhead, maximizing performance
    - *Recognizes URI patterns*
        > Essential for multi protocol checking

#### File format
Below examples of correct file format:
```
user:pass@host:port
scheme://host:port
scheme://user:pass@port
host:port
```
> Scheme can only be of ``http|https|socks4|socks5``<br>
> [Additional documentation](https://pkg.go.dev/net/url#URL)

### Get Started

We are very happy because of your interest in ProxyBeast. This guide is dedicated
to help you setup Proxybeast. 

Do not worry, our software is easy to use, so we won't be covering usage tutorials here. You may find tutorials on our [YouTube](https://www.youtube.com/@z3ntl3wip) channel.


- ### Installation
    
    There are two possible ways to install our software. To build from scratch, using ``Go`` tools. Or to install from a packaged installer.
  
    > We have a CI/CD to automatically deploy changes made to Github packages.
    > Additonal information can be found on the packages page.

    #### Precompiled installers and or executables

    | Platform      | File | Type |
    | ----------- | ----------- | ----------- |
    | Windows 10/11 (windows/amd64)      | [Installer](https://github.com/Z3NTL3/ProxyBeast/releases/download/v1.0.0/ProxyBeast-amd64-installer.exe)       | Windows installer |

    > Installers and or bins for other platforms aren't available, you can build from scratch
    > or await upon completion of the WASM (web) variant of ProxyBeast

    #### Build from scratch

    We assume you already have Go and it's toolchains installed. If not, follow the steps on this page.
    > [Install Go](https://go.dev/doc/installhttps://go.dev/doc/install)

    ##### Clone this repository
    We use ``git`` to clone this repo. This should download ProxyBeast in the current folder location.
  
    ```
    git clone https://github.com/Z3NTL3/ProxyBeast
    ```

    ##### Navigate into ProxyBeast workspace
    In the previous step we did install ProxyBeast. Execute the following to navigate into the folder.
    ```
    cd ProxyBeast
    ```

    ### Building
    
    For building, it's important to first, install all the dependencies of ProxyBeast. Execute the snippet below to continue.
    > We assume that you are located in the ProxyBeast workspace (better said, projects folder), as of the previous step.
    ```
    go get .
    ```

    - ##### Installating required tools
    First we need to install Wails. For this execute the following command.
    ```
    go install github.com/wailsapp/wails/v2/cmd/wails@latest
    ```
    > - **Important note**<br>
    > You require to have NPM installed. Which fits with Node.js. To install follow given directions here
    > [Install](https://nodejs.org/en)

    - ##### Validating tools
    To validate that you're all set, execute the following command. If everything seems OK, go further with the last procedures.
    ```
    wails doctor
    ```
    - ##### Final step 
    To build a binary (generating executables from source) execute the following:
    ```
    wails build
    ```
    > **NOTE**<br>
    > If you want to build an executable with GUI and a terminal for logs, execute the following command instead:
    > ```
    > wails build -windowsconsole
    >```
    > Attaching console to the GUI is only possible on Windows

<br>

> [!WARNING]
> Currently we do only support a packaged installer for Windows.<br>
> Installers for MacOS and debain based linux distributions aren't supported. It's essential noting there's no plan for this, as you can just build from scratch using our instructions guide. Which should be relatively easy to follow.

<br>
<hr>

# FAQ
Find an answer to most of your questions here. If it is not covered ask in [Discord](#todo).

* #### What do you mean with "zero dependency"
    With zero-dependency, we mean that ProxyBeast is ported together with a low level module for all of it's networking requirements. Which is a native module and is built without any additional dependency.
    
    > **Proxifier**<br>
    A module/library that is especially written to be ported with ProxyBeast, built by the author of ProxyBeast.<br>
    [Source](https://github.com/z3ntl3/Proxifier)

