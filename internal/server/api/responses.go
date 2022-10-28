package api

import (
	"pavel-fokin/images-storage/internal/imagesstorage"
)

type ResponseImage struct {
	Name       string `json:"name"`
	FileType   string `json:"filetype"`
	Width      int    `json:"width"`
	Height     int    `json:"height"`
	FileSize   int    `json:"filesize"`
	UploadedAt string `json:"uploadedAt"`
}

type ImagesGetResp struct {
	Data struct {
		Images []ResponseImage `json:"images"`
	} `json:"data"`
}

func asResponseImage(image imagesstorage.Image) ResponseImage {
	return ResponseImage{
		Name:       image.Name,
		FileType:   image.ContentType,
		Width:      image.Width,
		Height:     image.Height,
		FileSize:   image.Size,
		UploadedAt: image.UploadedAt,
	}
}

func asImagesGetResponse(images []imagesstorage.Image) ImagesGetResp {
	resp := ImagesGetResp{}

	for _, image := range images {
		resp.Data.Images = append(resp.Data.Images, asResponseImage(image))
	}

	return resp
}
