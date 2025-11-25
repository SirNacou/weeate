import { ImageWithFallback } from "@/components/image-with-fallback";
import LucideForkKnifeCrossed from "~icons/lucide/fork-knife-crossed";

type Props = {
  src: string;
  alt: string;
  className?: string;
  placeholderClassName?: string;
};

export function FoodImage({
  src,
  alt,
  className,
  placeholderClassName,
}: Props) {
  return (
    <ImageWithFallback
      src={src}
      alt={alt}
      className={className}
      placeholderClassName={placeholderClassName}
    >
      <ImageWithFallback.ErrorContent>
        <LucideForkKnifeCrossed className="size-6 sm:size-8 mb-2" />
        <span className="text-sm sm:text-base">Image not available</span>
      </ImageWithFallback.ErrorContent>
    </ImageWithFallback>
  );
}
