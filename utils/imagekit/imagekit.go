package imagekit

import (
	"os"

	"github.com/imagekit-developer/imagekit-go/v2"
	"github.com/imagekit-developer/imagekit-go/v2/option"
)

var ImageKitClient imagekit.Client

func InitImageKit() {
	prvtKey := os.Getenv("IMAGEKIT_PRIVATE_KEY")
	client := imagekit.NewClient(
		option.WithPrivateKey(prvtKey),
	)
	ImageKitClient = client
}
