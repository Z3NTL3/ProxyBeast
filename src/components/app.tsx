import { memo, useContext, useEffect, useRef, useState } from "react";
import { listen, UnlistenFn } from "@tauri-apps/api/event";
import { AiOutlineCloudUpload } from "react-icons/ai";
import { open, save } from "@tauri-apps/plugin-dialog";
import { TiMediaPlayOutline } from "react-icons/ti";
import { PiDownloadSimple } from "react-icons/pi";
import { invoke, Channel } from "@tauri-apps/api/core";
import { FaStop } from "react-icons/fa";
import { motion, useAnimate } from "motion/react";
import moment from "moment";
import "../App.css";
import { toast } from "sonner";
import { ScreenContext } from "@/screen.context";
import { Badge } from "./ui/badge";
import { RxFile } from "react-icons/rx";
import { writeTextFile } from "@tauri-apps/plugin-fs";
import { AppSettings } from "@/@types/app";
import { ContextMenu, ContextMenuContent, ContextMenuItem, ContextMenuTrigger } from "./ui/context-menu";
import { MdDeleteForever } from "react-icons/md";

const App = function () {
  let [logs, setLogs] = useState<
    Array<{
      timeFormat: string;
      msg: string;
    }>
  >([]);
  let [didStart, setDidStart] = useState(false);
  let [scope, animate] = useAnimate();
  let [id, setId] = useState<NodeJS.Timeout | null>(null);
  let [live_pane, dead_pane, progress] = [useRef(0), useRef(0), useRef(0)];
  let [load_pane, setLoad] = useState(0);
  // @ts-ignore
  let [proxies, setProxies] = useState<string[]>([]);
  let filePath = useRef("");

  let [settings, setSettings] = useState<AppSettings | null>(null);
  let screen = useContext(ScreenContext);

  useEffect(() => {
    invoke("retrieve_settings")
      .then((payload) => {
        let cast = payload as typeof settings;
        if(cast !== null)
          setSettings((_) => {
            return {
              ...cast,
            };
          });
      })
      .catch((err) => {
        toast.error(String(err));
      });
  }, []);

  useEffect(() => {
    if (typeof screen.setData !== "undefined") {
      screen.setData((screen_) => {
        return {
          ...screen_,
          current: "Overview",
        };
      });
    }
  }, [screen.current !== "Overview"]);

  useEffect(() => {
    const unlisten: Array<Promise<UnlistenFn>> = [];
    const activity = listen("activity", (ev) => {
      let segments = (ev.payload as string).split("]");
      let timeSeg = segments[0].replace("[", "");
      let msgSeg = segments[1];

      let log = {
        timeFormat: timeSeg,
        msg: msgSeg,
      };

      sessionStorage.setItem("boot", JSON.stringify(log));
      setLogs((logs) => [...logs, log]);
    });

    unlisten.push(activity);

    let bootAt = sessionStorage.getItem("boot");
    if (bootAt !== null && sessionStorage.getItem("bootSet") === null) {
      setLogs((logs) => [...logs, JSON.parse(bootAt)]);
      sessionStorage.setItem("bootSet", JSON.stringify(1));
    }

    return () => {
      unlisten.forEach(async (v) => v.then((cleanup) => cleanup()));
    };
  }, []);

  const startChecker = async () => {
    animate("#checker-btn", {
      rotate: [0, 360],
    });

    if (didStart) {
      await invoke("stop_check");
      return;
    }

    const channel = new Channel<string>();
    channel.onmessage = (message) => {
      switch (true) {
        case message === "proxy-checker:end": {
          filePath.current = "";
          setLoad(0);
          setDidStart((_) => false);
          setLogs((logs) => [
            ...logs,
            {
              timeFormat: moment().format("HH:mm:ss"),
              msg: "Proxy checker terminated",
            },
          ]);

          if (typeof id === "number") {
            clearInterval(id);
            setId(null);
          }
          break;
        }
        case message === "proxy-checker:start": {
          progress.current = 0;
          live_pane.current = 0;
          dead_pane.current = 0;

          let id = setInterval(() => {
            let log_pane = document.getElementById("pane");
            log_pane?.lastElementChild?.scrollIntoView(false);
          }, 24);

          setId(id);
          setDidStart((_) => true);
          setLogs((logs) => [
            ...logs,
            {
              timeFormat: moment().format("HH:mm:ss"),
              msg: "Proxy checker initialized",
            },
          ]);
          break;
        }
        case message.includes("proxy|good"): {
          live_pane.current += 1;
          calcProgress();
          let ack = message
            .split("proxy|good|")[1]
            .split("|")
            .filter((v) => v !== "");

          let proxy = ack[0];
          let latency = ack[2];

          setLogs((logs) => [
            ...logs,
            {
              timeFormat: moment().format("HH:mm:ss"),
              msg: `live: ${proxy} - latency ${latency}ms`,
            },
          ]);
          setProxies((list) => [...list, proxy]);
          break;
        }
        case message.includes("proxy|bad"): {
          dead_pane.current += 1;

          calcProgress();
          let proxy = message.split("proxy|bad|")[1];
          if (proxy.length === 0) return;

          setLogs((logs) => [
            ...logs,
            {
              timeFormat: moment().format("HH:mm:ss"),
              msg: `dead: ${proxy}`,
            },
          ]);
          break;
        }
      }
    };

    invoke("check_proxy_list", {
      chan: channel,
    }).catch((err) => {
      toast.error(String(err));
    });
  };

  const exportProxies = async () => {
    if (proxies.length < 1)
      toast.warning("No live proxies were found to be saved");
    let path = await save({
      filters: [
        {
          name: "Proxy file(txt)",
          extensions: ["txt"],
        },
      ],
    });

    if (path !== null)
      writeTextFile(path, proxies.join("\n"))
        .then(console.log)
        .catch((err) =>
          toast.error("Failed exporting proxies to given location"),
        );
  };

  const calcProgress = () => {
    let current_progress = live_pane.current + dead_pane.current;
    let p = (current_progress / load_pane ) * 100;
    progress.current = p;
  };

  return (
    <div className="w-full h-full overflow-hidden" ref={scope}>
      <div className="flex p-4 w-full h-fit">

        <ContextMenu>
          <div
            onClick={() => {
              open({
                multiple: false,
                directory: false,
                filters: [
                  {
                    name: "(txt) Proxy file",
                    extensions: ["txt"],
                  },
                ],
              }).then((path) => {
                if (typeof path === "string" && path.length > 1) {
                  invoke("read_file", { path })
                    .then((v) => {
                      setLoad(settings?.scheme === "MULTI" ? v as number * 4: v as number);
                      filePath.current = path;
                      toast.info("Selected proxy file");
                    })
                    .catch((err) => {
                      toast.error(String(err));
                    });
                }
              });
            }}
            className="cursor-pointer rounded-md flex flex-col justify-center items-center border-white/20  border-2 border-dotted w-[78%] h-120"
          >
            <div className="p-3 bg-[#DBDBFD] rounded-full mb-5">
              <AiOutlineCloudUpload color="#4D6AF0" size={30} />
            </div>

            <h3>Drop your proxy list here</h3>
            <p className="text-[13px] text-white/50">
              Press inside to select your proxy list file
            </p>
            {filePath.current !== "" ? (
              <ContextMenuTrigger>
                <motion.div
                  className="z-50"
                  whileInView={{
                    opacity: [0, 1],
                  }}
                >
                  <Badge className="mt-5 hover:text-white/60 bg-transparent text-white/40 border-white/10 px-5 h-fit py-1">
                    <RxFile size={20} className="text-white" /> Selected:{" "}
                    {filePath.current}
                  </Badge>
                </motion.div>
              </ContextMenuTrigger>
            ) : null}
          </div>
          <ContextMenuContent>
            <ContextMenuItem onClick={async () => {
              await invoke("stop_check")
              filePath.current = "";
              toast.info("Deleted attached proxy list file")
              }} className={"text-red-500"}><MdDeleteForever /> Delete</ContextMenuItem>
          </ContextMenuContent>
        </ContextMenu>
        <div className="flex flex-col">
          <div className="flex flex-col items-start justify-start bg-[#2A2A45] p-2 w-50 h-50 ml-2 rounded-md border border-white/10">
            <div className="flex items-center w-full text-white/40 text-xs">
              Total Loaded
              <div className="flex grow justify-end items-center mr-2">
                <h4 className="text-white text-lg font-semibold">
                  {load_pane}
                </h4>
              </div>
            </div>

            <div className="flex w-full items-center gap-x-2 mt-7">
              <div className="bg-green-500 rounded-full w-2 h-2"></div>

              <p className="text-xs text-white/40">Live</p>

              <h2 className="flex grow items-center justify-end mr-1 text-green-600 font-semibold text-2xl">
                {live_pane.current}
              </h2>
            </div>

            <div className="flex w-full items-center gap-x-2 mt-1">
              <div className="bg-red-500 rounded-full w-2 h-2"></div>

              <p className="text-xs text-white/40">Dead</p>

              <h2 className="flex grow items-center justify-end mr-1 text-red-600 font-semibold text-2xl">
                {dead_pane.current}
              </h2>
            </div>

            <div className="flex flex-col w-full grow justify-end">
              <div className="flex items-center p-1">
                <h2 className="text-xs text-white/40">Progress</h2>

                <div className="flex w-full justify-end text-white/40 text-xs">
                  {progress.current.toFixed(0)}%
                </div>
              </div>
              <div className="border border-white/20 bg-[#2A2A45] w-full h-2 rounded-full">
                <motion.div
                  layout
                  initial={{
                    width: "0%",
                  }}
                  animate={{
                    width: `${progress.current.toFixed(0)}%`,
                  }}
                  className={`bg-blue-500 h-full rounded-md`}
                ></motion.div>
              </div>
            </div>
          </div>

          <div
            onClick={startChecker}
            className={`${!didStart ? "hover:shadow-[1px_1px_30px_0.1px_rgba(53,120,236,0.4)]" + " bg-[#0A84FF]" : "hover:shadow-[1px_1px_30px_0.1px_rgba(255,68,2,0.3)]" + " bg-red-500"} cursor-pointer flex items-center gap-x-1 text-sm justify-center w-full ml-2 rounded-lg p-4 mt-5 text-center font-semibold`}
          >
            <motion.div
              id="checker-btn"
              whileInView={{ rotate: [0, 360] }}
              layout
            >
              {!didStart ? (
                <TiMediaPlayOutline size={20} />
              ) : (
                <FaStop size={20} />
              )}
            </motion.div>
            {!didStart ? "Start Check" : "Stop Check"}
          </div>

          <div
            onClick={() => exportProxies()}
            className="hover:shadow-[1px_1px_30px_0.1px_rgba(255,255,255,0.02)] cursor-pointer flex flex-col gap-y-1 items-center justify-center border border-white/10 text-center rounded-lg ml-2 mt-4 p-2 text-white/40 text-xs"
          >
            <PiDownloadSimple className="text-white/70 font-bold" size={18} />
            Export Proxies
          </div>
        </div>
      </div>

      <div className="flex flex-col w-[76%] p-2 rounded-md mx-4 h-fit overflow-hidden mb-5 bg-[#2A2A45] border border-white/10 mt-2">
        <div className="flex items-center w-full h-fit">
          <h2 className="font-semibold font-inter">Live Logs</h2>

          <div className="flex grow justify-end items-center mr-2">
            <div className="bg-blue-500 rounded-full w-1 h-1 animate-pulse"></div>
          </div>
        </div>

        <div
          id="pane"
          className="w-full h-20 overflow-y-scroll flex flex-col scrollbar-thumb-gray-800"
        >
          {logs.map((log, i) => (
            <div key={`logview-item-${i}`} className="flex">
              <p className="text-white/40">
                <span className="text-[#4D6AF0]/85">[{log.timeFormat}]</span>
                <span className="ml-1">
                  {log.msg.includes("live") || log.msg.includes("dead") ? (
                    <span
                      className={`${log.msg.includes("live") ? "text-green-300" : "text-red-400"}`}
                    >
                      {log.msg}
                    </span>
                  ) : (
                    log.msg
                  )}
                </span>
              </p>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};
export default memo(App);
