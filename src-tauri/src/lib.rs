use async_channel::{Receiver, Sender, bounded};
use chrono::Local;
use proxifier_rs::{ClientConfig, RootCertStore};
use std::fs;
use std::io::stdout;
use std::sync::atomic::AtomicBool;
use std::sync::{Arc, OnceLock};
use std::time::Duration;
use tauri::{AppHandle, Emitter, Listener, Manager, async_runtime};
use tokio::sync::RwLock;
use tracing::field::Visit;
use tracing::{Subscriber, info};
use tracing_appender::non_blocking::WorkerGuard;
use tracing_subscriber::filter::filter_fn;
use tracing_subscriber::prelude::*;
use tracing_subscriber::{Layer, Registry, fmt};

static LIVE_LOGS: &str = "live-logs";
static APP_HANDLE: OnceLock<AppHandle> = OnceLock::new();

struct AppState {
    proxy_checker: ProxyChecker,
    tls_config: Arc<ClientConfig>,
    app_config: RwLock<models::AppConfig>,
    #[allow(unused)]
    log_guard: WorkerGuard,
}

struct ProxyChecker {
    fd_state: AtomicBool,
    workers_state: AtomicBool,
    signal: RwLock<CancellationToken>,
    pipe: (Sender<String>, Receiver<String>),
}

pub(crate) use tokio_util::sync::CancellationToken;

use crate::models::AppConfig;

pub(crate) mod commands;
pub(crate) mod models;
pub(crate) mod events {
    pub const WINDOW_LOADED: &str = "window_loaded";
    pub const WINDOW_LOAD_PROGRESS: &str = "load_progress";
    pub const APP_VERSION: &str = "app_version";
}

struct LiveLogs;
impl<S: Subscriber> Layer<S> for LiveLogs {
    fn enabled(
        &self,
        _: &tracing::Metadata<'_>,
        _: tracing_subscriber::layer::Context<'_, S>,
    ) -> bool {
        true
    }

    // avoid unrelated events
    fn event_enabled(
        &self,
        _event: &tracing::Event<'_>,
        _ctx: tracing_subscriber::layer::Context<'_, S>,
    ) -> bool {
        _event.metadata().target() == LIVE_LOGS
    }

    fn on_event(&self, event: &tracing::Event<'_>, _: tracing_subscriber::layer::Context<'_, S>) {
        let mut visitor: MessageVisitor = Default::default();
        event.record(&mut visitor);

        let t: chrono::DateTime<Local> = Local::now();
        let log = format!("[{}] {}", t.format("%d/%m/%Y %H:%M:%S"), visitor.0);

        APP_HANDLE
            .get()
            .unwrap()
            .app_handle()
            .emit_to("main", "activity", log)
            .unwrap();
    }
}

#[derive(Default)]
#[repr(transparent)]
struct MessageVisitor(String);

impl Visit for MessageVisitor {
    fn record_debug(&mut self, field: &tracing::field::Field, value: &dyn core::fmt::Debug) {
        if field.name() == "message" {
            self.0 = format!("{:?}", value);
        }
    }
}

