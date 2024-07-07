package utils

import (
	"encoding/json"
	"log"
	"os"

	"ospm/config"
	"ospm/internal/api/routes"
	"ospm/internal/repository/database/cockroachdb"
	OSPMLogger "ospm/internal/service/log"

	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// StartOSPM manages the startup sequence
func StartOSPM() {
	var OSPMWG sync.WaitGroup

	// 1.
	// reading the configs and settings from config.env file
	// and load it to memory so the configs be accessible in the entire program
	config.LoadOSPMConfigs()

	// 2.
	// init the logger
	OSPMLogger.InitLogger()

	configs, _ := json.MarshalIndent(config.OSPM, "", "  ")
	OSPMLogger.Log.Debugf("%+v", string(configs))

	//3.
	// init the database
	cockroachdb.InitialDB()

	//4.
	// starting the api server
	OSPMWG.Add(1)
	StartAPIServer(&OSPMWG)

	OSPMWG.Wait()
}

func StartAPIServer(wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: config.OSPM.API.AllowOrigins,
		AllowMethods: config.OSPM.API.AllowMethods,
		AllowHeaders: config.OSPM.API.AllowHeaders,
	}))

	app.Use(logger.New(logger.Config{
		Format:       "[${time}] status:${status} - Latency: ${latency} Method: ${method} Path: ${path}\n",
		TimeFormat:   time.RFC3339Nano,
		TimeZone:     "Local",
		TimeInterval: 500 * time.Millisecond,
		Output:       os.Stdout,
		Done: func(c *fiber.Ctx, logString []byte) {
			if c.Response().StatusCode() != fiber.StatusOK {
				log.Println(string(logString))
			}
		},
	}))

	routes.Setup(app)

	OSPMLogger.Log.Fatal(app.Listen(config.OSPM.API.GetListenAddress()))

}
