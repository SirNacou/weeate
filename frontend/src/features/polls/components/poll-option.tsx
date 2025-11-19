import { Label } from "@/components/ui/label";
import { RadioGroupItem } from "@/components/ui/radio-group";
import { Image } from "@imagekit/react";
import TablerCurrencyDong from "~icons/tabler/currency-dong?width=2em&height=2em";

export type Option = {
  id: string;
  foodName: string;
  foodImageUrl: string;
  price: number;
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
  return (
    <div
      key={option.id}
      className={`relative rounded-xl border-2 transition-all overflow-hidden cursor-pointer ${
        isSelected ?
          "border-primary shadow-lg"
        : "border-slate-200 hover:border-primary-300"
      } ${disabled ? "opacity-60" : ""}`}
      onClick={() => !disabled && onSelect?.(option.id)}
    >
      <div className="h-32 sm:h-40 bg-slate-100 overflow-hidden relative">
        <Image
          src={option.foodImageUrl}
          alt={option.foodName}
          className="w-full h-full object-cover text-center"
          loading="lazy"
          decoding="async"
          onError={() => {
            console.log("Image failed to load");
          }}
        />
        <div className="absolute top-2 left-2 bg-white/90 backdrop-blur px-2 py-0.5 sm:px-3 sm:py-1 rounded-full flex items-center gap-0.5 shadow">
          <span className="text-slate-900 font-medium text-xs sm:text-sm">
            {new Intl.NumberFormat("vi-VN").format(option.price)}
          </span>
          <TablerCurrencyDong className="size-3 sm:size-3.5 text-slate-600" />
        </div>
      </div>
      <div className="p-3 sm:p-4 bg-white">
        <div
          className="flex items-center gap-2 sm:gap-3 mb-1 sm:mb-2"
          onClick={(e) => e.stopPropagation()}
        >
          <RadioGroupItem value={option.id} id={`option-${option.id}`} />
          <Label
            htmlFor={`option-${option.id}`}
            className={`flex-1 cursor-pointer text-sm sm:text-base ${false ? "text-slate-900" : "text-slate-700"}`}
          >
            {option.foodName}
          </Label>
        </div>

        {/* Show current voters */}
        {[].length > 0 && (
          <div className="flex items-center gap-2 mt-2 pt-2 border-t border-slate-100">
            <div className="flex -space-x-2">
              {/* {choice.voters.slice(0, 3).map((voter, idx) => (
              <Avatar
                key={idx}
                className="w-6 h-6 border-2 border-white"
              >
                <AvatarImage
                  src={voter.avatar}
                  alt={voter.name}
                />
                <AvatarFallback className="text-xs">
                  {voter.initials}
                </AvatarFallback>
              </Avatar>
            ))}
            {choice.voters.length > 3 && (
              <div className="w-6 h-6 rounded-full bg-slate-200 border-2 border-white flex items-center justify-center">
                <span className="text-xs text-slate-600">
                  +{choice.voters.length - 3}
                </span>
              </div>
            )} */}
            </div>
            {/* <span className="text-slate-500">
              {choice.voters.length} voted
            </span> */}
          </div>
        )}
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
