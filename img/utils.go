package img

import "strings"

const (
	IMAGE_TYPE_JPEG = "JPEG"
	IMAGE_TYPE_PNG  = "PNG"
)

func GetImageType(name string) string {
	arr := strings.Split(name, ".")
	suffix := arr[len(arr)-1]
	switch suffix {
	case "jpg", "JPG", "jpeg", "JPEG":
		return IMAGE_TYPE_JPEG
	case "png", "PNG":
		return IMAGE_TYPE_PNG
	default:
		return suffix
	}
}
