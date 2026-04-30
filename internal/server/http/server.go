package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"subscriptions_rest/internal/domain"
	"subscriptions_rest/internal/repository"

	"github.com/valyala/fasthttp"
)

type Server struct {
	repository repository.Repository
	server     *fasthttp.Server
}

func New(repository repository.Repository) *Server {
	server := &fasthttp.Server{}
	return &Server{
		repository: repository,
		server:     server,
	}
}

func (s *Server) Start(port int) error {
	s.server.Handler = s.getSubscribe
	log.Printf("start listening http-server at %d port\n", port)
	return s.server.ListenAndServe(fmt.Sprintf("localhost:%d", port))
}

func (s *Server) Stop() error {
	log.Println("stop http-server")
	return s.server.Shutdown()
}

func (s *Server) getSubscribe(ctx *fasthttp.RequestCtx) {
	if string(ctx.Request.Header.Method()) != http.MethodGet || string(ctx.Request.URI().Path()) != "/subscribes" {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}

	var request domain.RepositoryRequest
	err := json.Unmarshal(ctx.Request.Body(), &request)
	if err != nil {
		log.Printf("failed to parse json: %s\n", err.Error())
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	// if request.ClientId == "" || request.ServiceName == "" {
	// 	log.Printf("invalid data\n")
	// 	ctx.SetStatusCode(fasthttp.StatusBadRequest)
	// 	return
	// }

	subscribe, err := s.repository.GetSubscribe(request.ClientId, request.ServiceName)
	if err != nil {
		log.Printf("failed to get subscribe: %s\n", err.Error())
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	responseBody, err := json.Marshal(subscribe)
	if err != nil {
		log.Printf("failed to marshal data: %s\n", err.Error())
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(responseBody)
}

func (s *Server) getSubscribesPrice(ctx *fasthttp.RequestCtx) {
	if string(ctx.Request.Header.Method()) != http.MethodGet || string(ctx.Request.URI().Path()) != "/subscribes/prices" {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}

	var request domain.RepositoryRequest
	err := json.Unmarshal(ctx.Request.Body(), &request)
	if err != nil {
		log.Printf("failed to parse json: %s\n", err.Error())
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	// if request.ClientId == "" || request.ServiceName == "" || request.Start.Equal(time.Time{}) {
	// 	log.Printf("invalid data\n")
	// 	ctx.SetStatusCode(fasthttp.StatusBadRequest)
	// 	return
	// }

	price, err := s.repository.GetSubscribesPrice(request)
	if err != nil {
		log.Printf("failed to get price: %s\n", err.Error())
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	responseBody, err := json.Marshal(domain.PriceResponse{
		Price: price,
	})
	if err != nil {
		log.Printf("failed to marshal data: %s\n", err.Error())
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(responseBody)
}

func (s *Server) createSubscribe(ctx *fasthttp.RequestCtx) {
	if string(ctx.Request.Header.Method()) != http.MethodPost || string(ctx.Request.URI().Path()) != "/subscribes" {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}

	// TODO
	var request domain.Subscribe
	err := json.Unmarshal(ctx.Request.Body(), &request)
	if err != nil {
		log.Printf("failed to parse json: %s\n", err.Error())
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	// if request.ClientId == "" || request.ServiceName == "" || request.Start.Equal(time.Time{}) {
	// 	log.Printf("invalid data\n")
	// 	ctx.SetStatusCode(fasthttp.StatusBadRequest)
	// 	return
	// }

	err = s.repository.AddSubscribe(request)
	if err != nil {
		log.Printf("failed to add subscribe: %s\n", err.Error())
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusCreated)
}

func (s *Server) updateSubscribe(ctx *fasthttp.RequestCtx) {
	if string(ctx.Request.Header.Method()) != http.MethodPut || string(ctx.Request.URI().Path()) != "/subscribes" {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}

	// TODO
	var request domain.Subscribe
	err := json.Unmarshal(ctx.Request.Body(), &request)
	if err != nil {
		log.Printf("failed to parse json: %s\n", err.Error())
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	// if request.ClientId == "" || request.ServiceName == "" || request.Start.Equal(time.Time{}) {
	// 	log.Printf("invalid data\n")
	// 	ctx.SetStatusCode(fasthttp.StatusBadRequest)
	// 	return
	// }

	err = s.repository.UpdateSubscribe(request)
	if err != nil {
		log.Printf("failed to update subscribe: %s\n", err.Error())
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
}

func (s *Server) deleteSubscribe(ctx *fasthttp.RequestCtx) {
	if string(ctx.Request.Header.Method()) != http.MethodDelete || string(ctx.Request.URI().Path()) != "/subscribes" {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}

	// TODO
	var request domain.RepositoryRequest
	err := json.Unmarshal(ctx.Request.Body(), &request)
	if err != nil {
		log.Printf("failed to parse json: %s\n", err.Error())
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	// if request.ClientId == "" || request.ServiceName == "" || request.Start.Equal(time.Time{}) {
	// 	log.Printf("invalid data\n")
	// 	ctx.SetStatusCode(fasthttp.StatusBadRequest)
	// 	return
	// }

	err = s.repository.DeleteSubscribe(request.ClientId, request.ServiceName)
	if err != nil {
		log.Printf("failed to delete subscribe: %s\n", err.Error())
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
}