#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
    tauri::Builder::default()
        .plugin(tauri_plugin_fs::init())
        .plugin(tauri_plugin_process::init())
        .plugin(tauri_plugin_os::init())
        .plugin(tauri_plugin_dialog::init())
        .plugin(tauri_plugin_opener::init())
        .invoke_handler(tauri::generate_handler![
            commands::check_proxy_list,
            commands::stop_check,
            commands::read_file,
            commands::retrieve_settings,
            commands::save_settings
        ])
        .setup(|app| {
            let file_appender =
                tracing_appender::rolling::daily(app.path().app_log_dir()?, "diagnostics.log");
            let (non_blocking_appender, _guard) = tracing_appender::non_blocking(file_appender);

            let subscriber = Registry::default()
                .with(LiveLogs.with_filter(filter_fn(|metadata| metadata.target() == LIVE_LOGS)))
                .with(
                    fmt::Layer::default()
                        // $APPLOG file
                        .with_writer(non_blocking_appender)
                        .with_line_number(true)
                        .pretty()
                        .with_ansi(true),
                )
                .with(
                    fmt::Layer::default()
                        .with_writer(stdout)
                        .with_line_number(true)
                        .pretty()
                        .with_ansi(true),
                );
            tracing::subscriber::set_global_default(subscriber)
                .expect("failed setting default subscriber");

            let _ = std::fs::create_dir(app.path().app_config_dir()?);
            let config_file = app.path().app_config_dir()?.join("config.json");
            let mut app_config: models::AppConfig = AppConfig {
                timeout: Duration::from_millis(5000),
                pool_size: 1000,
                judge: "google.com".into(),
                enforce_scheme: crate::models::Scheme::Uri,
                use_tls: true,
                retry: true,
            };

            match fs::exists(config_file.clone()) {
                Ok(exists) => {
                    if !exists {
                        fs::write(config_file, serde_json::to_vec(&app_config)?).unwrap();
                    } else {
                        let config = fs::read_to_string(config_file).unwrap();
                        app_config = serde_json::from_str(&config[..]).unwrap();
                    }
                }
                Err(_) => {
                    fs::write(config_file, serde_json::to_vec(&app_config)?).unwrap();
                }
            };

            info!("App config {:?}", app_config);

            APP_HANDLE.set(app.app_handle().to_owned()).unwrap();
            let mut root_cert_store = RootCertStore::empty();
            root_cert_store.extend(webpki_roots::TLS_SERVER_ROOTS.iter().cloned());

            let config = Arc::new(
                ClientConfig::builder()
                    .with_root_certificates(root_cert_store)
                    .with_no_client_auth(),
            );

            let checker = ProxyChecker {
                fd_state: AtomicBool::new(false),
                workers_state: AtomicBool::new(false),
                signal: RwLock::new(CancellationToken::new()),
                pipe: bounded(app_config.pool_size as usize),
            };

            app.manage(AppState {
                log_guard: _guard,
                proxy_checker: checker,
                tls_config: config,
                app_config: RwLock::new(app_config),
            });

            let windows = app.webview_windows();
            let main_window = windows.get("main").unwrap();

            main_window.hide()?;

            let boot_window = tauri::webview::WebviewWindowBuilder::new(
                app,
                "bootstrapper",
                tauri::WebviewUrl::App("/bootstrap".into()),
            )
            .center()
            .closable(false)
            .focused(true)
            .resizable(false)
            .title("Bootstrap")
            .maximizable(false)
            .inner_size(520_f64, 380_f64)
            .decorations(false)
            .visible(false)
            .always_on_top(true)
            .build()?;

            let bootstrap = boot_window.clone();
            let main = main_window.clone();

            boot_window.once(events::WINDOW_LOADED, move |_| bootstrap.show().unwrap());
            main_window.once(events::WINDOW_LOADED, move |_| {
                let bootstrap_ = boot_window.clone();
                let main_ = main.clone();

                // tweak application loader so that progress bar is realtime.
                // FOR EXAMPLE: oneshot channel that measures every 200 milisecond and based on time duration
                // it calculates progress percentage, when reaching WINDOW_LOADED, it's calculated to a 1:1 ratio in 0-100%
                async_runtime::spawn(async move {
                    for num in [
                        rand::random_range(0.0..=0.2),
                        rand::random_range(0.2..0.4),
                        rand::random_range(0.4..0.8),
                    ] {
                        bootstrap_
                            .emit_to("bootstrapper", events::WINDOW_LOAD_PROGRESS, num)
                            .unwrap();
                        tokio::time::sleep(Duration::from_millis(200)).await;
                    }

                    tokio::time::sleep(Duration::from_millis(100)).await;
                    bootstrap_
                        .emit_to("bootstrapper", events::WINDOW_LOAD_PROGRESS, 1.0)
                        .unwrap();

                    tokio::time::sleep(Duration::from_millis(200)).await;
                    bootstrap_.close().unwrap();
                    main_.show().unwrap();
                    main_.set_focus().unwrap();

                    APP_HANDLE
                        .get()
                        .unwrap()
                        .emit_to("main", events::APP_VERSION, env!("CARGO_PKG_VERSION"))
                        .unwrap();
                    info!(target: LIVE_LOGS ,"Application bootstrapped.");
                });
            });

            Ok(())
        })
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
