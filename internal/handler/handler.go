package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/yomek33/talki/internal/services"
)

const frontendURI = "http://localhost:5173"

type Handlers struct {
	UserHandler
	MaterialHandler
	PhraseHandler
	ChatHandler
	jwtSecretKey string
	Firebase     *Firebase
}

func NewHandler(s *services.Services, jwtSecretKey string, firebase *Firebase) *Handlers {
	return &Handlers{
		UserHandler:     &userHandler{UserService: s.UserService, jwtSecretKey: jwtSecretKey, Firebase: firebase},
		MaterialHandler: &materialHandler{MaterialService: s.MaterialService, PhraseService: s.PhraseService},
		PhraseHandler:   &phraseHandler{PhraseService: s.PhraseService},
		ChatHandler:     &chatHandler{ChatService: s.ChatService},
		jwtSecretKey:    jwtSecretKey,
		Firebase:        firebase,
	}
}

func (h *Handlers) SetDefault(e *echo.Echo) {
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome to our API")
	})
}

func (h *Handlers) SetAPIRoutes(e *echo.Echo) {
	api := e.Group("/api")

	// Handle OPTIONS request for CORS preflight
	api.OPTIONS("/auth", handleOptions)
	api.POST("/auth", h.GetGoogleLoginSignin)

	materialRoutes := api.Group("/materials")
	materialRoutes.POST("", h.CreateMaterial)
	materialRoutes.GET("", h.GetAllMaterials)
	materialRoutes.GET("/:id", h.GetMaterialByID)
	materialRoutes.PUT("/:id", h.UpdateMaterial)
	materialRoutes.DELETE("/:id", h.DeleteMaterial)
	materialRoutes.GET("/:id/status", h.CheckMaterialStatus)
	materialRoutes.GET("/:id/phrases", h.GetProcessedPhrases)

	userRoutes := api.Group("/users")
	userRoutes.PUT("/:id", h.UpdateUser)
	userRoutes.DELETE("/:id", h.DeleteUser)

	chatRoutes := api.Group("/chat")
	chatRoutes.POST("", h.CreateChat)
	chatRoutes.GET("/:id", h.GetChat)
	chatRoutes.POST("/:id/messages", h.CreateMessage)
	chatRoutes.GET("/:id/messages", h.GetMessages)
}

func handleOptions(c echo.Context) error {
	setCORSHeaders(c)
	return c.NoContent(http.StatusNoContent)
}

func Echo() *echo.Echo {
	e := echo.New()

	// Set up middleware
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} ${method} ${uri} ${status} ${latency_human}\n",
	}))
	e.Use(middleware.Recover())
	e.Pre(middleware.RemoveTrailingSlash())

	// Configure CORS settings
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{frontendURI},
		AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))

	e.Use(middleware.Secure())

	// Custom HTTP error handler
	e.HTTPErrorHandler = customHTTPErrorHandler

	return e
}

func customHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	message := echo.Map{"message": "Internal Server Error"}

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		if message, ok = he.Message.(echo.Map); !ok {
			// he.Message was not of type echo.Map
			if messageStr, ok := he.Message.(string); ok {
				message = echo.Map{"message": messageStr}
			} else {
				message = echo.Map{"message": http.StatusText(code)}
			}
		}
		if he.Internal != nil {
			message = echo.Map{"message": fmt.Sprintf("%v, %v", message, he.Internal)}
		}
	}

	// Log the error
	c.Logger().Error(err)

	// Send JSON response
	if !c.Response().Committed {
		if c.Request().Method == http.MethodHead {
			c.NoContent(code)
		} else {
			c.JSON(code, message)
		}
	}
}

func setCORSHeaders(c echo.Context) {
	c.Response().Header().Set("Access-Control-Allow-Origin", frontendURI)
	c.Response().Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Credentials", "true")
}
