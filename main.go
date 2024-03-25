package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/lightwell/cms_template_go_v2/controller"
	"github.com/lightwell/cms_utilities_go_v2/cloud"
	"github.com/lightwell/cms_utilities_go_v2/cms"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/rs/cors"

	"github.com/joho/godotenv"
	auth "github.com/lightwell/cms_utilities_go_v2/auth"
	"github.com/lightwell/cms_utilities_go_v2/config"
	"github.com/lightwell/cms_utilities_go_v2/db"
	"github.com/lightwell/cms_utilities_go_v2/models"
)

// https://docs.aws.amazon.com/elasticbeanstalk/latest/dg/https-singleinstance-go.html
const (
	certPath    = "/etc/pki/tls/certs/server.crt"
	privKeyPath = "/etc/pki/tls/certs/server.key"
)

func main() {
	var err error
	godotenv.Load(".env")

	// Setup
	isDev, err := strconv.ParseBool(os.Getenv("IS_DEV"))
	if err != nil {
		panic(err)
	}
	err = config.LoadConfig(isDev)
	if err != nil {
		panic(err)
	}

	err = models.LoadModels()
	if err != nil {
		panic(err)
	}

	err = db.CreateCon(isDev)
	if err != nil {
		panic(err)
	}
	defer db.CloseCon()

	err = cloud.InitS3()
	if err != nil {
		panic(err)
	}

	// fs := http.FileServer(http.Dir("./_frontend/dist"))
	// http.Handle("/static/", http.StripPrefix("/static/", fs))

	go setupFrontend()

	// Run pending migrations
	output, err := db.RunMigrations(isDev)
	if err != nil {
		log.Println(err)
	}
	log.Print(output)

	// Initialise rowToModelLookups and customEndpoints
	controller.InitialiseLookupsEndpoints()

	// Register "Pages" and "Handlers"
	auth.Setup()
	mux := cms.RegisterHandlers()

	handler := cors.Default().Handler(mux)

	c := cors.New(cors.Options{
		AllowedOrigins:   config.GetCORSOptions().AllowedOrigins,
		AllowCredentials: true,
		Debug:            true,
		AllowedMethods:   config.GetCORSOptions().AllowedMethods,
		AllowedHeaders:   config.GetCORSOptions().AllowedHeaders,
	})

	// Get Port for Server and Run it
	port := config.GetHosting().BackendPort // Set port to 5000 if running on AWS or locally
	if port < 5000 || port > 8090 {
		log.Fatal("port must be in range 5000 <= port <= 8090")
	}

	isHTTPS, _ := strconv.ParseBool(os.Getenv("IS_HTTPS"))
	if isHTTPS {
		port = 443
		cms.IsHTTPS = true

		// Setup HTTPS connection
		log.Printf("Running on port %d...", port)
		log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%d", port), certPath, privKeyPath, c.Handler(handler)))
	} else {
		// Setup HTTP connection
		log.Printf("Running on port %d...", port)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), c.Handler(handler)))
	}
}

func setupFrontend() {

	// React Front end setup
	engine := html.New("./react_frontend/build/", ".html")
	app := fiber.New(fiber.Config{Views: engine})

	// Setup React Routes
	for _, route := range config.GetFrontendRoutes() {
		app.Static(route, "./react_frontend/build/")
		app.Get(route, home)
	}

	app.Listen(fmt.Sprintf(":%d", config.GetHosting().FrontendPort))
}

func home(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{})
}
