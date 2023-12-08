package handler

import (
	"time"
	"encoding/json"
	"net/http"
	"strconv"
	"os"
	"os/signal"
	"syscall"
	"context"
	"fmt"

	"github.com/gorilla/mux"

	"github.com/go-rbac/internal/service"
	"github.com/go-rbac/internal/core"
	"github.com/aws/aws-xray-sdk-go/xray"

)

type HttpWorkerAdapter struct {
	workerService 	*service.RBACService
}

func NewHttpWorkerAdapter(workerService *service.RBACService) *HttpWorkerAdapter {
	childLogger.Debug().Msg("NewHttpWorkerAdapter")
	return &HttpWorkerAdapter{
		workerService: workerService,
	}
}

type HttpServer struct {
	httpAppServer 	core.HttpAppServer
}

func NewHttpAppServer(httpAppServer core.HttpAppServer) HttpServer {
	childLogger.Debug().Msg("NewHttpAppServer")

	return HttpServer{	
						httpAppServer: httpAppServer,
					}
}

func (h HttpServer) StartHttpAppServer(ctx context.Context, httpWorkerAdapter *HttpWorkerAdapter) {
	childLogger.Info().Msg("StartHttpAppServer")
		
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		childLogger.Debug().Msg("/")
		json.NewEncoder(rw).Encode(h.httpAppServer)
	})

	myRouter.HandleFunc("/info", func(rw http.ResponseWriter, req *http.Request) {
		childLogger.Debug().Msg("/info")
		json.NewEncoder(rw).Encode(h.httpAppServer)
	})
	
	health := myRouter.Methods(http.MethodGet, http.MethodOptions).Subrouter()
    health.HandleFunc("/health", httpWorkerAdapter.Health)

	live := myRouter.Methods(http.MethodGet, http.MethodOptions).Subrouter()
    live.HandleFunc("/live", httpWorkerAdapter.Live)

	header := myRouter.Methods(http.MethodGet, http.MethodOptions).Subrouter()
    header.HandleFunc("/header", httpWorkerAdapter.Header)

	addRole := myRouter.Methods(http.MethodPost, http.MethodOptions).Subrouter()
	addRole.Handle("/role/add", 
						xray.Handler(xray.NewFixedSegmentNamer(fmt.Sprintf("%s%s%s", "rbac:", h.httpAppServer.InfoPod.AvailabilityZone, "./role/add")), 
						http.HandlerFunc(httpWorkerAdapter.AddRole),
						),
	)

	getRole := myRouter.Methods(http.MethodGet, http.MethodOptions).Subrouter()
	getRole.Handle("/role/get/{id}",
						xray.Handler(xray.NewFixedSegmentNamer(fmt.Sprintf("%s%s%s", "rbac:", h.httpAppServer.InfoPod.AvailabilityZone, "./role/getId")),
						http.HandlerFunc(httpWorkerAdapter.GetRole),
						),
	)

	addPolicy := myRouter.Methods(http.MethodPost, http.MethodOptions).Subrouter()
	addPolicy.Handle("/policy/add", 
						xray.Handler(xray.NewFixedSegmentNamer(fmt.Sprintf("%s%s%s", "rbac:", h.httpAppServer.InfoPod.AvailabilityZone, "./policy/add")), 
						http.HandlerFunc(httpWorkerAdapter.AddPolicy),
						),
	)

	getPolicy := myRouter.Methods(http.MethodGet, http.MethodOptions).Subrouter()
	getPolicy.Handle("/policy/get/{id}",
						xray.Handler(xray.NewFixedSegmentNamer(fmt.Sprintf("%s%s%s", "rbac:", h.httpAppServer.InfoPod.AvailabilityZone, "./policy/getId")),
						http.HandlerFunc(httpWorkerAdapter.GetPolicy),
						),
	)
	srv := http.Server{
		Addr:         ":" +  strconv.Itoa(h.httpAppServer.Server.Port),      	
		Handler:      myRouter,                	          
		ReadTimeout:  time.Duration(h.httpAppServer.Server.ReadTimeout) * time.Second,   
		WriteTimeout: time.Duration(h.httpAppServer.Server.WriteTimeout) * time.Second,  
		IdleTimeout:  time.Duration(h.httpAppServer.Server.IdleTimeout) * time.Second, 
	}

	childLogger.Info().Str("Service Port : ", strconv.Itoa(h.httpAppServer.Server.Port)).Msg("Service Port")

	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			childLogger.Error().Err(err).Msg("Cancel http mux server !!!")
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch

	if err := srv.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
		childLogger.Error().Err(err).Msg("WARNING Dirty Shutdown !!!")
		return
	}
	childLogger.Info().Msg("Stop Done !!!!")
}