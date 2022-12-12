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
	"github.com/mrusme/xbsapi/ent"
	"github.com/mrusme/xbsapi/lib"

	"go.uber.org/zap"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/recover"

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
	var err error
	var xbsctx lib.XBSContext
	var entClient *ent.Client

	entClient, err = ent.Open(config.Database.Type, config.Database.Connection)
	if err != nil {
		logger.Error(
			"Failed initializing database",
			zap.Error(err),
		)
	}
	defer entClient.Close()
	if err := entClient.Schema.Create(context.Background()); err != nil {
		logger.Error(
			"Failed initializing schema",
			zap.Error(err),
		)
	}

	xbsctx = lib.XBSContext{
		EntClient: entClient,
		Logger:    logger,
	}

	fiberApp = fiber.New(fiber.Config{
		Prefork:                 config.Server.Prefork,
		ServerHeader:            config.Server.ServerHeader,
		StrictRouting:           config.Server.StrictRouting,
		CaseSensitive:           config.Server.CaseSensitive,
		ETag:                    config.Server.ETag,
		Concurrency:             config.Server.Concurrency,
		ProxyHeader:             config.Server.ProxyHeader,
		EnableTrustedProxyCheck: config.Server.EnableTrustedProxyCheck,
		TrustedProxies:          config.Server.TrustedProxies,
		DisableStartupMessage:   config.Server.DisableStartupMessage,
		AppName:                 config.Server.AppName,
		ReduceMemoryUsage:       config.Server.ReduceMemoryUsage,
		Network:                 config.Server.Network,
		EnablePrintRoutes:       config.Server.EnablePrintRoutes,
	})
	logger.Info(
		"initialized fiber",
		zap.Any("config", config.Server),
	)

	fiberApp.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(ctx *fiber.Ctx, e interface{}) {
			logger.Error(
				"PANIC",
				zap.Any("error", e),
			)
		},
	}))

	fiberApp.Use(fiberzap.New(fiberzap.Config{
		Logger: logger,
	}))
	logger.Info(
		"initialized logger middleware",
	)

	fiberApp.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // TODO: Config
	}))
	logger.Info(
		"initialized compress middleware",
	)

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
