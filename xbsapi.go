package main

import (
	"context"
	"embed"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	fiberadapter "github.com/awslabs/aws-lambda-go-api-proxy/fiber"
	"github.com/mrusme/journalist/middlewares/fiberzap"
	"github.com/mrusme/xbsapi/api"
	"github.com/mrusme/xbsapi/lib"

	"go.uber.org/zap"

	"github.com/gofiber/fiber/v2"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed favicon.ico
var favicon embed.FS

var fiberApp *fiber.App
var fiberLambda *fiberadapter.FiberLambda

var config lib.Config
var logger *zap.Logger

func init() {
	var err error

	fiberLambda = fiberadapter.New(fiberApp)
	config, err = lib.Cfg()
	if err != nil {
		panic(err)
	}

	if config.Debug == "true" {
		logger, _ = zap.NewDevelopment()
	} else {
		logger, _ = zap.NewProduction()
	}
	defer logger.Sync()
	// TODO: Use sugarLogger
	// sugar := logger.Sugar()
}

func AWSLambdaHandler(
	ctx context.Context,
	req events.APIGatewayProxyRequest,
) (events.APIGatewayProxyResponse, error) {
	return fiberLambda.ProxyWithContext(ctx, req)
}

func GCFHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	err := CloudFunctionRouteToFiber(fiberApp, w, r)
	if err != nil {
		logger.Error(
			"Handler error",
			zap.Error(err),
		)
		return
	}
}

func main() {
	// var err error
	var xbsctx lib.XBSContext
	// var entClient *ent.Client

	// entClient, err = ent.Open(config.Database.Type, config.Database.Connection)
	// if err != nil {
	// 	logger.Error(
	// 		"Failed initializing database",
	// 		zap.Error(err),
	// 	)
	// }
	// defer entClient.Close()
	// if err := entClient.Schema.Create(context.Background()); err != nil {
	// 	logger.Error(
	// 		"Failed initializing schema",
	// 		zap.Error(err),
	// 	)
	// }

	xbsctx = lib.XBSContext{
		Logger: logger,
	}

	fiberApp = fiber.New(fiber.Config{
		Prefork:                 false, // TODO: Make configurable
		ServerHeader:            "",    // TODO: Make configurable
		StrictRouting:           false,
		CaseSensitive:           false,
		ETag:                    false,      // TODO: Make configurable
		Concurrency:             256 * 1024, // TODO: Make configurable
		ProxyHeader:             "",         // TODO: Make configurable
		EnableTrustedProxyCheck: false,      // TODO: Make configurable
		TrustedProxies:          []string{}, // TODO: Make configurable
		DisableStartupMessage:   true,
		AppName:                 "xbsapi",
		ReduceMemoryUsage:       false,            // TODO: Make configurable
		Network:                 fiber.NetworkTCP, // TODO: Make configurable
		EnablePrintRoutes:       false,
	})
	fiberApp.Use(fiberzap.New(fiberzap.Config{
		Logger: logger,
	}))

	api.Register(
		&xbsctx,
		fiberApp,
	)

	fiberApp.Get("/favicon.ico", func(ctx *fiber.Ctx) error {
		fi, err := favicon.Open("favicon.ico")
		if err != nil {
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}
		return ctx.SendStream(fi)
	})

	functionName := os.Getenv("AWS_LAMBDA_FUNCTION_NAME")
	if functionName == "" {
		listenAddr := fmt.Sprintf(
			"%s:%s",
			config.Server.BindIP,
			config.Server.Port,
		)
		logger.Fatal(
			"Server failed",
			zap.Error(fiberApp.Listen(listenAddr)),
		)
	} else {
		lambda.Start(AWSLambdaHandler)
	}
}
