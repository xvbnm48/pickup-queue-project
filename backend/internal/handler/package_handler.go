package handler

import (
	"net/http"
	"pickup-queue/internal/domain"
	"pickup-queue/internal/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PackageHandler struct {
	packageUsecase *usecase.PackageUsecase
}

func NewPackageHandler(packageUsecase *usecase.PackageUsecase) *PackageHandler {
	return &PackageHandler{
		packageUsecase: packageUsecase,
	}
}

// CreatePackage creates a new package
// @Summary Create a new package
// @Description Create a new package in the pickup queue
// @Tags packages
// @Accept json
// @Produce json
// @Param package body domain.CreatePackageRequest true "Package details"
// @Success 201 {object} domain.Package
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /packages [post]
func (h *PackageHandler) CreatePackage(c *gin.Context) {
	var req domain.CreatePackageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	pkg, err := h.packageUsecase.CreatePackage(&req)
	if err != nil {
		if err == usecase.ErrDuplicateOrderRef {
			c.JSON(http.StatusConflict, ErrorResponse{Error: "Order reference already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, SuccessResponse{Data: pkg})
}

// GetPackage gets a package by ID
// @Summary Get a package by ID
// @Description Get package details by package ID
// @Tags packages
// @Produce json
// @Param id path string true "Package ID"
// @Success 200 {object} domain.Package
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /packages/{id} [get]
func (h *PackageHandler) GetPackage(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid package ID"})
		return
	}

	pkg, err := h.packageUsecase.GetPackage(id)
	if err != nil {
		if err == usecase.ErrPackageNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "Package not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Data: pkg})
}

// GetPackageByOrderRef gets a package by order reference
// @Summary Get a package by order reference
// @Description Get package details by order reference
// @Tags packages
// @Produce json
// @Param orderRef path string true "Order Reference"
// @Success 200 {object} domain.Package
// @Failure 404 {object} ErrorResponse
// @Router /packages/order/{orderRef} [get]
func (h *PackageHandler) GetPackageByOrderRef(c *gin.Context) {
	orderRef := c.Param("orderRef")

	pkg, err := h.packageUsecase.GetPackageByOrderRef(orderRef)
	if err != nil {
		if err == usecase.ErrPackageNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "Package not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Data: pkg})
}

// ListPackages lists packages with pagination and filtering
// @Summary List packages
// @Description Get a list of packages with pagination and optional status filtering
// @Tags packages
// @Produce json
// @Param limit query int false "Limit" default(50)
// @Param offset query int false "Offset" default(0)
// @Param status query string false "Filter by status"
// @Success 200 {object} PackageListResponse
// @Failure 400 {object} ErrorResponse
// @Router /packages [get]
func (h *PackageHandler) ListPackages(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")
	statusStr := c.Query("status")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	var status *domain.PackageStatus
	if statusStr != "" {
		s := domain.PackageStatus(statusStr)
		if s == domain.StatusWaiting || s == domain.StatusPicked ||
			s == domain.StatusHandedOver || s == domain.StatusExpired {
			status = &s
		}
	}

	packages, err := h.packageUsecase.ListPackages(limit, offset, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	response := PackageListResponse{
		Data:   packages,
		Limit:  limit,
		Offset: offset,
		Count:  len(packages),
	}

	c.JSON(http.StatusOK, response)
}

// UpdatePackageStatus updates package status
// @Summary Update package status
// @Description Update the status of a package
// @Tags packages
// @Accept json
// @Produce json
// @Param id path string true "Package ID"
// @Param status body domain.UpdatePackageStatusRequest true "New status"
// @Success 200 {object} domain.Package
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /packages/{id}/status [patch]
func (h *PackageHandler) UpdatePackageStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid package ID"})
		return
	}

	var req domain.UpdatePackageStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	pkg, err := h.packageUsecase.UpdatePackageStatus(id, req.Status)
	if err != nil {
		if err == usecase.ErrPackageNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "Package not found"})
			return
		}
		if err == usecase.ErrInvalidStatusTransition {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid status transition"})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Data: pkg})
}

// DeletePackage deletes a package
// @Summary Delete a package
// @Description Delete a package from the system
// @Tags packages
// @Param id path string true "Package ID"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /packages/{id} [delete]
func (h *PackageHandler) DeletePackage(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid package ID"})
		return
	}

	err = h.packageUsecase.DeletePackage(id)
	if err != nil {
		if err == usecase.ErrPackageNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "Package not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetPackageStats gets package statistics
// @Summary Get package statistics
// @Description Get aggregated statistics for all packages
// @Tags packages
// @Produce json
// @Success 200 {object} domain.PackageStats
// @Failure 500 {object} ErrorResponse
// @Router /packages/stats [get]
func (h *PackageHandler) GetPackageStats(c *gin.Context) {
	stats, err := h.packageUsecase.GetPackageStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Data: stats})
}

// Response models
type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Data interface{} `json:"data"`
}

type PackageListResponse struct {
	Data   []*domain.Package `json:"data"`
	Limit  int               `json:"limit"`
	Offset int               `json:"offset"`
	Count  int               `json:"count"`
}
