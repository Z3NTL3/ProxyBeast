use proxifier_rs::{TcpStream, TlsStream};
use serde::{Deserialize, Serialize};
use serde_with::{DisplayFromStr, DurationMilliSeconds, serde_as};
use std::time::Duration;

#[derive(Serialize, Deserialize, Debug, Clone, Copy)]
pub enum Scheme {
    Uri,
    Multi,
    Http,
    Https,
    Socks4,
    Socks5,
}

pub enum MaybeTLS {
    Plain(TcpStream),
    Tls(Box<TlsStream<TcpStream>>),
}

impl std::fmt::Display for Scheme {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            Scheme::Uri => write!(f, "URI"),
            Scheme::Multi => write!(f, "MULTI"),
            Scheme::Http => write!(f, "HTTP"),
            Scheme::Https => write!(f, "HTTPS"),
            Scheme::Socks4 => write!(f, "SOCKS4"),
            Scheme::Socks5 => write!(f, "SOCKS5"),
        }
    }
}

impl core::str::FromStr for Scheme {
    type Err = anyhow::Error;

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        Ok(match s {
            "URI" => Scheme::Uri,
            "MULTI" => Scheme::Multi,
            "HTTP" => Scheme::Http,
            "HTTPS" => Scheme::Https,
            "SOCKS4" => Scheme::Socks4,
            "SOCKS5" => Scheme::Socks5,
            _ => Scheme::Uri,
        })
    }
}

#[serde_as]
#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct AppConfig {
    #[serde(rename(serialize = "poolSize", deserialize = "poolSize"))]
    pub pool_size: u64,

    #[serde_as(as = "DurationMilliSeconds<u64>")]
    #[serde(rename(serialize = "timeoutMS", deserialize = "timeoutMS"))]
    pub timeout: Duration,

    // Fields with defaults
    #[serde(default = "default_judge")]
    pub judge: String,

    #[serde_as(as = "DisplayFromStr")]
    #[serde(
        default = "default_scheme",
        rename(serialize = "scheme", deserialize = "scheme")
    )]
    pub enforce_scheme: Scheme,

    #[serde(default = "default_tls")]
    pub use_tls: bool,
    #[serde(default = "default_retry")]
    pub retry: bool,
}

fn default_judge() -> String {
    "google.com".into()
}

fn default_scheme() -> Scheme {
    Scheme::Uri
}

fn default_tls() -> bool {
    true
}

fn default_retry() -> bool {
    true
}
