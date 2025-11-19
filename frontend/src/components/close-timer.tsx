import React, { useState, useEffect } from "react";
import { format, differenceInSeconds } from "date-fns";
import { Clock, AlertCircle, Hourglass, CheckCircle2 } from "lucide-react";

// 1. Import the 'cn' utility (standard in shadcn projects)
import { cn } from "@/lib/utils";

import { Badge } from "@/components/ui/badge";
import { Progress } from "@/components/ui/progress";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";

type CloseTimerProps = {
  closesAt: Date | string;
  className?: string; // 2. Accept className as a prop
};

// 2. Accept className as a prop
const CloseTimer = ({ closesAt, className }: CloseTimerProps) => {
  const [now, setNow] = useState(new Date());

  useEffect(() => {
    const timerId = setInterval(() => setNow(new Date()), 1000);
    return () => clearInterval(timerId);
  }, []);

  const closeDate = new Date(closesAt);
  const totalSecondsLeft = differenceInSeconds(closeDate, now);
  const absoluteTime = format(closeDate, "h:mm a");

  // --- STATE 1: CLOSED ---
  if (totalSecondsLeft <= 0) {
    return (
      <Badge
        variant="secondary"
        // 3. Apply cn() to merge defaults with your custom prop
        className={cn("gap-1.5 pl-1.5 cursor-default", className)}
      >
        <CheckCircle2 className="h-3.5 w-3.5 text-muted-foreground" />
        <span className="text-muted-foreground">Poll Closed</span>
      </Badge>
    );
  }

  // --- STATE 2: CRITICAL (< 5 mins) ---
  if (totalSecondsLeft < 300) {
    const minutes = Math.floor(totalSecondsLeft / 60);
    const seconds = totalSecondsLeft % 60;
    const formattedSeconds = seconds < 10 ? `0${seconds}` : seconds;
    const progressValue = (totalSecondsLeft / 300) * 100;

    return (
      // Applied to the wrapper div so layout classes (like margin/width) work
      <div className={cn("flex flex-col gap-1.5 w-fit", className)}>
        <Badge
          variant="destructive"
          className="gap-1.5 pl-1.5 animate-pulse justify-center"
        >
          <AlertCircle className="h-3.5 w-3.5" />
          <span>
            Closing in {minutes}:{formattedSeconds}
          </span>
        </Badge>
        <Progress value={progressValue} className="h-1.5 w-full bg-red-100" />
      </div>
    );
  }

  // --- STATE 3: URGENT (< 1 hour) ---
  if (totalSecondsLeft < 3600) {
    const minutesLeft = Math.floor(totalSecondsLeft / 60);
    return (
      <TooltipProvider>
        <Tooltip>
          <TooltipTrigger asChild>
            <Badge
              className={cn(
                "gap-1.5 pl-1.5 bg-orange-500 hover:bg-orange-600 text-white border-orange-600",
                className // Merge here
              )}
            >
              <Hourglass className="h-3.5 w-3.5" />
              <span>{minutesLeft} mins left</span>
            </Badge>
          </TooltipTrigger>
          <TooltipContent>
            <p>Closes at {absoluteTime}</p>
          </TooltipContent>
        </Tooltip>
      </TooltipProvider>
    );
  }

  // --- STATE 4: NORMAL (> 1 hour) ---
  return (
    <TooltipProvider>
      <Tooltip>
        <TooltipTrigger asChild>
          <Badge
            variant="outline"
            className={cn("gap-1.5 pl-1.5", className)} // Merge here
          >
            <Clock className="h-3.5 w-3.5" />
            <span>Closes at {absoluteTime}</span>
          </Badge>
        </TooltipTrigger>
        <TooltipContent>
          <p>Ends on {format(closeDate, "MMM d, yyyy")}</p>
        </TooltipContent>
      </Tooltip>
    </TooltipProvider>
  );
};

export default CloseTimer;
