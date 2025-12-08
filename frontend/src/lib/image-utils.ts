import { env } from "@/env/client"
import { getResponsiveImageAttributes, ResponsiveImageAttributes, Transformation } from "@imagekit/react"

/**
 * Check if src is a standard image URL that doesn't need ImageKit processing
 * Returns true for absolute URLs (http/https), blob URLs, and data URLs
 */
export const isStandardImageSrc = (src: string): boolean => {
  return src.startsWith("http://") || 
         src.startsWith("https://") || 
         src.startsWith("blob:") || 
         src.startsWith("data:")
}

export const getImageUrl = ({src, 
  width = 100,
  sizes = "100vw",
  transformation = {
  width: 100, height: 100, crop: "maintain_ratio", quality: 100,aspectRatio: "1:1"
}}: {
  src: string | null,
  width?: number,
  sizes?: string,
  transformation?: Transformation
}): ResponsiveImageAttributes | null => {
  if (!src) return null;
  // If the path is already a standard image source (URL, blob, or data), return it as is
  if (isStandardImageSrc(src)) {
    return {
      src: src,
    }
  }
  // Otherwise, prepend the ImageKit URL endpoint
  const imageKitUrl = getResponsiveImageAttributes({
    src: src,
    urlEndpoint: env.VITE_IMAGEKIT_URL,
    transformation: [
      transformation
    ],
    width,
    sizes
  });
  

  return imageKitUrl;
}