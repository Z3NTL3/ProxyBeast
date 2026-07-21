import { useEffect, useState } from "react";
import logo from "../assets/logo.png";
import useLoad from "../hooks/useLoad";
import { listen } from "@tauri-apps/api/event";
import { motion } from "motion/react";

export default function Bootstrap() {
  let [progress, setProgress] = useState<Number>(0);
  useLoad("bootstrapper");
  useEffect(() => {
    listen("load_progress", (progress) => {
      setProgress(progress.payload as Number);
    });
  });

  return (
    <div
      data-tauri-drag-region
      className="flex flex-col justify-center items-center h-screen w-screen bg-[#13131F] z-50"
    >
      <div className="px-4 py-3 bg-[#222236]/70 rounded-lg border border-[#808080]/30">
        <img
          width={40}
          src={logo}
          style={{
            boxShadow:
              "4px 4px 60px 10px rgba(0,88,180, 0.8), inset 4px 4px 60px 10px rgba(0,88,180, 0.3)",
          }}
        />
      </div>

      <h1 className="font-bold text-[26px] mt-7">ProxyBeast</h1>
      <p className="text-gray-500 text-[14px] text-xs">
        THE ULTIMATE PROXY CHECKER
      </p>

      <div className="flex w-full justify-center items-center">
        <div className="flex bg-gray-600 h-1 w-[60%] self-end mt-5 rounded-md">
          <motion.div
            initial={{ width: "0%" }}
            animate={{ width: `${Number(progress) * 100}%` }}
            className="bg-[#4D6AF0] h-1 self-end mt-5 rounded-md"
          ></motion.div>
        </div>
      </div>

      <p className="absolute bottom-5 justify-end text-gray-500 text-[12px] text-xs mt-10">
        SOFTWARE WRITTEN BY REAL HUMANS
      </p>
    </div>
  );
}
