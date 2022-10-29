package server

import (
	"pavel-fokin/images-storage/internal/server/api"
)

func (s *Server) SetupImagesAPIRoutes(imagesstorage api.ImagesStorage) {
	s.router.Get("/v1/images", api.ImagesGet(imagesstorage))
	s.router.Post("/v1/images", api.ImagesPost(imagesstorage))
	s.router.Get("/v1/images/{id}", api.ImagesGetByID(imagesstorage))
}
