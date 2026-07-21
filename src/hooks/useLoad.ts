import { emitTo } from "@tauri-apps/api/event";
import { useEffect, useState } from "react";

export default function useLoad(to?: string): Boolean {
  let [domLoaded, setDomLoaded] = useState(false);
  useEffect(() => {
    if (!domLoaded) setDomLoaded(true);
    if (typeof to !== "undefined") emitTo(to, "window_loaded");
    else emitTo("main", "window_loaded");
  }, [domLoaded]);

  return domLoaded;
}
