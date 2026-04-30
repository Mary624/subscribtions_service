package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	_ "subscriptions_rest/docs"
	"subscriptions_rest/internal/config"
	"subscriptions_rest/internal/domain"
	"subscriptions_rest/internal/repository"

	swagger "github.com/swaggo/fasthttp-swagger"

	"github.com/valyala/fasthttp"
)

type Server struct {
	repository repository.Repository
	cfg        *config.Config
	server     *fasthttp.Server
}

func New(cfg *config.Config, repository repository.Repository) *Server {
	server := &fasthttp.Server{}
	return &Server{
		cfg:        cfg,
		repository: repository,
		server:     server,
	}
}

func (s *Server) Start(port int) error {
	s.server.Handler = s.Handler
	log.Printf("start listening http-server at %d port\n", port)
	return s.server.ListenAndServe(fmt.Sprintf("localhost:%d", port))
}

func (s *Server) Stop() error {
	log.Println("stop http-server")
	return s.server.Shutdown()
}

func (s *Server) Handler(ctx *fasthttp.RequestCtx) {
	if string(ctx.Request.Header.Method()) == http.MethodGet && string(ctx.Request.URI().Path()) == "/subscribtions/all" {
		s.getSubscribtions(ctx)
		return
	} else if string(ctx.Request.Header.Method()) == http.MethodGet && string(ctx.Request.URI().Path()) == "/subscribtions" {
		s.getSubscribtion(ctx)
		return
	} else if string(ctx.Request.Header.Method()) == http.MethodGet && string(ctx.Request.URI().Path()) == "/subscribtions/price" {
		s.getSubscribtionsPrice(ctx)
		return
	} else if string(ctx.Request.Header.Method()) == http.MethodPost && string(ctx.Request.URI().Path()) == "/subscribtions" {
		s.createSubscribtion(ctx)
		return
	} else if string(ctx.Request.Header.Method()) == http.MethodPut && string(ctx.Request.URI().Path()) == "/subscribtions" {
		s.updateSubscribtion(ctx)
		return
	} else if string(ctx.Request.Header.Method()) == http.MethodDelete && string(ctx.Request.URI().Path()) == "/subscribtions" {
		s.deleteSubscribtion(ctx)
		return
	} else if string(ctx.Request.Header.Method()) == http.MethodGet && string(ctx.Request.URI().Path()) == "/swagger/index.html" ||
		string(ctx.Request.URI().Path()) == "/swagger/doc.json" ||
		string(ctx.Request.URI().Path()) == "/swagger/" {
		swagger.WrapHandler(swagger.InstanceName("swagger"))(ctx)
	}
	ctx.SetStatusCode(fasthttp.StatusNotFound)
}

