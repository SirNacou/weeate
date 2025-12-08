import { Image as ImageIcon, Pencil, Upload } from "lucide-react"
import { useCallback, useState } from "react"
import { useDropzone } from "react-dropzone"
import Cropper, { Area } from "react-easy-crop"

import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { Slider } from "@/components/ui/slider"
import { getCroppedImg } from '@/utils/canvasUtil'
import { UserAvatar } from "./user-avatar"

interface AvatarUploadProps {
  initialImage?: string
  onSave?: (croppedImage: string) => void
}

export default function AvatarUpload({ initialImage, onSave }: AvatarUploadProps) {
  const [imageSrc, setImageSrc] = useState<string | null>(null)
  const [croppedImage, setCroppedImage] = useState<string | null>(initialImage || null)

  // Crop State
  const [crop, setCrop] = useState({ x: 0, y: 0 })
  const [zoom, setZoom] = useState(1)
  const [croppedAreaPixels, setCroppedAreaPixels] = useState<Area | null>(null)

  // Dialog State
  const [isDialogOpen, setIsDialogOpen] = useState(false)

  // 1. Handle File Drop/Selection
  const onDrop = useCallback((acceptedFiles: File[]) => {
    if (acceptedFiles && acceptedFiles.length > 0) {
      const file = acceptedFiles[0]
      const reader = new FileReader()
      reader.addEventListener("load", () => {
        setImageSrc(reader.result?.toString() || "")
        setIsDialogOpen(true)
      })
      reader.readAsDataURL(file)
    }
  }, [])

  const { getRootProps, getInputProps, isDragActive } = useDropzone({
    onDrop,
    accept: { "image/*": [] },
    multiple: false,
  })

  // 2. Handle Crop Completion
  const onCropComplete = useCallback((croppedArea: Area, croppedAreaPixels: Area) => {
    setCroppedAreaPixels(croppedAreaPixels)
  }, [])

  // 3. Save the Cropped Image
  const handleSave = async () => {
    try {
      if (imageSrc && croppedAreaPixels) {
        const croppedImageBase64 = await getCroppedImg(imageSrc, croppedAreaPixels)
        if (croppedImageBase64) {
          setCroppedImage(croppedImageBase64)
          setIsDialogOpen(false)
          // Clean up local state
          setZoom(1)
          setCrop({ x: 0, y: 0 })
          // Trigger parent callback
          if (onSave) onSave(croppedImageBase64)
        }
      }
    } catch (e) {
      console.error(e)
    }
  }

  return (
    <div className="flex flex-col items-center gap-4">
      {/* 4. The Dropzone Avatar Trigger */}
      <div
        {...getRootProps()}
        className={`relative group cursor-pointer rounded-full p-1 transition-all duration-200 
          ${isDragActive ? "ring-4 ring-primary/30" : "hover:ring-4 hover:ring-muted"}`}
      >
        <input {...getInputProps()} />
        <UserAvatar src={croppedImage || "/miku monitoring_EbEgOC38U"}
          className="h-32 w-32 border-2 border-border shadow-sm"
          fallback={<ImageIcon className="h-10 w-10 opacity-50" />} />
        {/* <Avatar className="h-32 w-32 border-2 border-border shadow-sm">
          <AvatarImage src={croppedImage || ""} className="object-cover" />
          <AvatarFallback className="bg-muted text-muted-foreground">
            <ImageIcon className="h-10 w-10 opacity-50" />
          </AvatarFallback>
        </Avatar> */}

        {/* Hover Overlay */}
        <div className="absolute inset-0 flex items-center justify-center rounded-full bg-black/40 opacity-0 transition-opacity group-hover:opacity-100">
          <Upload className="h-8 w-8 text-white" />
        </div>

        {/* Helper Badge */}
        <div className="absolute bottom-0 right-0 rounded-full bg-primary p-2 shadow-md border-2 border-background">
          <Pencil className="h-4 w-4 text-primary-foreground" />
        </div>
      </div>

      <p className="text-sm text-muted-foreground">
        Click or drag to upload photo
      </p>

      {/* 5. Cropping Dialog */}
      <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>Crop your new profile picture</DialogTitle>
          </DialogHeader>

          <div className="relative mt-4 h-64 w-full overflow-hidden rounded-md bg-muted">
            {imageSrc && (
              <Cropper
                image={imageSrc}
                crop={crop}
                zoom={zoom}
                aspect={1}
                cropShape="round"
                showGrid={false}
                onCropChange={setCrop}
                onCropComplete={onCropComplete}
                onZoomChange={setZoom}
              />
            )}
          </div>

          <div className="mt-4 flex items-center gap-4">
            <span className="text-sm font-medium">Zoom</span>
            <Slider
              value={[zoom]}
              min={1}
              max={3}
              step={0.1}
              onValueChange={(value) => setZoom(value[0])}
              className="flex-1"
            />
          </div>

          <DialogFooter className="mt-4 sm:justify-between">
            <Button
              variant="secondary"
              onClick={() => setIsDialogOpen(false)}
            >
              Cancel
            </Button>
            <Button onClick={handleSave}>Save changes</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}