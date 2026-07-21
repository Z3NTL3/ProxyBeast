import { useContext, useEffect, useState } from "react";
import { ScreenContext } from "../screen.context";
import { LuSettings2 } from "react-icons/lu";
import { Switch } from "@/components/ui/switch";
import { SiTraefikproxy } from "react-icons/si";
import { Slider } from "@/components/ui/slider";
import { Badge } from "@/components/ui/badge";
import { Tooltip, TooltipContent, TooltipTrigger } from "./ui/tooltip";
import { motion } from "motion/react";
import { invoke } from "@tauri-apps/api/core";
import { toast } from "sonner";
import { relaunch } from "@tauri-apps/plugin-process";
import { CiFileOn } from "react-icons/ci";
import { openPath } from "@tauri-apps/plugin-opener";
import * as path from "@tauri-apps/api/path";
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "./ui/select";
import { AppSettings } from "@/@types/app";

const APPLOG_DIR = await path.appLogDir();
const JUDGES: Array<{ label: string; value: string }> = [
  {
    label: "Google",
    value: "google.com",
  },
  {
    label: "Cloudflare",
    value: "one.one.one.one",
  }
];

type SchemeKind = "URI" | "MULTI" | "HTTP" | "HTTPS" | "SOCKS4" | "SOCKS5"
const SCHEMES: Array<{ label: string; value: SchemeKind }> = [
  {
    label: "URI",
    value: "URI",
  },
  {
    label: "MULTI",
    value: "MULTI",
  },
  {
    label: "HTTP",
    value: "HTTP",
  },
  {
    label: "HTTPS",
    value: "HTTPS",
  },
  {
    label: "SOCKS4",
    value: "SOCKS4",
  },
  {
    label: "SOCKS5",
    value: "SOCKS5",
  },
];

const TOOLTIP_CONTENT_MAP: {[k in SchemeKind]: string} = {
  URI: "Expects URI format per line in your proxy list file",
  MULTI: "Scans for all protocols on every proxy",
  HTTP: "Scan only for HTTP protocol",
  HTTPS: "Scan only for HTTPS protocol",
  SOCKS4: "Scan only for SOCKS4 protocol",
  SOCKS5: "Scan only for SOCKS5 protocol"
}

