import { FileWithPreview } from "@/hooks/use-file-upload"
import React, { useCallback, useEffect, useRef, useState } from "react"
import { toast } from "sonner"
import {
	Dialog,
	DialogContent,
	DialogHeader,
	DialogTitle,
} from "./animate-ui/components/radix/dialog"
import AvatarUpload from "./avatar-upload"
import {
	Cropper,
	CropperArea,
	CropperAreaData,
	CropperImage,
	CropperPoint,
} from "./ui/cropper"
import { Slider } from "./ui/slider"

async function createCroppedImage(
	imageSrc: string,
	cropData: CropperAreaData,
	fileName: string,
): Promise<File> {
	const image = new Image()
	image.crossOrigin = "anonymous"

	return new Promise((resolve, reject) => {
		image.onload = () => {
			const canvas = document.createElement("canvas")
			const ctx = canvas.getContext("2d")

			if (!ctx) {
				reject(new Error("Could not get canvas context"))
				return
			}

			canvas.width = cropData.width
			canvas.height = cropData.height

			ctx.drawImage(
				image,
				cropData.x,
				cropData.y,
				cropData.width,
				cropData.height,
				0,
				0,
				cropData.width,
				cropData.height,
			)

			canvas.toBlob((blob) => {
				if (!blob) {
					reject(new Error("Canvas is empty"))
					return
				}

				const croppedFile = new File([blob], `cropped-${fileName}`, {
					type: "image/png",
				})
				resolve(croppedFile)
			}, "image/png")
		}

		image.onerror = () => reject(new Error("Failed to load image"))
		image.src = imageSrc
	})
}

interface Props {
	src?: string
	alt?: string
	fallback?: React.ReactNode
	onImageChange?: (imageUrl: string) => void
	className?: string
	size?: number
}

const EditableAvatar = ({
	src,
	onImageChange,
}: Props) => {
	const [isDialogOpen, setIsDialogOpen] = useState(false)
	const [image, setImage] = useState<FileWithPreview | null>(null)
	const avatarRef = useRef<{ openFileDialog: () => void }>(null)

	const handleFileChange = (file: FileWithPreview | null) => {
		console.log("Files selected:", file)
		onImageChange?.(file?.preview || "")
		if (file) {
			setImage(file)
			setIsDialogOpen(true)
		}
	}

	return (
		<>
			<AvatarUpload defaultAvatar={src} maxSize={5 * 1024 * 1024} onFileChange={handleFileChange} />
			<CropperDialog
				open={isDialogOpen}
				onOpenChange={setIsDialogOpen}
				imageUrl={image?.preview}
				file={image?.file as File}
			/>
		</>
	)
}

interface FileWithCrop {
	original: File
	cropped?: File
}

function CropperDialog({
	open,
	onOpenChange,
	file,
	imageUrl,
}: {
	open: boolean
	onOpenChange: (open: boolean) => void
	file?: File
	imageUrl?: string
}) {
	const [crop, setCrop] = useState<CropperPoint>({ x: 0, y: 0 })
	const [croppedArea, setCroppedArea] = useState<CropperAreaData | null>(null)
	const [zoom, setZoom] = useState(1)
	const [fileWithCrop, setFileWithCrop] = useState<FileWithCrop | null>(null)

	useEffect(() => {
		return () => {
			// Clean up object URLs
			if (imageUrl) {
				URL.revokeObjectURL(imageUrl)
			}
		}
	}, [imageUrl])

	const onCropApply = useCallback(async () => {
		if (!file || !croppedArea || !imageUrl) return

		try {
			const croppedFile = await createCroppedImage(
				imageUrl,
				croppedArea,
				file.name,
			)

			setFileWithCrop({ original: file, cropped: croppedFile })

			onOpenChange(false)
		} catch (error) {
			toast.error(
				error instanceof Error ? error.message : "Failed to crop image",
			)
		}
	}, [file, croppedArea, imageUrl, onOpenChange])

	return (
		<Dialog open={open} onOpenChange={onOpenChange}>
			<DialogContent>
				<DialogHeader>
					<DialogTitle>Crop Image</DialogTitle>
				</DialogHeader>
				<div className="flex flex-col gap-6">
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
						className="min-h-[500px]"
						onCropSizeChange={(cropSize) => {
							console.log("Crop size:", cropSize)
						}}
						onMediaLoaded={(mediaSize) => {
							console.log("Media size:", mediaSize)
						}}
					>
						<CropperImage src={imageUrl || undefined} alt="Cropper Image" />
						<CropperArea />
					</Cropper>

					<Slider
						value={[zoom]}
						onValueChange={(value) => setZoom(value[0])}
						min={1}
						max={3}
						step={0.01}
					/>
				</div>
			</DialogContent>
		</Dialog>
	)
}

export default EditableAvatar
