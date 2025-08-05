package utils

import (
  "io"
  "errors"
  "strings"
  "bytes"
  "net/http"
  "mime/multipart"
  "github.com/disintegration/imaging"
)

type RouteConversions struct {
  Main      []byte `json:"main,omitempty"`
  Thumbnail []byte `json:"thumbnail,omitempty"`
}

func GenerateRouteImages(input multipart.File) (*RouteConversions, error) {
	var buf bytes.Buffer
	_, err := io.Copy(&buf, input)
	if err != nil {
		return nil, err
	}

	contentType := http.DetectContentType(buf.Bytes())
	if !strings.HasPrefix(contentType, "image/jpeg") && !strings.HasPrefix(contentType, "image/png") {
		return nil, errors.New("unsupported file type: only JPEG and PNG are allowed")
	}

	imgReader := bytes.NewReader(buf.Bytes())
	img, err := imaging.Decode(imgReader)
	if err != nil {
		return nil, errors.New("failed to decode image")
	}

	thumb := imaging.Fill(img, 150, 150,  imaging.Center, imaging.Lanczos)
	main := imaging.Fill(img, 500, 600,  imaging.Center, imaging.Lanczos)

	mainBuf := new(bytes.Buffer)
	if err := imaging.Encode(mainBuf, main, getImagingFormat(contentType)); err != nil {
		return nil, err
	}

	thumbBuf := new(bytes.Buffer)
	if err := imaging.Encode(thumbBuf, thumb, getImagingFormat(contentType)); err != nil {
		return nil, err
	}

	return &RouteConversions{
		Main: mainBuf.Bytes(),
		Thumbnail: thumbBuf.Bytes(),
	}, nil
}

func getImagingFormat(contentType string) imaging.Format {
	switch contentType {
	case "image/jpeg":
		return imaging.JPEG
	case "image/png":
		return imaging.PNG
	default:
		return imaging.PNG
	}
}
