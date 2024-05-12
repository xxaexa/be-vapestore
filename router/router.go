package router

import (
	"clean-architecture/src/user/userDelivery"
	"clean-architecture/src/user/userRepository"
	"clean-architecture/src/user/userUseCase"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func InitRoute(v1Group *gin.RouterGroup, db *sql.DB) {
	userRepo := userRepository.NewUserRepository(db)
	userUc := userUseCase.NewUserUseCase(userRepo)
	userDelivery.NewUserDelivery(v1Group, userUc)
}
