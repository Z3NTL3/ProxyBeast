import { Tooltip, TooltipContent, TooltipTrigger } from "./ui/tooltip";

export default function URI_Tooltip() {
  return (
    <Tooltip>
      <TooltipTrigger className="text-blue-400 underline">URI</TooltipTrigger>
      <TooltipContent className="text-[12px] font-semibold text-gray-800">
        socks5://user:pass@192.168.1.1:9000
      </TooltipContent>
    </Tooltip>
  )
}
