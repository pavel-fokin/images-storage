package api

import "github.com/pavel-fokin/images-storage/internal/imagesstorage"

type ResponseImage struct {
	UUID       string `json:"uuid"`
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

type ImagesPostResp struct {
	Data struct {
		Image ResponseImage `json:"image"`
	} `json:"data"`
}

func asResponseImage(image imagesstorage.Image) ResponseImage {
	return ResponseImage{
		UUID:       image.UUID,
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

func asImagesPostResponse(image imagesstorage.Image) ImagesPostResp {
	resp := ImagesPostResp{}

	resp.Data.Image = asResponseImage(image)

	return resp
}
