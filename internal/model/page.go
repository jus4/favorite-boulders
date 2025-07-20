package model

type Metadata struct {
	Title       string
	Description *string
	OgImage     *Image
}

type Image struct {
	Url     string
	AltText *string
}
