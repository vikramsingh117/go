package routes

import (
    "github.com/gofiber/fiber/v2"
    "wav-to-flac-service/handlers"
)

func SetupRoutes(app *fiber.App) {
    app.Get("/", func(c *fiber.Ctx) error {
        return c.SendString("WAV to FLAC Audio Conversion Service")
    })

    // WebSocket route for streaming
    app.Use("/ws", handlers.UpgradeHandler)
}
