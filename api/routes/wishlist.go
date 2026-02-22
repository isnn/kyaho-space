package routes

import (
	"kyaho-space/api/handlers"
	"kyaho-space/pkg/wishlist"

	"github.com/gofiber/fiber/v3"
)

func WishlistRouter(api fiber.Router, service wishlist.Service) {
	api.Post("/wishlists", handlers.CreateWishlist(service))
	api.Get("/wishlists", handlers.ListWishlists(service))
	api.Get("/wishlists/:id", handlers.GetDetailWishlist(service))
	api.Put("/wishlists/:id", handlers.PutWishlist(service))
	api.Delete("/wishlists/:id", handlers.DeleteWishlist(service))
}
