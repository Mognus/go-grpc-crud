package crud

import "github.com/gofiber/fiber/v2"

// RouteConfig defines which middleware to apply per HTTP verb for a CRUDProvider.
// A nil/empty middleware slice means no protection for that verb.
type RouteConfig struct {
	Provider CRUDProvider
	Path     string         // URL path segment, defaults to Provider.GetModelName()
	List     []fiber.Handler
	Get      []fiber.Handler
	Create   []fiber.Handler
	Update   []fiber.Handler
	Delete   []fiber.Handler
}

// Public returns a RouteConfig with no middleware on any route.
func Public(provider CRUDProvider) RouteConfig {
	return RouteConfig{Provider: provider}
}

// Protected returns a RouteConfig applying the same middleware to all routes.
func Protected(provider CRUDProvider, middlewares ...fiber.Handler) RouteConfig {
	return RouteConfig{
		Provider: provider,
		List:     middlewares,
		Get:      middlewares,
		Create:   middlewares,
		Update:   middlewares,
		Delete:   middlewares,
	}
}

// ReadPublic returns a RouteConfig where read routes (List, Get) are public
// and write routes (Create, Update, Delete) require the given middlewares.
func ReadPublic(provider CRUDProvider, writeMiddlewares ...fiber.Handler) RouteConfig {
	return RouteConfig{
		Provider: provider,
		Create:   writeMiddlewares,
		Update:   writeMiddlewares,
		Delete:   writeMiddlewares,
	}
}
