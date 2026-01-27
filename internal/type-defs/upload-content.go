package typedefs

type UploadContentType struct {
	Text  string `validate:"required"`
	Image string
	Video string
}
