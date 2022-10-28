package server

import (
	"pavel-fokin/images-storage/internal/images"
	"pavel-fokin/images-storage/internal/server/api"
)

func (s *Server) SetupImagesAPIRoutes(images *images.Images) {
	s.router.Get("/v1/images", api.ImagesGet(images))
	s.router.Post("/v1/images", api.ImagesPost(images))
}
