<div align="center">
  <img width="400" alt="proxybeast gif" src="https://github.com/user-attachments/assets/3bd8a58e-958c-42ac-b0ae-6f7149d495e2" />

</div>

# Proxybeast

ProxyBeast is a powerful, complete and free proxy checker with advanced capabilities.

> This branch will be merged with the upstream repo on release

- [Trailer](https://www.youtube.com/watch?v=GOW_JKMfr9U)

> [!NOTE]  
> WRITTEN BY REAL HUMAN BEINGS

#### Features
- High-performance
- Reliable
- Rapid scanning
- Support for all proxy protocols (HTTP/HTTPS, SOCKS4/SOCKS5)
- Proxy Statistics
  - Latency
- Authentication Methods
  - ``User/Pass`` and ``NoAuth`` negotiations supported
- Easy to use
- Small & minimal app size
  

> Check out our propietrary proxy client crate we made for production use with ProxyBeast-v2 <br>
> [proxifier-rs crate](https://github.com/z3ntl3/proxifier-rs)


#### Awesome ProxyBeast-v2

##### Did you know?
- ProxyBeast has a ``LiveLog`` pane in the GUI, which shows all the neccessary logs. This was realised using ``Tokio-Tracing``. Specifically by composing multiple [Layer](https://docs.rs/tracing-subscriber/latest/tracing_subscriber/layer/index.html)'s together. One [Layer](https://docs.rs/tracing-subscriber/latest/tracing_subscriber/layer/index.html)'s responsibility is to watch log event data and send it efficiently to the GUI.

<img width="200" alt="image" src="https://github.com/user-attachments/assets/5dd2861a-fd6f-4b95-adbe-c2f5f1a90e4d" />


- ProxyBeast has it's own proxy client library implementation. It's based on ``CONNECT`` methods and allows for proxy relay to a destination target. Read more about it [here](https://github.com/z3ntl3/proxifier-rs)


- ProxyBeast preserves all application diagnostics, which roll daily to make it easier to identify problems. You can access those logs by pressing on the Diagnostics pane via the Settings application screen. We plan on making a bug tracker so we can patch bugs quicker, whereby you would attach the log file.
