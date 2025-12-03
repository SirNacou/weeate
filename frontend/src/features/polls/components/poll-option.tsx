import { PollOption } from "@/client";
import {
  AvatarGroup,
  AvatarGroupTooltip,
} from "@/components/animate-ui/components/animate/avatar-group";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Label } from "@/components/ui/label";
import { FoodImage } from "@/features/foods/components/food-image";
import { cn } from "@/lib/utils";
import { Circle } from "lucide-react";
import TablerCurrencyDong from "~icons/tabler/currency-dong?width=2em&height=2em";

type Props = {
  option: PollOption;
  onSelect?: (id: string) => void;
  isSelected?: boolean;
  disabled?: boolean;
};

function PollOptionRadio({
  option,
  onSelect,
  isSelected = false,
  disabled = false,
}: Props) {
  return (
    <div
      key={option.id}
      className={cn(
        "relative rounded-xl border-2 transition-all overflow-hidden cursor-pointer border-slate-200 ",
        {
          "border-primary shadow-lg": isSelected,
          "hover:border-zinc-400": !isSelected && !disabled,
          "opacity-50 cursor-not-allowed": disabled,
        }
      )}
      onClick={() => !disabled && onSelect?.(option.id)}
    >
      <div className="h-40 sm:h-48 bg-slate-100 overflow-hidden relative">
        <FoodImage src={""} alt={option.food.name} />
        <div className="absolute top-2 left-2 bg-white/90 backdrop-blur px-2 py-0.5 sm:px-3 sm:py-1 rounded-full flex items-center gap-0.5 shadow">
          <span className="text-slate-900 font-medium text-xs sm:text-sm">
            {new Intl.NumberFormat("vi-VN").format(option.price_at_creation)}
          </span>
          <TablerCurrencyDong className="size-4 sm:size-4.5 text-green-500" />
        </div>
      </div>
      <div className="p-3 sm:p-4 bg-white flex justify-between items-start gap-2 relative">
        <div className="flex items-center gap-2 sm:gap-3 mb-1 sm:mb-2">
          <div
            className={cn(
              "aspect-square h-6 w-6 rounded-full border border-primary text-primary ring-offset-background flex items-center justify-center",
              disabled ? "cursor-not-allowed opacity-50" : "cursor-pointer"
            )}
          >
            {isSelected && (
              <Circle className="h-3.5 w-3.5 fill-current text-current" />
            )}
          </div>
          <Label
            htmlFor={`option-${option.id}`}
            className={cn(
              "flex-1 text-base sm:text-lg",
              disabled ?
                "cursor-not-allowed text-slate-400"
              : "cursor-pointer text-slate-900"
            )}
          >
            {option.food.name}
          </Label>
        </div>
        {/* Show current voters */}
        {option.votes.length > 0 && (
          <div className="flex flex-col gap-1 self-end items-end">
            <AvatarGroup>
              {option.votes.map((vote, index) => (
                <Avatar key={index}>
                  <AvatarImage
                    src={vote.voter.avatar_url}
                    alt={vote.voter.avatar_url}
                  />
                  <AvatarFallback>
                    {vote.voter.display_name.substring(0, 2).toUpperCase()}
                  </AvatarFallback>
                  <AvatarGroupTooltip>
                    {vote.voter.display_name}
                  </AvatarGroupTooltip>
                </Avatar>
              ))}
            </AvatarGroup>

            <span className="text-slate-500">{option.votes.length} voted</span>
          </div>
        )}

        {/* <div className="flex items-center gap-2 mt-2 pt-2 border-t border-slate-100">
            <div className="flex -space-x-2">
              {option.votes.slice(0, 3).map((vote, idx) => (
                <Avatar
                  key={idx}
                  className="size-6 sm:size-10 border-2 border-white"
                >
                  <AvatarImage src={vote.userAvatarUrl} alt={vote.userName} />
                  <AvatarFallback className="text-xs">
                    {vote.userName.substring(0, 2).toUpperCase()}
                  </AvatarFallback>
                </Avatar>
              ))}
              {option.votes.length > 3 && (
                <div className="w-6 h-6 rounded-full bg-slate-200 border-2 border-white flex items-center justify-center">
                  <span className="text-xs text-slate-600">
                    +{option.votes.length - 3}
                  </span>
                </div>
              )}
            </div>
          </div> */}
      </div>
      {isSelected && (
        <div
          className={`absolute top-2 right-2 px-3 py-1 rounded-full shadow-lg bg-primary text-white`}
        >
          ✓ Selected
        </div>
      )}
    </div>
  );
}

export default PollOptionRadio;