// GetSubscribtions godoc
//
//	@Summary		Get all subscriptions
//	@Description	Получает все подписки
//	@Produce		json
//	@Success		200	{object}	[]domain.Subscribtion
//	@Failure		500	{object}	nil
//	@Router			/subscribtions/all [get]
func (s *Server) getSubscribtions(ctx *fasthttp.RequestCtx) {
	ctxPostgres, cancel := context.WithTimeout(context.Background(), s.cfg.Postgres.Timeout)
	defer cancel()

	subscribtions, err := s.repository.GetSubscribtions(ctxPostgres)
	if err != nil {
		log.Printf("failed to get subscribtions: %s\n", err.Error())
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	responseBody, err := json.Marshal(subscribtions)
	if err != nil {
		log.Printf("failed to marshal data: %s\n", err.Error())
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(responseBody)
}

// GetSubscribtion godoc
//
//	@Summary		Get subscription
//	@Description	Получает подписку
//	@Produce		json
//	@Success		200	{object}	domain.Subscribtion
//	@Failure		500	{object}	nil
//	@Failure		400	{object}	nil
//	@Router			/subscribtions [get]
func (s *Server) getSubscribtion(ctx *fasthttp.RequestCtx) {
	request, err := validateAndReturnStruct(ctx)
	if err != nil {
		return
	}

	ctxPostgres, cancel := context.WithTimeout(context.Background(), s.cfg.Postgres.Timeout)
	defer cancel()

	subscribtion, err := s.repository.GetSubscribtion(ctxPostgres, request.UserId, request.ServiceName)
	if err != nil {
		log.Printf("failed to get subscribtion: %s\n", err.Error())
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	responseBody, err := json.Marshal(subscribtion)
	if err != nil {
		log.Printf("failed to marshal data: %s\n", err.Error())
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(responseBody)
}

// GetSubscribtionsPrice godoc
//
//	@Summary		Get subscriptions price
//	@Description	Получает стоимость подписок за указанный период
//	@Produce		json
//	@Success		200	{object}	domain.PriceResponse
//	@Failure		500	{object}	nil
//	@Failure		400	{object}	nil
//	@Router			/subscribtions/price [get]
func (s *Server) getSubscribtionsPrice(ctx *fasthttp.RequestCtx) {
	request, err := validateAndReturnStruct(ctx)
	if err != nil {
		return
	}
	if request.End.IsValid() && request.Start.IsValid() && request.Start.After(*request.End) {
		log.Printf("invalid dates\n")
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	ctxPostgres, cancel := context.WithTimeout(context.Background(), s.cfg.Postgres.Timeout)
	defer cancel()
	price, err := s.repository.GetSubscribtionsPrice(ctxPostgres, request)
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

// CreateSubscribtion godoc
//
//		@Summary		Create subscription
//		@Description	Создает подписку
//		@Produce		json
//		@Success		201	{object}	nil
//	 	@Failure		500	{object}	nil
//	 	@Failure		400	{object}	nil
//		@Router			/subscribtions [get]
func (s *Server) createSubscribtion(ctx *fasthttp.RequestCtx) {
	var request domain.Subscribtion
	err := json.Unmarshal(ctx.Request.Body(), &request)
	if err != nil {
		log.Printf("failed to parse json: %s\n", err.Error())
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}
	if request.End.Month == 0 && request.End.Year == 0 {
		request.End = nil
	}
	if request.UserId == "" || request.ServiceName == "" {
		log.Printf("invalid json\n")
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	ctxPostgres, cancel := context.WithTimeout(context.Background(), s.cfg.Postgres.Timeout)
	defer cancel()

	err = s.repository.AddSubscribtion(ctxPostgres, request)
	if err != nil {
		log.Printf("failed to add subscribtion: %s\n", err.Error())
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusCreated)
}

// UpdateSubscribtion godoc
//
//	@Summary		Update subscription
//	@Description	Обновляет подписку
//	@Produce		json
//	@Success		200	{object}	nil
//	@Failure		500	{object}	nil
//	@Failure		400	{object}	nil
//	@Router			/subscribtions [put]
func (s *Server) updateSubscribtion(ctx *fasthttp.RequestCtx) {
	// TODO
	var request domain.Subscribtion
	err := json.Unmarshal(ctx.Request.Body(), &request)
	if err != nil {
		log.Printf("failed to parse json: %s\n", err.Error())
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}
	if request.End.Month == 0 && request.End.Year == 0 {
		request.End = nil
	}
	if request.UserId == "" || request.ServiceName == "" {
		log.Printf("invalid json\n")
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	if !request.Start.IsValid() || request.End != nil && !request.End.IsValid() {
		log.Printf("invalid date\n")
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	ctxPostgres, cancel := context.WithTimeout(context.Background(), s.cfg.Postgres.Timeout)
	defer cancel()

	err = s.repository.UpdateSubscribtion(ctxPostgres, request)
	if err != nil {
		log.Printf("failed to update subscribtion: %s\n", err.Error())
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
}

// DeleteSubscribtion godoc
//
//	@Summary		Delete subscription
//	@Description	Удаляет подписку
//	@Produce		json
//	@Success		200	{object}	nil
//	@Failure		500	{object}	nil
//	@Failure		400	{object}	nil
//	@Router			/subscribtions [delete]
func (s *Server) deleteSubscribtion(ctx *fasthttp.RequestCtx) {
	request, err := validateAndReturnStruct(ctx)
	if err != nil {
		return
	}

	ctxPostgres, cancel := context.WithTimeout(context.Background(), s.cfg.Postgres.Timeout)
	defer cancel()

	err = s.repository.DeleteSubscribtion(ctxPostgres, request.UserId, request.ServiceName)
	if err != nil {
		log.Printf("failed to delete subscribtion: %s\n", err.Error())
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
}

func validateAndReturnStruct(ctx *fasthttp.RequestCtx) (domain.Subscribtion, error) {
	var request domain.Subscribtion
	err := json.Unmarshal(ctx.Request.Body(), &request)
	if err != nil {
		log.Printf("failed to parse json: %s\n", err.Error())
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return domain.Subscribtion{}, err
	}
	if request.UserId == "" || request.ServiceName == "" {
		log.Printf("invalid json\n")
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return domain.Subscribtion{}, err
	}

	return request, nil
}
