use crate::models::{self, MaybeTLS, Scheme};
use anyhow::anyhow;
use base64::prelude::*;
use fancy_regex::Regex;
use http::Uri;
use proxifier_rs::{
    Context, NetworkTarget, Port, ServerName, auth::Auth, http::tunnel as http, https::HttpsProxy,
    socks_with_tls, socks4::connect as socks4, socks5::connect as socks5,
};
use std::net::{SocketAddr, SocketAddrV4};
use std::sync::Arc;
use std::sync::atomic::AtomicU64;
use std::sync::atomic::Ordering::SeqCst;
use tauri::path::BaseDirectory;
use tauri::{AppHandle, Manager};
use tokio::io::{AsyncBufReadExt, AsyncWriteExt};
use tokio::net::lookup_host;
use tokio::task::JoinHandle;
use tokio::time::{Instant, timeout};
use tokio::{fs, select};
use tokio_util::sync::CancellationToken;
use tracing::{debug, error, info};
use url::Url;

#[derive(Debug)]
struct Ack {
    proxy: String,
    /// Latency in miliseconds [`core::time::Duration::as_millis`]
    latency: u128,
}

#[tauri::command]
pub async fn save_settings(app: AppHandle, payload: models::AppConfig) -> Result<bool, String> {
    let resource_path = app
        .path()
        .resolve("config.json", BaseDirectory::AppConfig)
        .map_err(|err| err.to_string())?;
    let ser = serde_json::to_vec(&payload).map_err(|err| err.to_string())?;
    tokio::fs::write(resource_path, &ser)
        .await
        .map_err(|err| err.to_string())?;

    let state = app.state::<crate::AppState>();
    let mut guard = state.app_config.write().await;
    *guard = payload;

    Ok(true)
}

#[tauri::command]
pub async fn retrieve_settings(app: AppHandle) -> Result<models::AppConfig, String> {
    let state = app.state::<crate::AppState>();
    let guard = state.app_config.read().await;

    Ok(guard.clone())
}

#[tauri::command]
pub async fn stop_check(state: tauri::State<'_, crate::AppState>) -> Result<(), ()> {
    if state.proxy_checker.workers_state.load(SeqCst) {
        state.proxy_checker.signal.read().await.cancel();
    }
    let recv = state.proxy_checker.pipe.1.clone();
    while !recv.is_empty() {
        let _ = recv.recv().await;
    }

    state.proxy_checker.fd_state.store(false, SeqCst);
    Ok(())
}

#[tauri::command]
pub async fn read_file(handle: AppHandle, path: String) -> Result<u16, String> {
    let state = handle.state::<crate::AppState>();
    let sender = state.proxy_checker.pipe.0.clone();

    let is_set = state.proxy_checker.fd_state.load(SeqCst);
    if is_set {
        info!("ongoing so not sending anything");
        return Err(anyhow!("Clear your proxy file before uploading a new one").to_string());
    }

    {
        let receiver = state.proxy_checker.pipe.1.clone();
        while !receiver.is_empty() {
            let _ = receiver.try_recv();
        }
    }

    let contents = fs::read(path).await;
    let mut load: u16 = 0;

    match contents {
        Ok(contents) => {
            state.proxy_checker.fd_state.store(true, SeqCst);
            let mut lines = contents.lines();
            while let Some(line) = lines.next_line().await.ok().flatten() {
                load += 1;

                let sender_clone = sender.clone();
                tokio::spawn(async move {
                    let _ = sender_clone.send(line).await;
                });
            }
        }
        Err(err) => {
            return Err(err.to_string());
        }
    }
    Ok(load)
}

