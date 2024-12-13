package server

import (
	songsHttp "music-library/internal/songs/delivery/http"
	songRepository "music-library/internal/songs/repository"
	songsUseCase "music-library/internal/songs/usecase"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (s *Server) MapHandlers(g *gin.Engine) error {

	// Init repositories
	sRepo := songRepository.NewSongsRepository(s.db)

	// Init useCases
	songsUC := songsUseCase.NewSongsUseCase(s.cfg, sRepo, s.logger)

	// Init handlers
	songsHandlers := songsHttp.NewSongsHandlers(s.cfg, songsUC, s.logger)

	// URL: /swagger/index.html
	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	v1 := g.Group("/api")

	songsGroup := v1.Group("/songs")

	songsHttp.MapSongsRoutes(songsGroup, songsHandlers)
	return nil
}
