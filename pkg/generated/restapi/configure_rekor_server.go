/*
Copyright © 2020 The Sigstore Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"context"
	"crypto/tls"
	"net/http"
	"strconv"
	"time"

	// using embed to add the static html page duing build time
	_ "embed"

	"github.com/go-chi/chi/middleware"
	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/cors"
	"github.com/spf13/viper"

	pkgapi "github.com/sigstore/rekor/pkg/api"
	"github.com/sigstore/rekor/pkg/generated/restapi/operations"
	"github.com/sigstore/rekor/pkg/generated/restapi/operations/entries"
	"github.com/sigstore/rekor/pkg/generated/restapi/operations/index"
	"github.com/sigstore/rekor/pkg/generated/restapi/operations/pubkey"
	"github.com/sigstore/rekor/pkg/generated/restapi/operations/tlog"
	"github.com/sigstore/rekor/pkg/log"
	"github.com/sigstore/rekor/pkg/util"

	"golang.org/x/exp/slices"
)

//go:generate swagger generate server --target ../../generated --name RekorServer --spec ../../../openapi.yaml --principal interface{} --exclude-main

type contextKey string

var (
	ctxKeyAPIToRecord = contextKey("apiToRecord")
)

// Context payload for recording metrics.
type apiToRecord struct {
	method *string // Method to record in metrics, if any.
	path   *string // Path to record in metrics, if any.
}

func configureFlags(api *operations.RekorServerAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.RekorServerAPI) http.Handler {
	// configure the api here
	api.ServeError = logAndServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf
	api.Logger = log.Logger.Infof

	// api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()
	api.JSONProducer = runtime.JSONProducer()

	api.ApplicationXPemFileProducer = runtime.TextProducer()

	// disable all endpoints to start
	api.IndexSearchIndexHandler = index.SearchIndexHandlerFunc(pkgapi.SearchIndexNotImplementedHandler)
	api.EntriesCreateLogEntryHandler = entries.CreateLogEntryHandlerFunc(pkgapi.CreateLogEntryNotImplementedHandler)
	api.EntriesGetLogEntryByIndexHandler = entries.GetLogEntryByIndexHandlerFunc(pkgapi.GetLogEntryByIndexNotImplementedHandler)
	api.EntriesGetLogEntryByUUIDHandler = entries.GetLogEntryByUUIDHandlerFunc(pkgapi.GetLogEntryByUUIDNotImplementedHandler)
	api.EntriesSearchLogQueryHandler = entries.SearchLogQueryHandlerFunc(pkgapi.SearchLogQueryNotImplementedHandler)
	api.PubkeyGetPublicKeyHandler = pubkey.GetPublicKeyHandlerFunc(pkgapi.GetPublicKeyNotImplementedHandler)
	api.TlogGetLogProofHandler = tlog.GetLogProofHandlerFunc(pkgapi.GetLogProofNotImplementedHandler)

	enabledAPIEndpoints := viper.GetStringSlice("enabled_api_endpoints")
	if !slices.Contains(enabledAPIEndpoints, "searchIndex") && viper.GetBool("enable_retrieve_api") {
		enabledAPIEndpoints = append(enabledAPIEndpoints, "searchIndex")
	}

	for _, enabledAPI := range enabledAPIEndpoints {
		log.Logger.Infof("Enabling API endpoint: %s", enabledAPI)
		switch enabledAPI {
		case "searchIndex":
			api.IndexSearchIndexHandler = index.SearchIndexHandlerFunc(pkgapi.SearchIndexHandler)
		case "getLogInfo":
			api.TlogGetLogInfoHandler = tlog.GetLogInfoHandlerFunc(pkgapi.GetLogInfoHandler)
		case "getPublicKey":
			api.PubkeyGetPublicKeyHandler = pubkey.GetPublicKeyHandlerFunc(pkgapi.GetPublicKeyHandler)
		case "getLogProof":
			api.TlogGetLogProofHandler = tlog.GetLogProofHandlerFunc(pkgapi.GetLogProofHandler)
		case "createLogEntry":
			api.EntriesCreateLogEntryHandler = entries.CreateLogEntryHandlerFunc(pkgapi.CreateLogEntryHandler)
		case "getLogEntryByIndex":
			api.EntriesGetLogEntryByIndexHandler = entries.GetLogEntryByIndexHandlerFunc(pkgapi.GetLogEntryByIndexHandler)
		case "getLogEntryByUUID":
			api.EntriesGetLogEntryByUUIDHandler = entries.GetLogEntryByUUIDHandlerFunc(pkgapi.GetLogEntryByUUIDHandler)
		case "searchLogQuery":
			api.EntriesSearchLogQueryHandler = entries.SearchLogQueryHandlerFunc(pkgapi.SearchLogQueryHandler)
		default:
			log.Logger.Panicf("Unknown API endpoint requested: %s", enabledAPI)
		}
	}

	// all handlers need to be set before a call to api.AddMiddlewareFor
	for _, enabledAPI := range enabledAPIEndpoints {
		switch enabledAPI {
		case "searchIndex":
			recordMetricsForAPI(api, "POST", "/api/v1/index/retrieve") // add metrics
		case "getLogInfo":
			api.AddMiddlewareFor("GET", "/api/v1/log", middleware.NoCache) // not cacheable
			recordMetricsForAPI(api, "GET", "/api/v1/log")                 // add metrics
		case "getPublicKey":
			api.AddMiddlewareFor("GET", "/api/v1/log/publicKey", middleware.NoCache) // not cacheable
			recordMetricsForAPI(api, "GET", "/api/v1/log/publicKey")                 // add metrics
		case "getLogProof":
			api.AddMiddlewareFor("GET", "/api/v1/log/proof", middleware.NoCache) // not cacheable
			recordMetricsForAPI(api, "GET", "/api/v1/log/proof")                 // add metrics
		case "createLogEntry":
			recordMetricsForAPI(api, "POST", "/api/v1/log/entries") // add metrics
		case "getLogEntryByIndex":
			api.AddMiddlewareFor("GET", "/api/v1/log/entries", middleware.NoCache) // not cacheable
			recordMetricsForAPI(api, "GET", "/api/v1/log/entries")                 // add metrics
		case "getLogEntryByUUID":
			api.AddMiddlewareFor("GET", "/api/v1/log/entries/{entryUUID}", middleware.NoCache) // not cacheable
			recordMetricsForAPI(api, "GET", "/api/v1/log/entries/{entryUUID}")                 // add metrics
		case "searchLogQuery":
			recordMetricsForAPI(api, "POST", "/api/v1/log/entries/retrieve") // add metrics
		}
	}
	api.RegisterFormat("signedCheckpoint", &util.SignedNote{}, util.SignedCheckpointValidator)

	api.PreServerShutdown = func() {}
	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// We need this type to act as an adapter between zap and the middleware request logger.
type logAdapter struct {
}

func (l *logAdapter) Print(v ...interface{}) {
	log.Logger.Info(v...)
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	middleware.DefaultLogger = middleware.RequestLogger(
		&middleware.DefaultLogFormatter{Logger: &logAdapter{}})
	returnHandler := middleware.Logger(handler)
	returnHandler = middleware.Recoverer(returnHandler)
	returnHandler = middleware.Heartbeat("/ping")(returnHandler)
	returnHandler = serveStaticContent(returnHandler)

	handleCORS := cors.Default().Handler
	returnHandler = handleCORS(returnHandler)

	returnHandler = wrapMetrics(returnHandler)

	return middleware.RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		r = r.WithContext(log.WithRequestID(ctx, middleware.GetReqID(ctx)))
		defer func() {
			_ = log.ContextLogger(ctx).Sync()
		}()

		returnHandler.ServeHTTP(w, r)
	}))
}

// Populates the the apiToRecord for this method/path so metrics are emitted.
func recordMetricsForAPI(api *operations.RekorServerAPI, method string, path string) {
	metricsHandler := func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			if apiInfo, ok := ctx.Value(ctxKeyAPIToRecord).(*apiToRecord); ok {
				apiInfo.method = &method
				apiInfo.path = &path
			} else {
				log.ContextLogger(ctx).Warn("Could not attach api info - endpoint may not be monitored.")
			}
			handler.ServeHTTP(w, r)
		})
	}

	api.AddMiddlewareFor(method, path, metricsHandler)
}

func wrapMetrics(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		apiInfo := apiToRecord{}
		ctx = context.WithValue(ctx, ctxKeyAPIToRecord, &apiInfo)
		r = r.WithContext(ctx)

		start := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		defer func() {
			// Only record metrics for APIs that need instrumentation.
			if apiInfo.path != nil && apiInfo.method != nil {
				code := strconv.Itoa(ww.Status())
				labels := map[string]string{
					"path": *apiInfo.path,
					"code": code,
				}
				// This logs latency broken down by URL path and response code
				// TODO(var-sdk): delete these metrics once the new metrics are safely rolled out.
				pkgapi.MetricLatency.With(labels).Observe(float64(time.Since(start)))
				pkgapi.MetricLatencySummary.With(labels).Observe(float64(time.Since(start)))

				pkgapi.MetricRequestLatency.With(
					map[string]string{
						"path":   *apiInfo.path,
						"method": *apiInfo.method,
					}).Observe(float64(time.Since(start)))

				pkgapi.MetricRequestCount.With(
					map[string]string{
						"path":   *apiInfo.path,
						"method": *apiInfo.method,
						"code":   code,
					}).Inc()
			}
		}()

		handler.ServeHTTP(ww, r)

	})
}

func logAndServeError(w http.ResponseWriter, r *http.Request, err error) {
	ctx := r.Context()
	if apiErr, ok := err.(errors.Error); ok && apiErr.Code() == http.StatusNotFound {
		log.ContextLogger(ctx).Warn(err)
	} else {
		log.ContextLogger(ctx).Error(err)
	}
	requestFields := map[string]interface{}{}
	if err := mapstructure.Decode(r, &requestFields); err == nil {
		log.ContextLogger(ctx).Debug(requestFields)
	}
	errors.ServeError(w, r, err)
}

//go:embed rekorHomePage.html
var homePageBytes []byte

func serveStaticContent(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			w.Header().Add("Content-Type", "text/html")
			w.WriteHeader(200)
			_, _ = w.Write(homePageBytes)
			return
		}
		handler.ServeHTTP(w, r)
	})
}
