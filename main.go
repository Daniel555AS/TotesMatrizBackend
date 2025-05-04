package main

import (
	"totesbackend/app"
	_ "totesbackend/docs"
)

// @title           Totes Backend API
// @version         1.0
// @description     This is a sample API documentation for the Totes backend using Gin and Swagger.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@totesbackend.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost
// @BasePath  /

// @schemes http https
func main() {
	// Load environment variables and run the application
	app.SetupAndRunApp()
}
