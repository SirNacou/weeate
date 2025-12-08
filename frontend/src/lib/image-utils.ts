import { env } from "@/env/client"
import { getResponsiveImageAttributes, ResponsiveImageAttributes, Transformation } from "@imagekit/react"

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
  // If the path is already a full URL, return it as is
  if (src.startsWith("http://") || src.startsWith("https://")) {
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