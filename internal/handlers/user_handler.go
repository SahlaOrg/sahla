package handlers

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gofrs/uuid"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	
	domain "github.com/mohamed2394/sahla/internal/domains"
    dto "github.com/mohamed2394/sahla/internal/dtos"
	repository "github.com/mohamed2394/sahla/internal/repositories"


	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type UserHandler struct {
	userRepository repository.UserRepository
	minioClient    *minio.Client
}

func NewUserHandler(userRepository repository.UserRepository) *UserHandler {
	// Initialize MinIO client
	endpoint := "localhost:9000"
	accessKeyID := "your-access-key"
	secretAccessKey := "your-secret-key"
	useSSL := false

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		panic(err)
	}

	return &UserHandler{
		userRepository: userRepository,
		minioClient:    minioClient,
	}
}
func (h *UserHandler) UploadIDImage(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.FromString(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user ID"})
	}

	// Get the file from the request
	file, err := c.FormFile("id_image")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "No file uploaded"})
	}

	// Open the file
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error opening file"})
	}
	defer src.Close()

	// Upload file to MinIO
	bucketName := "user-id-images"
	objectName := fmt.Sprintf("%s%s", id, filepath.Ext(file.Filename))
	contentType := file.Header.Get("Content-Type")

	_, err = h.minioClient.PutObject(context.Background(), bucketName, objectName, src, file.Size, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error uploading file to MinIO"})
	}

	// Generate the URL for the uploaded image
	imageURL := fmt.Sprintf("http://localhost:9000/%s/%s", bucketName, objectName)

	// Update the user's ID image URL in the database
	err = h.userRepository.UpdateIDImage(id, imageURL)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error updating user ID image URL"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "ID image uploaded successfully", "url": imageURL})
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
