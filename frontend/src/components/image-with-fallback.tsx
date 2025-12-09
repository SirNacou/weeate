import { Skeleton } from "@/components/ui/skeleton";
import { cn } from "@/lib/utils";
import { Image, type IKImageProps } from "@imagekit/react";
import { useState, type ReactNode } from "react";
import LucideOctagonX from "~icons/lucide/octagon-x";

type ErrorContentProps = {
  children: ReactNode;
};

function ErrorContent({ children }: ErrorContentProps) {
  return <>{children}</>;
}

type Props = Omit<IKImageProps, "onLoad" | "onError"> & {
  placeholderClassName?: string;
  children?: React.ReactNode;
};

function ImageWithFallbackRoot({
  src,
  alt,
  className,
  placeholderClassName,
  children,
  ...imageProps
}: Props) {
  const [isImageError, setIsImageError] = useState(src !== "" ? false : true);
  const [isImageLoading, setIsImageLoading] = useState(true);

  // Extract ErrorContent from children
  let errorContent: ReactNode = null;
  if (children) {
    const childArray = Array.isArray(children) ? children : [children];
    errorContent = childArray.find(
      (child) =>
        typeof child === "object" &&
        child !== null &&
        "type" in child &&
        child.type === ErrorContent
    );
  }

  const defaultErrorContent = (
    <>
      <LucideOctagonX className="size-6 sm:size-8 mb-2" />
      <span className="text-sm sm:text-base">Image not available</span>
    </>
  );

  return (
    <>
      {isImageError ?
        <div
          className={cn(
            "w-full h-full flex flex-col items-center justify-center text-slate-400",
            placeholderClassName
          )}
        >
          {errorContent ?? defaultErrorContent}
        </div>
      : <>
          {isImageLoading && (
            <Skeleton className="w-full h-full absolute inset-0" />
          )}
          <Image
            src={src}
            alt={alt}
            className={cn(
              "w-full h-full object-cover text-center",
              isImageLoading ? "opacity-0" : "opacity-100",
              className
            )}
            loading="lazy"
            onLoad={() => setIsImageLoading(false)}
            onError={() => {
              setIsImageError(true);
              setIsImageLoading(false);
            }}
            {...imageProps}
          />
        </>
      }
    </>
  );
}

export const ImageWithFallback = Object.assign(ImageWithFallbackRoot, {
  ErrorContent,
});