export default function Settings() {
  let screen = useContext(ScreenContext);
  let [settings, setSettings] = useState<AppSettings>({
    poolSize: 1000,
    timeoutMS: 5000,
    judge: "google.com",
    scheme: "uri",
    use_tls: true,
    retry: true
  });
  let [poolSet, setPoolSet] = useState(false);
  console.log(settings)

  useEffect(() => {
    if (typeof screen.setData !== "undefined") {
      screen.setData((screen_) => {
        return {
          ...screen_,
          current: "Settings",
        };
      });
    }
  }, [screen.current !== "Settings"]);

  useEffect(() => {
    invoke("retrieve_settings")
      .then((payload) => {
        let cast = payload as typeof settings;
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

  const saveSettings = (restore?: boolean) => {
    invoke("save_settings", {
      payload: settings,
    })
      .then((_) => {
        if (poolSet) {
          toast(
            "Saved new settings must relaunch for new pool size to take effect",
            {
              action: {
                label: "Relaunch",
                onClick: async () => await relaunch(),
              },
            },
          );
          return;
        }
        toast.success(restore ? "Restored to defaults" : "Saved new settings");
      })
      .catch((err) => {
        console.error(err);
        toast.error("Something went wrong while saving the new settings");
      });
  };

  const restoreDefaults = () => {
    setSettings((_) => {
      return {
        poolSize: 1000,
        timeoutMS: 5000,
        judge: "google.com",
        scheme: "URI",
        use_tls: true,
        retry: true
      };
    });
    saveSettings(true);
  };

  return (
    <div className="flex flex-col p-5 mx-8 mt-10">
      <h1 className="text-2xl">Settings</h1>
      <p className="text-md text-gray-400">
        Configure your applications settings.
      </p>

      {/* settings pane */}
      <div className="grid grid-cols-2 gap-x-10 items-start justify-center mt-8">
        {/* pane item */}
        <div className="flex flex-col bg-[#2A2A45] w-[350px] h-fit rounded-md p-4 border border-white/20">
          <div className="flex items-center gap-x-1 mb-2">
            <LuSettings2 size={22} color="#4D6AF0" />
            <h2 className="text-[18px]">General</h2>
          </div>

          {/* sub item */}
          <div className="flex mt-2">
            <div className="flex flex-col">
              <h3 className="text-[14px] text-gray-300">Auto-Retry</h3>
              <p className="text-[12px] text-gray-400">
                Retry failed connections automatically.
              </p>
            </div>
            <div className="flex grow  justify-end items-center">
              <Switch onCheckedChange={((checked) => {
                setSettings((settings) => {
                  return {
                    ...settings,
                    retry: checked
                  }
                })
              })} checked={settings.retry} />
            </div>
          </div>
          {/* end */}

          {/* sub item */}
          <div className="flex mt-2">
            <div className="flex w-full flex-col">
              <h3 className="text-[14px] text-gray-300">Pool Size</h3>
              <p className="text-[12px] text-gray-400">
                Customize the worker pool size.
              </p>

              <div className="flex grow w-full justify-end items-end">
                <Badge className="text-[#0058BC] bg-[#D8E2FF] text-[12px]">
                  {settings.poolSize}
                </Badge>
              </div>
              <Slider
                onValueChange={(v) => {
                  setPoolSet(true);
                  setSettings((settings) => {
                    return { ...settings, poolSize: v as number };
                  });
                }}
                className="mt-2"
                value={[settings.poolSize]}
                min={100}
                max={10_000}
                step={10}
              />
            </div>
          </div>
          {/* end */}

          {/* sub item */}
          <div className="flex mt-5">
            <div className="flex flex-col">
              <h3 className="text-[14px] text-gray-300">Diagnostics</h3>
              <p className="text-[12px] text-gray-400">Application log file.</p>
            </div>
            <div className="flex grow  justify-end items-center">
              <motion.div
                onClick={() => openPath(APPLOG_DIR)}
                whileHover={{ borderWidth: 1, borderColor: "white" }}
                className="flex items-center gap-x-1 border border-white/60 text-xs cursor-pointer rounded-md px-5 py-1 hover:cursor-pointer"
              >
                <CiFileOn /> View
              </motion.div>
            </div>
          </div>
          {/* end */}

          <span className="underline mt-5 text-gray-400 text-[12px]">
            Must restart for new pool size to take effect.
          </span>
        </div>
        {/* end */}

        {/* pane item */}
        <div className="flex flex-col bg-[#2A2A45] w-[350px] h-fit rounded-md p-4 border border-white/20">
          <div className="flex items-center gap-x-1 mb-2">
            <SiTraefikproxy size={22} color="#4D6AF0" />
            <h2 className="text-[18px]">Connection</h2>
          </div>

          {/* sub item */}
          <div className="flex mt-2">
            <div className="flex flex-col w-full">
              <h3 className="text-[14px] text-gray-300">Timeout</h3>
              <p className="text-[12px] text-gray-400">
                Filter proxies within the maximum latency range.
              </p>
              <div className="flex grow w-full justify-end items-end">
                <Badge className="text-[#0058BC] bg-[#D8E2FF] text-[12px]">
                  {settings.timeoutMS}ms
                </Badge>
              </div>
            </div>
          </div>
          <Slider
            onValueChange={(v) => {
              setSettings((settings) => {
                return { ...settings, timeoutMS: v as number };
              });
            }}
            className="mt-2"
            value={[settings.timeoutMS]}
            max={20000}
            step={100}
          />

          {/* sub item */}
          <div className="flex flex-col mt-5">
            <div className="flex flex-col mb-2">
              <h3 className="text-[14px] text-gray-300">Protocol Scheme</h3>
              <p className="text-[12px] text-gray-400">Enforce certain schemes over your proxy list.</p>
            </div>

            <Select
              items={SCHEMES}
              onValueChange={(v, _) => {
                if (v === null || typeof v === "undefined") return;

                setSettings((settings) => {
                  return {
                    ...settings,
                    scheme: v as string,
                  };
                });
              }}
              multiple={false}
            >
              <SelectTrigger
                value={settings.scheme as string}
                className="w-[180px] text-gray-400"
              >
                <Tooltip>
                  <TooltipTrigger render={
                    <SelectValue placeholder={settings.scheme} />
                  } />
                  <TooltipContent>
                    {TOOLTIP_CONTENT_MAP[settings.scheme as SchemeKind]}
                  </TooltipContent>
                </Tooltip>
              </SelectTrigger>
              <SelectContent>
                <SelectGroup>
                  {SCHEMES.map((v, i) => (
                    <Tooltip>
                      <TooltipTrigger render={
                        <SelectItem
                          key={`select-item-${i}`}
                          disabled={settings.scheme === v.value ? true : false}
                          value={v.value}
                        >
                          {v.label}
                        </SelectItem>
                      } />
                      <TooltipContent>
                        {TOOLTIP_CONTENT_MAP[v.value]}
                      </TooltipContent>
                    </Tooltip>
                  ))}
                </SelectGroup>
              </SelectContent>
            </Select>
          </div>
          {/* end */}

          {/* sub item */}
          <div className="flex flex-col mt-5">
            <div className="flex flex-col mb-2">
              <h3 className="text-[14px] text-gray-300">Judge</h3>
              <p className="text-[12px] text-gray-400">Relay destination.</p>
            </div>

            <Select
              items={JUDGES}
              onValueChange={(v, _) => {
                if (v === null || typeof v === "undefined") return;

                setSettings((settings) => {
                  return {
                    ...settings,
                    judge: v as string,
                  };
                });
              }}
              multiple={false}
            >
              <SelectTrigger
                value={settings.judge as string}
                className="w-[180px]"
              >
                <SelectValue placeholder={settings.judge} />
              </SelectTrigger>
              <SelectContent>
                <SelectGroup>
                  {JUDGES.map((v, i) => (
                    <SelectItem
                      key={`select-item-${i}`}
                      disabled={settings.judge === v.value ? true : false}
                      value={v.value}
                    >
                      {v.label}
                    </SelectItem>
                  ))}
                </SelectGroup>
              </SelectContent>
            </Select>
          </div>
          {/* end */}

          {/* sub item */}
          <div className="flex mt-5">
            <div className="flex flex-col">
              <h3 className="text-[14px] text-gray-300">TLS</h3>
              <p className="text-[12px] text-gray-400">
                Wrap & wire connection with TLS.
              </p>
            </div>
            <div className="flex grow  justify-end items-center">
              <Switch disabled onCheckedChange={((checked) => {
                setSettings((settings) => {
                  return {
                    ...settings,
                    use_tls: checked
                  }
                })
              })} checked={settings.use_tls} />
            </div>
          </div>
          {/* end */}

          {/* end */}
        </div>
        {/* end */}
      </div>
      {/* end */}

      <div className="flex items-center grow gap-x-2 justify-end mt-10">
        <motion.div
          onClick={() => saveSettings()}
          whileHover={{
            scale: 1.04,
          }}
          className="bg-[#0A84FF] rounded-xl px-5 py-2 text-[14px] cursor-pointer"
        >
          Save Changes
        </motion.div>

        <div
          onClick={() => restoreDefaults()}
          className="cursor-pointer text-gray-400 text-[13px] hover:underline"
        >
          Restore defaults
        </div>
      </div>
    </div>
  );
}
