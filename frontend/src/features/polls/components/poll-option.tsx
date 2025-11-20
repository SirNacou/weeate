import { Label } from "@/components/ui/label";
import { RadioGroupItem } from "@/components/ui/radio-group";
import { Skeleton } from "@/components/ui/skeleton";
import { Image } from "@imagekit/react";
import TablerCurrencyDong from "~icons/tabler/currency-dong?width=2em&height=2em";
import LucideForkKnifeCrossed from "~icons/lucide/fork-knife-crossed";
import { useState } from "react";
import { cn } from "@/lib/utils";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import {
  AvatarGroup,
  AvatarGroupTooltip,
} from "@/components/animate-ui/components/animate/avatar-group";

export type Vote = {
  userId: string;
  userName: string;
  userAvatarUrl: string;
};

export type Option = {
  id: string;
  foodName: string;
  foodImageUrl: string;
  price: number;
  votes: Vote[];
};

type Props = {
  option: Option;
  onSelect?: (id: string) => void;
  isSelected?: boolean;
  disabled?: boolean;
};

function PollOption({
  option,
  onSelect,
  isSelected = false,
  disabled = false,
}: Props) {
  const [isImageError, setIsImageError] = useState(false);
  const [isImageLoading, setIsImageLoading] = useState(true);

  return (
    <div
      key={option.id}
      className={cn(
        "relative rounded-xl border-2 transition-all overflow-hidden cursor-pointer",
        {
          "border-primary shadow-lg": isSelected,
          "border-slate-200 hover:border-zinc-400": !isSelected,
          "opacity-60": disabled,
        }
      )}
      onClick={() => !disabled && onSelect?.(option.id)}
    >
      <div className="h-40 sm:h-48 bg-slate-100 overflow-hidden relative">
        {isImageError ?
          <div className="w-full h-full flex flex-col items-center justify-center text-slate-400">
            <LucideForkKnifeCrossed className="size-6 sm:size-8 mb-2" />
            <span className="text-sm sm:text-base">Image not available</span>
          </div>
        : <>
            {isImageLoading && (
              <Skeleton className="w-full h-full absolute inset-0" />
            )}
            <Image
              src={option.foodImageUrl}
              alt={option.foodName}
              className={`w-full h-full object-cover text-center ${isImageLoading ? "opacity-0" : "opacity-100"}`}
              loading="lazy"
              onLoad={() => setIsImageLoading(false)}
              onError={() => {
                setIsImageError(true);
                setIsImageLoading(false);
              }}
            />
          </>
        }
        <div className="absolute top-2 left-2 bg-white/90 backdrop-blur px-2 py-0.5 sm:px-3 sm:py-1 rounded-full flex items-center gap-0.5 shadow">
          <span className="text-slate-900 font-medium text-xs sm:text-sm">
            {new Intl.NumberFormat("vi-VN").format(option.price)}
          </span>
          <TablerCurrencyDong className="size-3 sm:size-3.5 text-slate-600" />
        </div>
      </div>
      <div className="p-3 sm:p-4 bg-white flex justify-between items-start gap-2 relative">
        <div
          className="flex items-center gap-2 sm:gap-3 mb-1 sm:mb-2"
          onClick={(e) => e.stopPropagation()}
        >
          <RadioGroupItem
            className="size-6 [&_svg]:size-4"
            value={option.id}
            id={option.id}
          />
          <Label
            htmlFor={`option-${option.id}`}
            className={`flex-1 cursor-pointer text-sm sm:text-base ${false ? "text-slate-900" : "text-slate-700"}`}
          >
            {option.foodName}
          </Label>
        </div>

        {/* Show current voters */}
        {option.votes.length > 0 && (
          <div className="flex flex-col gap-1 self-end items-end">
            <AvatarGroup>
              {option.votes.map((vote, index) => (
                <Avatar key={index}>
                  <AvatarImage src={vote.userAvatarUrl} />
                  <AvatarFallback>
                    {vote.userName.substring(0, 2).toUpperCase()}
                  </AvatarFallback>
                  <AvatarGroupTooltip>{vote.userName}</AvatarGroupTooltip>
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

export default PollOption;
