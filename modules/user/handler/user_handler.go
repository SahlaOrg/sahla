package handler

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gofrs/uuid"

	"github.com/labstack/echo/v4"
	"github.com/mohamed2394/sahla/modules/user/domain"
	"github.com/mohamed2394/sahla/modules/user/dto"
	"github.com/mohamed2394/sahla/modules/user/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	userRepository repository.UserRepository
}

func NewUserHandler(userRepository repository.UserRepository) *UserHandler {
	return &UserHandler{userRepository: userRepository}
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	var req dto.CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to hash password"})
	}
	uid, errU := uuid.NewV7()
	if errU != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": errU.Error()})
	}
	user := &domain.User{
		UniversalId:   uid,
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		Email:         req.Email,
		PhoneNumber:   req.PhoneNumber,
		Address:       req.Address,
		PasswordHash:  string(hashedPassword),
		LoyaltyPoints: req.LoyaltyPoints,
	}

	if err := h.userRepository.Create(user); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Convert the user struct to a map to modify the ID
	userResponse := map[string]interface{}{
		"ID":            user.ID,
		"UNIVERSAL_ID":  user.UniversalId,
		"CreatedAt":     user.CreatedAt,
		"UpdatedAt":     user.UpdatedAt,
		"DeletedAt":     user.DeletedAt,
		"FirstName":     user.FirstName,
		"LastName":      user.LastName,
		"Email":         user.Email,
		"PhoneNumber":   user.PhoneNumber,
		"Address":       user.Address,
		"LoyaltyPoints": user.LoyaltyPoints,
	}

	return c.JSON(http.StatusCreated, userResponse)
}

func stringToUUID(s string) (uuid.UUID, error) {
	var uuid uuid.UUID
	bytes, err := hex.DecodeString(s)
	if err != nil || len(bytes) != 16 {
		return uuid, fmt.Errorf("invalid UUID string")
	}
	copy(uuid[:], bytes)
	return uuid, nil
}

func (h *UserHandler) GetUserByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.FromString(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user ID"})
	}

	user, err := h.userRepository.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(c echo.Context) error {
	// Parse the UUID from the request parameter
	idStr := c.Param("id")
	id, err := uuid.FromString(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user ID"})
	}

	var req dto.UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	user, err := h.userRepository.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	// Update only the fields that are provided
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.PhoneNumber != "" {
		user.PhoneNumber = req.PhoneNumber
	}
	if req.Address != "" {
		user.Address = req.Address
	}
	if req.LoyaltyPoints != 0 {
		user.LoyaltyPoints = req.LoyaltyPoints
	}

	// If a new password is provided, hash it
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to hash password"})
		}
		user.PasswordHash = string(hashedPassword)
	}

	if err := h.userRepository.Update(user); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) DeleteUser(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.FromString(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user ID"})
	}

	if err := h.userRepository.Delete(uuid.UUID(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *UserHandler) ListUsers(c echo.Context) error {
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	users, err := h.userRepository.List(offset, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, users)
}
