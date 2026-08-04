import React, { memo, useContext } from "react";
import { MdOutlineClose } from "react-icons/md";
import { MdMinimize } from "react-icons/md";
import { Window } from "@tauri-apps/api/window";
import { Outlet } from "react-router";
import { motion } from "motion/react";
import { ScreenContext } from "./screen.context";
import { SCREENS } from "./screens.tsx";
import { invoke } from "@tauri-apps/api/core";
import { Toaster } from "@/components/ui/sonner";
import "./App.css";
import logo from "./assets/logo.png";
import { TooltipProvider } from "@/components/ui/tooltip";
import { CiStickyNote } from "react-icons/ci";
import { platform } from "@tauri-apps/plugin-os";
import useLoad from "./hooks/useLoad.ts";
const PLATFORM = platform();

export const Layout = memo(() => {
  useLoad();
  let screenData = useContext(ScreenContext);
  return (
    <TooltipProvider>
      <div className="flex w-screen h-screen bg-[#1E1E2E] overflow-hidden">
        <motion.div
          initial={{
            opacity: 0,
          }}
          whileInView={{ opacity: 1 }}
          layout
          className="bg-[#2A2A45] w-60 h-full p-5 border-r border-[#808080]/40"
        >
          {/*logo*/}
          <div
            data-tauri-drag-region
            className={`${PLATFORM === "macos" ? "mt-2" : null} flex items-center`}
          >
            <img width={40} src={logo} />
            <h2 className="text-[24px] font-bold">ProxyBeast</h2>
          </div>
          {/*end*/}

          {/*version*/}
          <p className="text-xs text-white/40 mt-2 ml-1">
            v{sessionStorage.getItem("version")}
          </p>
          {/*end*/}

          {/*items*/}
          <div className="flex flex-col gap-y-4 mt-10">
            {SCREENS.map((screen_, i) => (
              <React.Fragment key={`screen-${i}`}>
                <a
                  tabIndex={-1}
                  key={`screen-${i}`}
                  href={screen_.href}
                  onClick={async () => {
                    await invoke("stop_check");
                    return true;
                  }}
                >
                  {screenData.current === screen_.title ? (
                    <motion.div
                      initial={{
                        scaleY: 0,
                      }}
                      whileInView={{ scaleY: 1 }}
                      className="hover:shadow-[1px_1px_30px_0.1px_rgba(53,120,236,0.4)] bg-[#0A84FF] cursor-pointer flex items-center gap-x-2 px-2 py-1.5 rounded-lg font-inter text-[13px] text-left"
                    >
                      {screen_.node}
                    </motion.div>
                  ) : (
                    <div className="hover:text-white/60 text-white/40 flex items-center gap-x-2 px-2 py-1.5 cursor-pointer rounded-lg font-inter text-[13px] text-left">
                      {screen_.node}
                    </div>
                  )}
                </a>
              </React.Fragment>
            ))}
          </div>
          {/*end*/}

          <div className={`${PLATFORM === "macos" ? "mt-95" : "mt-100"}`}>
            <div className="flex w-full h-[1px] bg-white/20 rounded-full"></div>
            <p className="mt-3 flex items-center justify-center gap-x-1 text-center text-[12px] text-white/90">
              <CiStickyNote size={15} />
              Humanly Engineered
            </p>
          </div>
        </motion.div>
        <main className={`flex flex-col w-full h-full items-start`}>
          <motion.div
            data-tauri-drag-region
            layout
            initial={{
              opacity: 0,
            }}
            whileInView={{ opacity: 1 }}
            className="flex bg-[#2A2A45] w-full h-10 p-3 font-inter text-[13px] text-white/80 border-b border-[#808080]/40"
          >
            {screenData.current}
            <div
              data-tauri-drag-region
              className={`${PLATFORM === "macos" ? "hidden" : "flex grow items-center justify-end gap-x-1"}`}
            >
              <MdMinimize
                className="font-bold"
                fontSize={18}
                onClick={() => Window.getCurrent().minimize()}
              />
              <MdOutlineClose
                className="font-bold"
                fontSize={18}
                onClick={() => Window.getCurrent().close()}
              />
            </div>
          </motion.div>

          <motion.div
            className="w-full h-full"
            layout
            initial={{
              opacity: 0,
            }}
            whileInView={{ opacity: 1 }}
            transition={{
              type: "keyframes",
              ease: "easeIn",
            }}
          >
            <Outlet />
            <Toaster />
          </motion.div>
        </main>
      </div>
    </TooltipProvider>
  );
});
