package handler

import (
	"go-clean-arch/internal/user/domain/entity"
	"go-clean-arch/internal/user/usecase"
	"go-clean-arch/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UserHandler struct {
	userUseCase *usecase.UserUseCase
	logger      *zap.Logger
}

func NewUserHandler(userUseCase *usecase.UserUseCase, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
		logger:      logger,
	}
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	var req entity.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("Failed to parse request body", zap.Error(err))
		return response.Error(c, fiber.StatusBadRequest, "Invalid request", err.Error())
	}

	user, err := h.userUseCase.CreateUser(c.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to create user",
			zap.Error(err),
			zap.String("email", req.Email),
			zap.String("name", req.Name),
		)
		return response.Error(c, fiber.StatusInternalServerError, "Failed to create user", err.Error())
	}

	h.logger.Info("User created successfully",
		zap.String("user_id", user.ID.String()),
		zap.String("email", user.Email),
	)

	return response.Success(c, fiber.StatusCreated, "User created successfully", user)
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req entity.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("Failed to parse login request", zap.Error(err))
		return response.Error(c, fiber.StatusBadRequest, "Invalid request", err.Error())
	}

	token, err := h.userUseCase.Login(c.Context(), &req)
	if err != nil {
		h.logger.Warn("Login failed",
			zap.Error(err),
			zap.String("email", req.Email),
		)
		return response.Error(c, fiber.StatusUnauthorized, "Login failed", err.Error())
	}

	h.logger.Info("User logged in successfully", zap.String("email", req.Email))

	return response.Success(c, fiber.StatusOK, "Login successful", fiber.Map{"token": token})
}

func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return response.Error(c, fiber.StatusUnauthorized, "Unauthorized", "User ID not found")
	}

	id, err := uuid.Parse(userID.(string))
	if err != nil {
		h.logger.Error("Invalid user ID format", zap.Error(err))
		return response.Error(c, fiber.StatusBadRequest, "Invalid user ID", err.Error())
	}

	user, err := h.userUseCase.GetUserByID(c.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get user profile",
			zap.Error(err),
			zap.String("user_id", id.String()),
		)
		return response.Error(c, fiber.StatusNotFound, "User not found", err.Error())
	}

	return response.Success(c, fiber.StatusOK, "User profile retrieved", user)
}