/// Invoke this tauri command by [`check_proxy_list`] to initiate proxy checking, beaware that [`read_file`] must be called prior.
///
/// - [`read_file`] feeds proxy URIs non-blocking and asynchronously through a multi producer multi consumer [`async_channel::bounded`] channel.
/// - [`check_proxy_list`] consumes the feed after the user presses `Start Check` from GUI to invoke [`check_proxy_list`]
///
/// This command will exit instantly if no messages are found in the pipe. GUI alerts the user to provide a proxy list file.
///
/// # Procedures
/// Before starting any operation the following is asserted:
/// - If the given token was cancelled then the pipe should be drained out when not empty
/// - Subsequently the token should be renewed using [`tokio::sync::RWLock`] in [`crate::AppState`]
///
///
/// # Arguments
/// - timeout: Timeout in miliseconds for each task to complete, acts as max timeout for each proxy being asserted.
///
/// # Stop ongoing operation
/// Press `Stop Check` from the GUI, for more details read paragraph "Cancel Safety"
///
/// # Cancel Safety
/// - Safety via [`tokio_util::sync::CancellationToken>`]
///
/// This Tauri command is cancel safe. Invoke [`stop_check`] to gracefully terminate all ongoing operations on the multithreaded async runtime.
/// The user invokes this by pressing `Stop Check` button in the GUI.
#[tauri::command(rename_all = "snake_case")]
pub async fn check_proxy_list(
    app: AppHandle,
    chan: tauri::ipc::Channel<String>,
) -> Result<(), String> {
    let state = app.state::<crate::AppState>();
    let ongoing = state.proxy_checker.workers_state.load(SeqCst);
    let fd_state = state.proxy_checker.fd_state.load(SeqCst);
    if ongoing | !fd_state {
        if !fd_state {
            return Err("Upload a proxy list file".into());
        }

        return Err("Cannot upload new proxy file when there is ongoing operation.".into());
    }

    let d = state.app_config.read().await.timeout;
    info!("timeout set: {:?}", d);

    state.proxy_checker.workers_state.store(true, SeqCst);
    tokio::spawn(async move {
        let count = Arc::new(AtomicU64::new(0));
        let state = app.state::<crate::AppState>();
        let mut tasks: Vec<JoinHandle<()>> = vec![];

        info!("worker pool starting");
        let _ = chan.send("proxy-checker:start".into());

        let sender = state.proxy_checker.pipe.0.clone();
        let receiver = state.proxy_checker.pipe.1.clone();
        let token = state.proxy_checker.signal.read().await;

        // it's absolutely guaranteed that unwrap call will resolve to T
        // and will never fail
        for _ in 0..=sender.capacity().unwrap() {
            let count = count.clone();
            if token.is_cancelled() {
                break;
            }

            let app_clone = app.clone();
            let chan = chan.clone();
            let receiver = receiver.clone();
            let token = token.clone();

            let t = tokio::spawn(async move {
                let state = app_clone.state::<crate::AppState>();
                let config = state.app_config.read().await;

                't1: while !receiver.is_empty() && !token.is_cancelled() {
                    let proxy = receiver.try_recv();
                    if let Ok(mut proxy) = proxy {
                        let scheme = config.enforce_scheme;
                        // select the partitions after the protocol scheme
                        let re = Regex::new(
                            r"[^(socks4:\/\/)|(socks5:\/\/)|(http:\/\/)|(https:\/\/)](.*)",
                        );

                        if re.is_err() {
                            info!("is err");
                            continue 't1;
                        }

                        let alt_proxy = proxy.clone();
                        let re = re.unwrap();
                        match re.find(&proxy) {
                            Ok(Some(v)) => {
                                debug!("match {}", v.as_str());
                                proxy = v.as_str().to_owned();
                            }
                            Ok(None) => {
                                error!("regex didnt match break 't1");
                                let _ = chan.send("error|invalid-protocol-scheme".to_string());
                                continue 't1;
                            }
                            Err(err) => {
                                error!("{err} continuing 't1");
                                continue 't1;
                            }
                        };

                        let mut combine_scan: Vec<String> = vec![];
                        match scheme {
                            Scheme::Uri => combine_scan.push(alt_proxy),
                            Scheme::Multi => {
                                for proc in ["http", "https", "socks4", "socks5"] {
                                    combine_scan.push(format!("{}://{}", proc, proxy));
                                }
                            }
                            Scheme::Http => combine_scan.push(format!("{}://{}", "http", proxy)),
                            Scheme::Https => combine_scan.push(format!("{}://{}", "https", proxy)),
                            Scheme::Socks4 => {
                                combine_scan.push(format!("{}://{}", "socks4", proxy))
                            }
                            Scheme::Socks5 => {
                                combine_scan.push(format!("{}://{}", "socks5", proxy))
                            }
                        };

                        let tls_config = state.tls_config.clone();
                        debug!("combined proxy: {:?}", combine_scan);
                        for proxy in combine_scan {
                            let task = timeout(d, async {
                                let result: anyhow::Result<Ack> = async {
                                    let uri = proxy.parse::<Url>()?;
                                    let proxy_addr: SocketAddrV4 = format!("{}:{}",
                                        uri.host_str().ok_or_else(|| anyhow!("couldn't retrieve host portion"))?,
                                        uri.port().ok_or_else(|| anyhow!("couldn't retrieve port portion"))?
                                    ).parse()?;
                                    info!("proxy recv: {:?} addr: {proxy_addr}", proxy);

                                    let mut auth = Auth::NoAuth;
                                    let judge = &config.judge;
                                    #[allow(unused)]
                                    let retry = &config.retry;
                                    let tls = &config.use_tls;

                                    let mut host = {
                                        let judge = if *tls {
                                            judge.to_owned() + ":443"
                                        } else {
                                            judge.to_owned() + ":80"
                                        };

                                        lookup_host(judge)
                                            .await?
                                            .find_map(|addr| match addr {
                                                SocketAddr::V4(v4) => Some(v4),
                                                _ => None,
                                            })
                                            .ok_or_else(|| anyhow!("dns res didn't resolve to ipv4"))?
                                    };

                                    info!("host: {host}");
                                    info!("judge: {judge}");

                                    let now = Instant::now();
                                    match uri.scheme() {
                                        "http" => {
                                            if let Some(pass) = uri.password() {
                                                auth = Auth::HTTPAuthorizationHeader(format!("Proxy-Authorization: Basic {}", BASE64_STANDARD.encode(uri.username().to_owned() + ":" + pass)));
                                            }

                                            let dest = format!("http://{}:80", judge).parse::<Uri>()?;
                                            let mut conn = http(&dest, proxy_addr, auth).await?;

                                            conn.write_all(format!("GET / HTTP/1.1\r\nHost: {judge}:80\r\nConnection: close\r\n\r\n").as_bytes())
                                                .await?;

                                            info!("net ack")
                                        },
                                        "https" => {
                                            if let Some(pass) = uri.password() {
                                                auth = Auth::HTTPAuthorizationHeader(format!("Proxy-Authorization: Basic {}", BASE64_STANDARD.encode(uri.username().to_owned() + ":" + pass)));
                                            }

                                            let proxy = HttpsProxy::with_client_config(tls_config.clone());
                                            let mut conn = proxy.tunnel(&format!("https://{}:443", judge).parse()?, proxy_addr, auth).await?;

                                            conn.write_all(format!("GET / HTTP/1.1\r\nHost: {judge}:443\r\nConnection: close\r\n\r\n").as_bytes())
                                                .await?;

                                            info!("net ack")
                                        },
                                        "socks4" => {
                                            let conn: models::MaybeTLS = if *tls {
                                                host.set_port(443);
                                                MaybeTLS::Tls(Box::new(socks_with_tls(
                                                    socks4(
                                                        Context {
                                                            proxy: proxy_addr,
                                                            destination: host,
                                                        }
                                                    ).await?,
                                                    tls_config.clone(),
                                                    ServerName::try_from(judge.to_string())?
                                                ).await?))
                                            } else {
                                                host.set_port(80);
                                                MaybeTLS::Plain(
                                                    socks4(
                                                        Context {
                                                            proxy: proxy_addr,
                                                            destination: host,
                                                        }
                                                    ).await?
                                                )
                                            };

                                            match conn {
                                                MaybeTLS::Plain(mut conn) => {
                                                    conn.write_all(format!("GET / HTTP/1.1\r\nHost: {judge}:80\r\nConnection: close\r\n\r\n").as_bytes())
                                                        .await?;

                                                    info!("net ack")
                                                }
                                                MaybeTLS::Tls(mut conn) => {
                                                    conn.write_all(format!("GET / HTTP/1.1\r\nHost: {judge}:443\r\nConnection: close\r\n\r\n").as_bytes())
                                                        .await?;

                                                    info!("net ack")
                                                }
                                            }
                                        },
                                        "socks5" => {
                                            if let Some(pass) = uri.password() {
                                                auth = Auth::UserPass(uri.username().into(), pass.into())
                                            }

                                            let conn: models::MaybeTLS = if *tls {
                                                MaybeTLS::Tls(Box::new(socks_with_tls(
                                                    socks5(
                                                        Context {
                                                            proxy: proxy_addr,
                                                            destination: NetworkTarget::Domain(judge.parse()?, Port(443)),
                                                        },
                                                        auth,
                                                    ).await?,
                                                    tls_config.clone(),
                                                    ServerName::try_from(judge.to_string())?
                                                ).await?))
                                            } else {
                                                MaybeTLS::Plain(
                                                    socks5(
                                                        Context {
                                                            proxy: proxy_addr,
                                                            destination: NetworkTarget::Domain(judge.parse()?, Port(80)),
                                                        },
                                                        auth,
                                                    ).await?
                                                )
                                            };

                                            match conn {
                                                MaybeTLS::Plain(mut conn) => {
                                                    conn.write_all(format!("GET / HTTP/1.1\r\nHost: {judge}:80\r\nConnection: close\r\n\r\n").as_bytes())
                                                        .await?;

                                                    info!("net ack")
                                                }
                                                MaybeTLS::Tls(mut conn) => {
                                                    conn.write_all(format!("GET / HTTP/1.1\r\nHost: {judge}:443\r\nConnection: close\r\n\r\n").as_bytes())
                                                        .await?;

                                                    info!("net ack")
                                                }
                                            }
                                        }
                                        _  => {
                                            error!("skipping unknown proxy scheme '{:?}' in {:?}", uri.scheme(), uri);
                                        }
                                    }

                                    Ok(Ack { proxy: uri.to_string(), latency: now.elapsed().as_millis() })
                                }.await;
                                result
                            });
                            select! {
                                // prioritize cancellation
                                biased;
                                _ = token.cancelled() => {
                                     info!("task was cancelled");
                                     break 't1;
                                }
                                res = task => {
                                    match res {
                                        Ok(res) => {
                                            match res {
                                                Ok(ack) => {
                                                    info!("proxy:good:{}:latency:{}", ack.proxy, ack.latency);
                                                    let _ = chan.send(format!("proxy|good|{}|latency|{}", ack.proxy, ack.latency));
                                                }
                                                Err(err) => {
                                                    error!("proxy connection error: {err}");
                                                    info!("proxy:bad:{}", proxy);
                                                    let _ = chan.send(format!("proxy|bad|{}", proxy));
                                                }
                                            }
                                        }
                                        Err(err) => {
                                            error!("task aborted because it timed out: {:?}", err);
                                            let _ = chan.send(format!("proxy|bad|{}", proxy));
                                        }
                                    }

                                    count.fetch_add(1, SeqCst);
                                    info!("task finished");
                                }
                            };
                        }
                    }
                }
            });
            tasks.push(t);
        }

        for task in tasks {
            let _ = task.await;
        }

        info!("pool completed");
        drop(token);

        {
            let mut wlock = state.proxy_checker.signal.write().await;
            *wlock = CancellationToken::new();
        }

        // op done
        state.proxy_checker.fd_state.store(false, SeqCst);
        state.proxy_checker.workers_state.store(false, SeqCst);

        info!("worker pool OP completion");
        let _ = chan.send("proxy-checker:end".into());
        info!("COUNT: {}", count.load(SeqCst));
    });
    Ok(())
}
