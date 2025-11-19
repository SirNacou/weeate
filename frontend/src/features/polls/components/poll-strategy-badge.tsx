import LucideUserRound from "~icons/lucide/user-round";
import LucideTrophy from "~icons/lucide/trophy";
import { Badge } from "@/components/ui/badge";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import { cn } from "@/lib/utils";

type PollStrategyBadgeProps = {
  strategy: "ORDER_CONSENSUS_ITEM" | "ORDER_PERSONAL_CHOICE";
  className?: string;
};

const PollStrategyBadge = ({ strategy, className }: PollStrategyBadgeProps) => {
  const isConsensus = strategy === "ORDER_CONSENSUS_ITEM";

  // CONFIGURATION
  const config =
    isConsensus ?
      {
        label: "Majority Rules",
        // Icon: Shows a Trophy (Winning)
        Icon: LucideTrophy,
        // Visual: Uniformity (Same color/style)
        style:
          "bg-amber-100 text-amber-800 border-amber-200 hover:bg-amber-100",
        tooltipTitle: "Winner Takes All",
        tooltipDesc:
          "Everyone eats the item with the most votes. You might not get what you picked.",
      }
    : {
        label: "Personal Choice",
        // Icon: Shows User (Individual)
        Icon: LucideUserRound,
        // Visual: Variety (Calm color)
        style:
          "bg-slate-100 text-slate-700 border-slate-200 hover:bg-slate-100",
        tooltipTitle: "You Get What You Pick",
        tooltipDesc: "Your vote determines your specific meal.",
      };

  const MainIcon = config.Icon;

  return (
    <TooltipProvider>
      <Tooltip delayDuration={200}>
        <TooltipTrigger asChild>
          <div
            className={cn(
              "inline-flex items-center gap-2 cursor-help",
              className
            )}
          >
            {/* The Badge */}
            <Badge
              variant="outline"
              className={cn("gap-1.5 pl-2 py-1", config.style)}
            >
              <MainIcon className="h-3.5 w-3.5" />
              <span className="font-semibold text-xs uppercase tracking-wide">
                {config.label}
              </span>
            </Badge>
          </div>
        </TooltipTrigger>

        {/* The Explanation */}
        <TooltipContent className="max-w-[250px]">
          <p className="font-bold mb-1">{config.tooltipTitle}</p>
          <p className="text-xs">{config.tooltipDesc}</p>
        </TooltipContent>
      </Tooltip>
    </TooltipProvider>
  );
};

export default PollStrategyBadge;
