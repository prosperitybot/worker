package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prosperitybot/worker/domain"
	"github.com/prosperitybot/worker/services"
	"github.com/rs/cors"
	muxtrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gorilla/mux"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
)

// func init() {
// 	prometheus.Register(totalRequests)
// 	prometheus.Register(responseStatus)
// 	prometheus.Register(httpDuration)
// }

func Start() {

	// Init validator and register custom validation functions
	// internal.InitValidator()
	tracer.Start()
	defer tracer.Stop()

	if err := profiler.Start(
		profiler.WithProfileTypes(
			profiler.CPUProfile,
			profiler.HeapProfile,

			// The profiles below are disabled by
			// default to keep overhead low, but
			// can be enabled as needed.
			// profiler.BlockProfile,
			// profiler.MutexProfile,
			// profiler.GoroutineProfile,
		),
	); err != nil {
		log.Fatal(err)
	}
	defer profiler.Stop()

	router := muxtrace.NewRouter()
	testHandler := TestHandler{services.NewTestService(domain.NewTestRepositoryDatabase())}

	// Adding Middleware
	router.Use(NewMiddleware())

	// Handle Prometheus metrics endpoint
	// router.Path("/metrics").Handler(promhttp.Handler())

	// Loading Router
	HandleRequests(router, testHandler)

	handler := cors.Default().Handler(router)

	srv := &http.Server{
		Addr:    ":" + os.Getenv("SERVER_PORT"),
		Handler: handler,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			// logger.Error(context.Background(), fmt.Sprintf("An error occurred while serving the http server: %v", err))
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	// logger.Error(context.Background(), "Shutting down the http server")
	err := srv.Shutdown(ctx)
	if err != nil {
		// logger.Error(context.Background(), fmt.Sprintf("An error occurred while shutting down the http server: %v", err))
	}
	// logger.Error(context.Background(), "Shutting down the database connection")
	// internal.Database.Close()
	os.Exit(0)
}
