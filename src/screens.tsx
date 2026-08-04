import { AiOutlineDashboard } from "react-icons/ai";
import { BsCashStack } from "react-icons/bs";
import { IoSettingsOutline } from "react-icons/io5";

export const SCREENS = [
  {
    title: "Overview",
    href: "/",
    node: (
      <>
        <AiOutlineDashboard size={20} />
        Overview
      </>
    ),
  },
  {
    title: "Settings",
    href: "/settings",
    node: (
      <>
        <IoSettingsOutline size={20} />
        Settings
      </>
    ),
  },

  {
    title: "Credits",
    href: "/credits",
    node: (
      <>
        <BsCashStack size={20} />
        Credits
      </>
    ),
  },
];
