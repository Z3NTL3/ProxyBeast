import { useContext, useEffect } from "react";
import { ScreenContext } from "../screen.context";
import { Badge } from "./ui/badge";
import { PiGithubLogoDuotone } from "react-icons/pi";
import { motion, stagger } from "motion/react";
import { Avatar, AvatarBadge, AvatarImage } from "./ui/avatar";
import { MdVerified } from "react-icons/md";
import { openUrl } from "@tauri-apps/plugin-opener";
import z3ntl3Pfp from "@/assets/img/z3ntl3.png";
import filipPfp from "@/assets/img/filip.png";
import { FaGithub, FaLinkedinIn, FaSquareXTwitter } from "react-icons/fa6";

export default function Credits() {
  let screen = useContext(ScreenContext);

  useEffect(() => {
    if (typeof screen.setData !== "undefined") {
      screen.setData((screen_) => {
        return {
          ...screen_,
          current: "Credits",
        };
      });
    }
  }, [screen.current !== "Credits"]);

  return (
    <div className="flex flex-col p-5 mx-8 mt-10  ">
      <h1 className="text-2xl">Contributors</h1>
      <p className="text-md text-gray-400">
        These are the people behind ProxyBeast.
      </p>

      <div className="grid grid-cols-2 justify-center mt-10 gap-x-8 gap-y-10">
        <div className="flex flex-col p-5 bg-[#2A2A45] border border-white/15 min-h-50 h-fit rounded-md">
          <div className="flex items-start gap-x-5 h-fit w-fit">
            <Avatar size={"lg"} className="mt-5">
              <AvatarImage src={z3ntl3Pfp} alt="@z3ntl3" />
              <AvatarBadge>
                <MdVerified />
              </AvatarBadge>
            </Avatar>
            <div className="flex flex-col">
              <div className="flex items-center gap-x-2 mt-3">
                <h2 className="font-semibold text-xl">Efdal Sancak</h2>
                <Badge className="text-xs">Core Developer</Badge>
              </div>
              <p className="text-white/60 font-light text-md">
                A passionate allround software engineer with deep passion for
                computer science. Admires the art of reading. Realised and built
                all parts of ProxyBeast.
              </p>

              <motion.div
                onClick={() => openUrl("https://github.com/z3ntl3")}
                whileHover={{
                  scaleX: 1.04,
                }}
                className="mt-3 bg-[#2A2A3D] items-center gap-x-2 text-[14px] w-fit px-10 text-left flex justify-start py-1 rounded-lg border border-white/15 cursor-pointer"
              >
                <PiGithubLogoDuotone /> Github
              </motion.div>
            </div>
          </div>
        </div>
        <div className="flex flex-col p-5 bg-[#2A2A45] border border-white/15 min-h-50 h-fit rounded-md">
          <div className="flex items-start gap-x-5 h-fit w-fit">
            <Avatar size={"lg"} className="mt-5">
              <AvatarImage
                src={filipPfp}
                className="grayscale"
                alt="@terzicdsgn"
              />
              <AvatarBadge>
                <MdVerified />
              </AvatarBadge>
            </Avatar>
            <div className="flex flex-col">
              <div className="flex items-center gap-x-2 mt-3">
                <h2 className="font-semibold text-xl">Filip Terzic</h2>
                <Badge className="text-xs">Lead Designer</Badge>
              </div>
              <p className="text-white/60 font-light text-md">
                Talented UI/UX youngster. Focused on elevating creative designs
                to higher levels of profound designs. Contributed Figma
                prototypes for ProxyBeast.
              </p>

              <motion.div
                onClick={() => openUrl("https://x.com/terzicdsgn")}
                whileHover={{
                  scaleX: 1.04,
                }}
                className="mt-3 bg-[#2A2A3D] items-center gap-x-2 text-[14px] w-fit px-10 text-left flex justify-start py-1 rounded-lg border border-white/15 cursor-pointer"
              >
                <FaSquareXTwitter /> Twitter
              </motion.div>
            </div>
          </div>
        </div>

        <div className="flex flex-col p-5 bg-[#2A2A45] border border-white/15 min-h-fit h-fit rounded-md">
          <div className="flex items-start gap-x-5 h-fit w-fit">
            <Avatar size={"lg"} className="mt-5">
              <AvatarImage
                src={filipPfp}
                className="grayscale"
                alt="@terzicdsgn"
              />
              <AvatarBadge>
                <MdVerified />
              </AvatarBadge>
            </Avatar>
            <div className="flex flex-col">
              <div className="flex items-center gap-x-2 mt-3">
                <h2 className="font-semibold text-xl">Wisdom Moore</h2>
                <Badge className="text-xs">MacOS tester</Badge>
              </div>
              <p className="text-white/60 font-light text-md">
                Made it possible for ProxyBeast to reach a stable MacOS release
                by providing extensive feedback and QA resolution.
              </p>

              <motion.div
                onClick={() =>
                  openUrl("https://www.linkedin.com/in/wisdom-moore-a2741126b")
                }
                whileHover={{
                  scaleX: 1.04,
                }}
                className="mt-3 bg-[#2A2A3D] items-center gap-x-2 text-[14px] w-fit px-10 text-left flex justify-start py-1 rounded-lg border border-white/15 cursor-pointer"
              >
                <FaLinkedinIn />
              </motion.div>
            </div>
          </div>
        </div>
      </div>

      <motion.div
        onClick={() => openUrl("https://github.com/z3ntl3/ProxyBeast-v2")}
        layout
        transition={{
          repeat: Infinity,
          duration: 1.4,
          ease: "circInOut",
          delayChildren: stagger(0.3),
        }}
        animate={{
          translateY: [0, 10, 0],
        }}
        className="absolute cursor-pointer right-5 bottom-5 bg-primary text-black/90 p-2 rounded-full shadow-md shadow-primary/30"
      >
        <FaGithub />
      </motion.div>
    </div>
  );
}
