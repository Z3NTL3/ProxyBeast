import { createContext, Dispatch, SetStateAction } from "react";
export interface ScreenData {
  version: String;
  current?: String;
  setData?: Dispatch<SetStateAction<ScreenData>>;
}

export const ScreenContext = createContext<ScreenData>({
  version: "1.0.0",
});
