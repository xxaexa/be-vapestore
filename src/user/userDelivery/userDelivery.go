package userDelivery

import (
	"clean-architecture/model/dto/json"
	"clean-architecture/model/dto/userDto"
	"clean-architecture/pkg/middleware"
	"clean-architecture/pkg/validation"
	"clean-architecture/src/user"
	"clean-architecture/utils"
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
)

type userDelivery struct {
	userUC user.UserUseCase
}

func NewUserDelivery(v1Group *gin.RouterGroup, userUC user.UserUseCase) {
	handler := userDelivery{
		userUC: userUC,
	}

	// Group for operations that require Basic Auth
	basicAuthGroup := v1Group.Group("/users")
	{
		basicAuthGroup.POST("/login", middleware.BasicAuthEncode, handler.loginUser)
		basicAuthGroup.POST("/create", handler.registerUser)
	}

	// Group for operations that require JWT Auth
	jwtAuthGroup := v1Group.Group("/users", middleware.JwtAuth())
	{
		jwtAuthGroup.GET("", handler.getUsers)
		jwtAuthGroup.GET("/:id", handler.getUserByID)
		jwtAuthGroup.PUT("/:id", handler.updateUser)
		jwtAuthGroup.DELETE("/:id", handler.deleteUser)
	}
}

func (c *userDelivery) registerUser(ctx *gin.Context) {
	var userPayload *userDto.CreateUserRequest
	if err := ctx.ShouldBindJSON(&userPayload); err != nil {
		validationError := validation.GetValidationError(err)
		if len(validationError) > 0 {
			json.NewResponseBadRequest(ctx, validationError, "bad request", "01", "01")
			return
		}
	}

	existingUser, err := c.userUC.GetUserByEmail(userPayload.Email)
	if err != nil && err != sql.ErrNoRows {
		json.NewResponseError(ctx, "Failed to query database", "01", "03")
		return
	}

	if existingUser != nil {
		json.NewResponseForbidden(ctx, "email already in use", "01", "02")
		return
	}

	if !c.userUC.IsValidPassword(userPayload.Password) {
		json.NewResponseForbidden(ctx, "password must contain upper and lower case letters", "01", "03")
		return
	}

	hashedPassword, err := c.userUC.HashPassword(userPayload.Password)
	if err != nil {
		json.NewResponseError(ctx, "Internal server error", "01", "04")
	}

	err = c.userUC.CreateUser(&userDto.CreateUserRequest{
		Email:    userPayload.Email,
		FullName: userPayload.FullName,
		Password: hashedPassword,
	})
	if err != nil {
		json.NewResponseError(ctx, err.Error(), "01", "03")
		return
	}

	json.NewResponseSuccess(ctx, nil, "success", "01", "04")
}

func (c *userDelivery) loginUser(ctx *gin.Context) {
	var userPayload *userDto.LoginUserRequest
	if err := ctx.ShouldBindJSON(&userPayload); err != nil {
		validationError := validation.GetValidationError(err)
		if len(validationError) > 0 {
			json.NewResponseBadRequest(ctx, validationError, "bad request", "02", "01")
			return
		}
	}

	user, err := c.userUC.GetUserByEmail(userPayload.Email)
	if err != nil {
		json.NewResponseForbidden(ctx, "invalid email", "02", "02")
		return
	}

	if !c.userUC.ComparePasswords(user.Password, []byte(userPayload.Password)) {
		json.NewResponseForbidden(ctx, "invalid email or password", "02", "03")
		return
	}

	token, err := middleware.GenerateTokenJwt(user.ID, 3)
	if err != nil {
		json.NewResponseForbidden(ctx, "invalid email", "02", "04")
		return
	}

	json.NewResponseSuccess(ctx, token, "success", "02", "05")
}

func (c *userDelivery) getUsers(ctx *gin.Context) {

	pageStr := ctx.Query("page")
	limitStr := ctx.Query("size")
	email := ctx.Query("email")
	fullName := ctx.Query("fullname")

	page, _ := utils.StrToInt(pageStr)
	limit, _ := utils.StrToInt(limitStr)

	users, count, err := c.userUC.GetUsers(page, limit, email, fullName)
	if err != nil {
		json.NewResponseError(ctx, err.Error(), "03", "01")
		return
	}

	json.NewResponseSuccessPage(ctx, users, page, count, "success", "03", "02")
}

func (c *userDelivery) getUserByID(ctx *gin.Context) {

	ID := ctx.Param("id")

	users, err := c.userUC.GetUserByID(ID)
	if err != nil {
		json.NewResponseError(ctx, err.Error(), "04", "01")
		return
	}

	json.NewResponseSuccess(ctx, users, "success", "04", "02")
}

func (c *userDelivery) updateUser(ctx *gin.Context) {
	ID := ctx.Param("id")
	userPayload := &userDto.UpdateUserRequest{}
	userPayload.ID = ID
	if err := ctx.ShouldBindJSON(&userPayload); err != nil {
		validationError := validation.GetValidationError(err)
		if len(validationError) > 0 {
			json.NewResponseBadRequest(ctx, validationError, "bad request", "05", "01")
			return
		}
	}

	hashedPassword, err := c.userUC.HashPassword(userPayload.Password)
	if err != nil {
		json.NewResponseError(ctx, "Internal server error", "05", "02")
	}
	userPayload.Password = hashedPassword

	err = c.userUC.UpdateUser(userPayload)
	if err != nil {
		json.NewResponseError(ctx, err.Error(), "05", "03")
		return
	}

	json.NewResponseSuccess(ctx, nil, "success", "05", "02")
}

func (c *userDelivery) deleteUser(ctx *gin.Context) {

	ID := ctx.Param("id")
	IDLogged, _ := ctx.Get("userID")
	fmt.Println(IDLogged, ID)
	if ID == IDLogged.(string) {
		json.NewResponseForbidden(ctx, "cant delete yourself", "06", "01")
		return
	}

	err := c.userUC.DeleteUser(ID)
	if err != nil {
		json.NewResponseError(ctx, err.Error(), "06", "02")
		return
	}

	json.NewResponseSuccess(ctx, nil, "success", "06", "03")
}
