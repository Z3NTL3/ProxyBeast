import { Tooltip, TooltipContent, TooltipTrigger } from "./ui/tooltip";

export default function HelpfulTooltip({title, content} : { title: String, content: String}) {
  return (
    <Tooltip>
      <TooltipTrigger className="text-blue-400 underline">{title}</TooltipTrigger>
      <TooltipContent className="text-[12px] font-semibold text-gray-800">
        {content}
      </TooltipContent>
    </Tooltip>
  )
}
