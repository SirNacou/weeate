import {
  Avatar,
  AvatarFallback,
  AvatarImage,
} from "@/components/ui/avatar" // Your Shadcn imports
import { getImageUrl as getImageAttributes } from "@/lib/image-utils"

// 1. Re-use the helper logic
const isAbsoluteUrl = (urlString: string) => {
  try {
    new URL(urlString)
    return true
  } catch (e) {
    return false
  }
}

interface Props {
  src: string
  fallback?: React.ReactNode
  className?: string
}

export function UserAvatar({ src, fallback = "CN", className }: Props) {
  const isOidc = isAbsoluteUrl(src)
  const imageAttributes = getImageAttributes({
    src: src,
    width: 100,
    sizes: "(max-width: 600px) 48px, 100px",
  })

  return (
    <Avatar className={className}>
      {isOidc ? (
        // CASE A: Standard OIDC URL
        // Use default Shadcn behavior
        <AvatarImage src={src} alt="User Avatar" />
      ) : (
        // CASE B: ImageKit Relative Path
        // Use 'asChild' to inject IKImage while keeping Shadcn styling
        <AvatarImage
          {...imageAttributes}
          alt="User Avatar"
          loading="lazy"
        >
        </AvatarImage>
      )}

      {/* Fallback shows automatically if either image fails to load */}
      <AvatarFallback>{fallback}</AvatarFallback>
    </Avatar>
  )
}