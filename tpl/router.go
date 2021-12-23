package tpl

const (
	RegisterRouter = `router{{CammelCaseName}}(app)
	// router:register`

	HandlerInterface = `package {{PackageName}}

import "github.com/labstack/echo/v4"

// Handler defines all methods attached to the routes
type Handler interface {
	// Create handler for POST /{{SnakeCaseName}} route
	Create(ec echo.Context) error

	// Update handler for PUT /{{SnakeCaseName}}/:id route
	Update(ec echo.Context) error

	// Delete handler for DELETE /{{SnakeCaseName}}/:id route
	Delete(ec echo.Context) error

	// FindOne handler for GET /{{SnakeCaseName}}/:id route
	FindOne(ec echo.Context) error

	// Find handler for GET /{{SnakeCaseName}} route
	Find(ec echo.Context) error
}`

	HandlerStruct = `package {{PackageName}}

import (
	"github.com/Drafteame/framework/logger"
	"github.com/labstack/echo/v4"
)

type handler struct {
	logger logger.Logger
}

var _ Handler = &handler{}

func NewHandler(l logger.Logger) Handler {
	return &handler{
		logger: l,
	}
}

func (h *handler) Create(ec echo.Context) error {
	return ec.JSON(200, map[string]interface{}{
		"method": "create",
		"resource": "{{SnakeCaseName}}",
	})
}

func (h *handler) Update(ec echo.Context) error {
	return ec.JSON(200, map[string]interface{}{
		"method": "update",
		"resource": "{{SnakeCaseName}}",
	})
}

func (h *handler) Delete(ec echo.Context) error {
	return ec.JSON(200, map[string]interface{}{
		"method": "delete",
		"resource": "{{SnakeCaseName}}",
	})
}

func (h *handler) FindOne(ec echo.Context) error {
	return ec.JSON(200, map[string]interface{}{
		"method": "findOne",
		"resource": "{{SnakeCaseName}}",
	})
}

func (h *handler) Find(ec echo.Context) error {
	return ec.JSON(200, map[string]interface{}{
		"method": "find",
		"resource": "{{SnakeCaseName}}",
	})
}`

	JSONSchemas = `package schemas

import "github.com/Drafteame/framework/types"

var (
	Create{{CammelCaseName}} = types.JSONSchema{
		Type: "object",
		AdditionalProperties: true,
	}

	Update{{CammelCaseName}} = types.JSONSchema{
		Type: "object",
		AdditionalProperties: true,
	}
)`

	Router = `package routes

import (
	"github.com/Drafteame/framework/engine"
	em "github.com/Drafteame/framework/middlewares"
	handlers "{{Namespace}}/internal/handlers/{{PackageName}}"
	"{{Namespace}}/internal/schemas"
)

func router{{CammelCaseName}}(app engine.App) {
	handler := handlers.NewHandler(app.Logger())

	app.POST("/{{SnakeCaseName}}", handler.Create, em.JSONSchema(schemas.Create{{CammelCaseName}}))
	app.GET("/{{SnakeCaseName}}", handler.Find)
	app.PUT("/{{SnakeCaseName}}/:id", handler.Update, em.JSONSchema(schemas.Update{{CammelCaseName}}))
	app.DELETE("/{{SnakeCaseName}}/:id", handler.Delete)
	app.GET("/{{SnakeCaseName}}/:id", handler.FindOne)
}`
)
