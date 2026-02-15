package panel

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

const panelTestRequestTimeoutMS = 15000

func testFiberRequest(app *fiber.App, req *http.Request) (*http.Response, error) {
	return app.Test(req, panelTestRequestTimeoutMS)
}
