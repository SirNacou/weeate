import { useCallback, useEffect, useRef, useState } from "react";
import { toast } from "sonner";
import LucideCamera from "~icons/lucide/camera";
import LucideUserRound from "~icons/lucide/user-round";
import { Dialog, DialogContent } from "./animate-ui/components/radix/dialog";
import { Avatar, AvatarFallback, AvatarImage } from "./ui/avatar";
import {
  Cropper,
  CropperArea,
  CropperAreaData,
  CropperImage,
  CropperPoint,
} from "./ui/cropper";

async function createCroppedImage(
  imageSrc: string,
  cropData: CropperAreaData,
  fileName: string
): Promise<File> {
  const image = new Image();
  image.crossOrigin = "anonymous";

  return new Promise((resolve, reject) => {
    image.onload = () => {
      const canvas = document.createElement("canvas");
      const ctx = canvas.getContext("2d");

      if (!ctx) {
        reject(new Error("Could not get canvas context"));
        return;
      }

      canvas.width = cropData.width;
      canvas.height = cropData.height;

      ctx.drawImage(
        image,
        cropData.x,
        cropData.y,
        cropData.width,
        cropData.height,
        0,
        0,
        cropData.width,
        cropData.height
      );

      canvas.toBlob((blob) => {
        if (!blob) {
          reject(new Error("Canvas is empty"));
          return;
        }

        const croppedFile = new File([blob], `cropped-${fileName}`, {
          type: "image/png",
        });
        resolve(croppedFile);
      }, "image/png");
    };

    image.onerror = () => reject(new Error("Failed to load image"));
    image.src = imageSrc;
  });
}

interface Props {
  src?: string;
  alt?: string;
  fallback?: React.ReactNode;
  onImageChange?: (imageUrl: string) => void;
  className?: string;
  size?: number;
}

const EditableAvatar = ({
  src,
  alt,
  fallback,
  onImageChange,
  className,
  size = 128,
}: Props) => {
  const inputRef = useRef<HTMLInputElement>(null);
  const [isDialogOpen, setIsDialogOpen] = useState(false);

  const handleClick = () => {
    inputRef.current?.click();
  };

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      const imageUrl = URL.createObjectURL(file);
      onImageChange?.(imageUrl);

      setIsDialogOpen(true);
    }
  };

  // Calculate responsive sizes based on avatar size
  const fallbackIconSize = Math.max(24, size * 0.375); // 37.5% of avatar size, min 24px
  const fallbackTextSize = Math.max(16, size * 0.25); // 25% of avatar size, min 16px
  const cameraIconSize = Math.max(16, size * 0.1875); // 18.75% of avatar size, min 16px
  const editTextSize = Math.max(10, size * 0.09375); // 9.375% of avatar size, min 10px
  const overlayGap = Math.max(2, size * 0.03125); // 3.125% of avatar size, min 2px

  return (
    <>
      <div
        className={`relative group cursor-pointer ${className}`}
        style={{ width: size, height: size }}
        onClick={handleClick}
      >
        <Avatar className="h-full w-full border-2 border-slate-200 dark:border-slate-800 ring-offset-background transition-all group-hover:ring-2 group-hover:ring-slate-400 group-hover:ring-offset-2">
          <AvatarImage src={src} alt={alt} />
          <AvatarFallback
            className="font-semibold"
            style={{ fontSize: fallbackTextSize }}
          >
            {fallback || (
              <LucideUserRound
                style={{ width: fallbackIconSize, height: fallbackIconSize }}
              />
            )}
          </AvatarFallback>
        </Avatar>

        {/* Hover Overlay */}
        <div className="absolute inset-0 flex items-center justify-center rounded-full bg-black/60 opacity-0 transition-opacity duration-300 group-hover:opacity-100">
          <div
            className="flex flex-col items-center text-white"
            style={{ gap: overlayGap }}
          >
            <LucideCamera
              style={{ width: cameraIconSize, height: cameraIconSize }}
            />
            <span className="font-medium" style={{ fontSize: editTextSize }}>
              Edit
            </span>
          </div>
        </div>

        {/* Hidden File Input */}
        <input
          type="file"
          ref={inputRef}
          onChange={handleFileChange}
          className="hidden"
          accept="image/*"
        />
      </div>
      <CropperDialog open={isDialogOpen} onOpenChange={setIsDialogOpen} />
    </>
  );
};

interface FileWithCrop {
  original: File;
  cropped?: File;
}

function CropperDialog({
  open,
  onOpenChange,
  file,
  imageUrl,
}: {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  file?: File;
  imageUrl?: string;
}) {
  const [crop, setCrop] = useState<CropperPoint>({ x: 0, y: 0 });
  const [croppedArea, setCroppedArea] = useState<CropperAreaData | null>(null);
  const [zoom, setZoom] = useState(1);
  const [fileWithCrop, setFileWithCrop] = useState<FileWithCrop | null>(null);

  useEffect(() => {
    return () => {
      // Clean up object URLs
      if (imageUrl) {
        URL.revokeObjectURL(imageUrl);
      }
    };
  }, [imageUrl]);

  const onCropApply = useCallback(async () => {
    if (!file || !croppedArea || !imageUrl) return;

    try {
      const croppedFile = await createCroppedImage(
        imageUrl,
        croppedArea,
        file.name
      );

      setFileWithCrop({ original: file, cropped: croppedFile });

      onOpenChange(false);
    } catch (error) {
      toast.error(
        error instanceof Error ? error.message : "Failed to crop image"
      );
    }
  }, [file, croppedArea, imageUrl, fileWithCrop, onOpenChange]);

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <Cropper
          aspectRatio={1}
          crop={crop}
          zoom={zoom}
          shape={"circle"}
          objectFit={"contain"}
          onCropChange={setCrop}
          onZoomChange={setZoom}
          onCropAreaChange={setCroppedArea}
          onCropComplete={setCroppedArea}
          className="min-h-72"
        >
          <CropperImage src="" alt="Cropper Image" />
          <CropperArea />
        </Cropper>
      </DialogContent>
    </Dialog>
  );
}

export default EditableAvatar;
