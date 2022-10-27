package server

import (
	"pavel-fokin/images-storage/internal/server/api"
)

func (s *Server) SetupImagesAPIRoutes() {
	s.router.Get("/v1/images", api.ImagesGet())
	s.router.Post("/v1/images", api.ImagesPost())
}
