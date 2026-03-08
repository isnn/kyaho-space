package handlers

import (
	. "kyaho-space/pkg/common"
	"kyaho-space/pkg/entities"
	"kyaho-space/pkg/wishlist"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func CreateWishlist(service wishlist.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		var req entities.WishlistRequest
		if err := c.Bind().Body(&req); err != nil {
			return BadRequestResponse(c, "Invalid request body")
		}

		if errs := ValidateStruct(req); errs != nil {
			return ValidationErrorResponse(c, "Validation failed", errs)
		}

		w := entities.Wishlist{
			TargetAction: req.TargetAction,
			TargetYear:   req.TargetYear,
			Position:     req.Position,
			Status:       req.Status,
		}

		if err := service.AddWishlist(&w); err != nil {
			return InternalErrorResponse(c, "Failed to create wishlist")
		}

		return CreatedResponse(c, w)
	}
}

func ListWishlists(service wishlist.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		page, _ := strconv.Atoi(c.Query("page", "1"))
		limit, _ := strconv.Atoi(c.Query("limit", "10"))

		params := ListParams{
			Page:   page,
			Limit:  limit,
			Search: c.Query("search", ""),
			SortBy: c.Query("sort_by", "created_at"),
			Order:  c.Query("order", "desc"),
		}
		year := c.Query("year", "")

		items, total, err := service.GetWishlists(params, year)
		if err != nil {
			return InternalErrorResponse(c, "Failed to retrieve wishlists")
		}

		return SuccessResponse(c, NewPaginatedResponse(params, total, items))
	}
}

func GetDetailWishlist(service wishlist.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return BadRequestResponse(c, "Wishlist ID is required")
		}

		w, err := service.GetWishlistByID(id)
		if err != nil {
			return NotFoundResponse(c, "Wishlist not found")
		}

		return SuccessResponse(c, w)
	}
}

func PutWishlist(service wishlist.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return BadRequestResponse(c, "Wishlist ID is required")
		}

		var req entities.WishlistUpdateRequest
		if err := c.Bind().Body(&req); err != nil {
			return BadRequestResponse(c, "Invalid request body")
		}

		if errs := ValidateStruct(req); errs != nil {
			return ValidationErrorResponse(c, "Validation failed", errs)
		}

		updates := make(map[string]interface{})
		if req.TargetAction != nil {
			updates["target_action"] = *req.TargetAction
		}
		if req.TargetYear != nil {
			updates["target_year"] = *req.TargetYear
		}
		if req.Position != nil {
			updates["position"] = *req.Position
		}
		if req.Status != nil {
			status := *req.Status
			if status == "" {
				status = "planned"
			}
			updates["status"] = status
		}

		if len(updates) == 0 {
			return BadRequestResponse(c, "No fields to update")
		}

		w := entities.Wishlist{ID: id}
		if err := service.UpdateWishlist(&w, updates); err != nil {
			if err == gorm.ErrRecordNotFound {
				return NotFoundResponse(c, "Wishlist not found")
			}
			return InternalErrorResponse(c, "Failed to update wishlist")
		}

		return SuccessResponse(c, fiber.Map{"message": "Wishlist updated successfully"})
	}
}

func DeleteWishlist(service wishlist.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return BadRequestResponse(c, "Wishlist ID is required")
		}

		if err := service.DeleteWishlist(id); err != nil {
			if err == gorm.ErrRecordNotFound {
				return NotFoundResponse(c, "Wishlist not found")
			}
			return InternalErrorResponse(c, "Failed to delete wishlist")
		}

		return SuccessResponse(c, fiber.Map{"message": "Wishlist deleted successfully"})
	}
}
