package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	_ "subscriptions_rest/docs"
	"subscriptions_rest/internal/config"
	"subscriptions_rest/internal/domain"
	"subscriptions_rest/internal/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	fiberSwagger "github.com/swaggo/fiber-swagger"
)

type Server struct {
	repository repository.Repository
	cfg        *config.Config
	server     *fiber.App
}

type fiberError struct {
	Code    int
	Message string
	Fields  fiber.Map
}

func (e *fiberError) Error() string {
	return e.Message
}

func newError(err error, code int, fields fiber.Map) *fiberError {
	if fields == nil {
		fields = make(fiber.Map)
	}
	fields["error"] = err.Error()

	return &fiberError{
		Code:    code,
		Message: err.Error(),
		Fields:  fields,
	}
}

func New(cfg *config.Config, repository repository.Repository) *Server {
	server := fiber.New(fiber.Config{
		AppName: "searchbeat-coordinator",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			log.Println(err.Error())
			if handlerErr, ok := err.(*fiberError); ok {
				c.Status(handlerErr.Code).JSON(handlerErr.Fields) //nolint:errcheck
			}
			return err
		},
		Immutable: true,
	})

	server.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))
	return &Server{
		cfg:        cfg,
		repository: repository,
		server:     server,
	}
}

func (s *Server) Start(port int) error {
	s.registerV1GroupRoutes(s.server)
	log.Printf("start listening http-server at %d port\n", port)
	return s.server.Listen(fmt.Sprintf("localhost:%d", port))
}

func (s *Server) registerV1GroupRoutes(app *fiber.App) {
	app.Use("/swagger/*", fiberSwagger.WrapHandler)
	apiGroup := app.Group("/api")
	v1Group := apiGroup.Group("/v1")

	v1Group.Get("/subscribtions/all", s.getSubscribtions)
	v1Group.Get("/subscribtions", s.getSubscribtion)
	v1Group.Get("/subscribtions/price", s.getSubscribtionsPrice)
	v1Group.Post("/subscribtions", s.createSubscribtion)
	v1Group.Put("/subscribtions", s.updateSubscribtion)
	v1Group.Delete("/subscribtions", s.deleteSubscribtion)
}

func (s *Server) Stop() error {
	log.Println("stop http-server")
	return s.server.Shutdown()
}

// func (s *Server) Handler(ctx *fasthttp.RequestCtx) {
// 	if string(ctx.Request.Header.Method()) == http.MethodGet && string(ctx.Request.URI().Path()) == "/api/v1/subscribtions/all" {
// 		s.getSubscribtions(ctx)
// 		return
// 	} else if string(ctx.Request.Header.Method()) == http.MethodGet && string(ctx.Request.URI().Path()) == "/api/v1/subscribtions" {
// 		s.getSubscribtion(ctx)
// 		return
// 	} else if string(ctx.Request.Header.Method()) == http.MethodGet && string(ctx.Request.URI().Path()) == "/api/v1/subscribtions/price" {
// 		s.getSubscribtionsPrice(ctx)
// 		return
// 	} else if string(ctx.Request.Header.Method()) == http.MethodPost && string(ctx.Request.URI().Path()) == "/api/v1/subscribtions" {
// 		s.createSubscribtion(ctx)
// 		return
// 	} else if string(ctx.Request.Header.Method()) == http.MethodPut && string(ctx.Request.URI().Path()) == "/api/v1/subscribtions" {
// 		s.updateSubscribtion(ctx)
// 		return
// 	} else if string(ctx.Request.Header.Method()) == http.MethodDelete && string(ctx.Request.URI().Path()) == "/api/v1/subscribtions" {
// 		s.deleteSubscribtion(ctx)
// 		return
// 	} else if string(ctx.Request.Header.Method()) == http.MethodGet &&
// 		strings.HasPrefix(string(ctx.Request.URI().Path()), "/swagger") {
// 		swagger.WrapHandler(swagger.InstanceName("swagger"))(ctx)
// 	}
// 	ctx.SetStatusCode(fasthttp.StatusNotFound)
// }

// GetSubscribtions godoc
//
//	@Summary		Get all subscriptions
//	@Description	Получает все подписки
//	@Produce		json
//	@Success		200	{object}	[]domain.Subscribtion
//	@Failure		500	{object}	nil
//	@Router			/subscribtions/all [get]
func (s *Server) getSubscribtions(c *fiber.Ctx) error {
	ctxPostgres, cancel := context.WithTimeout(context.Background(), s.cfg.Postgres.Timeout)
	defer cancel()

	subscribtions, err := s.repository.GetSubscribtions(ctxPostgres)
	if err != nil {
		return newError(fmt.Errorf("failed to get subscribtions: %w", err),
			http.StatusInternalServerError, fiber.Map{})
	}

	return c.Status(http.StatusOK).JSON(subscribtions)
}

// GetSubscribtion godoc
//
//	@Summary		Get subscription
//	@Description	Получает подписку
//	@Param 			user_id query string true "get subscription by user id"
//	@Param 			service_name query string true "get subscription by service name"
//	@Produce		json
//	@Success		200	{object}	domain.Subscribtion
//	@Failure		500	{object}	nil
//	@Failure		400	{object}	nil
//	@Router			/subscribtions [get]
func (s *Server) getSubscribtion(c *fiber.Ctx) error {
	queries := c.Queries()
	userId, okUserId := queries["user_id"]
	serviceName, okServiceName := queries["service_name"]
	if !okUserId || !okServiceName {
		return newError(fmt.Errorf("invalid params"),
			http.StatusBadRequest, fiber.Map{})
	}

	ctxPostgres, cancel := context.WithTimeout(context.Background(), s.cfg.Postgres.Timeout)
	defer cancel()

	subscribtion, err := s.repository.GetSubscribtion(ctxPostgres, userId, serviceName)
	if err != nil {
		return newError(fmt.Errorf("failed to get subscribtion: %w", err),
			http.StatusInternalServerError, fiber.Map{})
	}

	return c.Status(http.StatusOK).JSON(subscribtion)
}

// GetSubscribtionsPrice godoc
//
//	@Summary		Get subscriptions price
//	@Description	Получает стоимость подписок за указанный период
//	@Param 			user_id query string true "get subscription price by user id"
//	@Param 			service_name query string true "get subscription price by service name"
//	@Param 			start_date query string false "filter by start date"
//	@Param 			end_date query string false "filter by end date"
//	@Produce		json
//	@Success		200	{object}	domain.PriceResponse
//	@Failure		500	{object}	nil
//	@Failure		400	{object}	nil
//	@Router			/subscribtions/price [get]
func (s *Server) getSubscribtionsPrice(c *fiber.Ctx) error {
	queries := c.Queries()
	userId, okUserId := queries["user_id"]
	serviceName, okServiceName := queries["service_name"]
	if !okUserId || !okServiceName {
		return newError(fmt.Errorf("invalid params"),
			http.StatusBadRequest, fiber.Map{})
	}
	startDateStr, okStart := queries["start_date"]
	endDateStr, okEnd := queries["end_date"]
	var startDate, endDate *domain.Date
	var err error
	if okStart {
		startDate, err = domain.DateFromString(startDateStr)
		if err != nil {
			return newError(fmt.Errorf("invalid start date"),
				http.StatusBadRequest, fiber.Map{})
		}
	}
	if okEnd {
		endDate, err = domain.DateFromString(endDateStr)
		if err != nil {
			return newError(fmt.Errorf("invalid end date"),
				http.StatusBadRequest, fiber.Map{})
		}
	}

	if endDate != nil && endDate.IsValid() && startDate != nil && startDate.IsValid() && startDate.After(*endDate) {
		return newError(fmt.Errorf("invalid dates"),
			http.StatusBadRequest, fiber.Map{})
	}

	ctxPostgres, cancel := context.WithTimeout(context.Background(), s.cfg.Postgres.Timeout)
	defer cancel()
	price, err := s.repository.GetSubscribtionsPrice(ctxPostgres, domain.Subscribtion{
		UserId:      userId,
		ServiceName: serviceName,
		Start:       startDate,
		End:         endDate,
	})
	if err != nil {
		return newError(fmt.Errorf("failed to get price: %w", err),
			http.StatusInternalServerError, fiber.Map{})
	}

	return c.Status(http.StatusOK).JSON(domain.PriceResponse{
		Price: price,
	})
}

// CreateSubscribtion godoc
//
//		@Summary		Create subscription
//		@Description	Создает подписку
//		@Param 			Entry body domain.SubscribtionRequest true "createSubscribtion" example({"service_name": "Yandex Plus", "user_id": "", "price": 400, "start_date": "07-2025", "end_date": ""} )
//		@Produce		json
//		@Success		201	{object}	nil
//	 	@Failure		500	{object}	nil
//	 	@Failure		400	{object}	nil
//		@Router			/subscribtions [post]
func (s *Server) createSubscribtion(c *fiber.Ctx) error {
	request, err := validateAndReturnStruct(c)
	if err != nil {
		return err
	}
	if request.Start == nil || request.End != nil && request.End.IsValid() && request.Start != nil && request.Start.IsValid() && request.Start.After(*request.End) {
		return newError(fmt.Errorf("invalid dates"),
			http.StatusBadRequest, fiber.Map{})
	}
	if request.End.Month == 0 && request.End.Year == 0 {
		request.End = nil
	}

	ctxPostgres, cancel := context.WithTimeout(context.Background(), s.cfg.Postgres.Timeout)
	defer cancel()

	err = s.repository.AddSubscribtion(ctxPostgres, request)
	if err != nil {
		return newError(fmt.Errorf("failed to add subscribtion: %w", err),
			http.StatusInternalServerError, fiber.Map{})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{})
}

// UpdateSubscribtion godoc
//
//	@Summary		Update subscription
//	@Description	Обновляет подписку
//	@Param 			Entry body domain.SubscribtionRequest true "updateSubscribtion" example({"service_name": "Yandex Plus", "user_id": "", "price": 400, "start_date": "07-2025", "end_date": ""} )
//	@Produce		json
//	@Success		200	{object}	nil
//	@Failure		500	{object}	nil
//	@Failure		400	{object}	nil
//	@Router			/subscribtions [put]
func (s *Server) updateSubscribtion(c *fiber.Ctx) error {
	request, err := validateAndReturnStruct(c)
	if err != nil {
		return err
	}
	if request.End != nil && request.End.Month == 0 && request.End.Year == 0 {
		request.End = nil
	}
	if request.End != nil && !request.End.IsValid() || request.Start == nil || request.End != nil && !request.Start.IsValid() || request.End != nil && request.End.IsValid() && request.Start.IsValid() && request.Start.After(*request.End) {
		return newError(fmt.Errorf("invalid dates"),
			http.StatusBadRequest, fiber.Map{})
	}

	ctxPostgres, cancel := context.WithTimeout(context.Background(), s.cfg.Postgres.Timeout)
	defer cancel()

	err = s.repository.UpdateSubscribtion(ctxPostgres, request)
	if err != nil {
		return newError(fmt.Errorf("failed to update subscribtion: %w", err),
			http.StatusInternalServerError, fiber.Map{})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{})
}

// DeleteSubscribtion godoc
//
//	@Summary		Delete subscription
//	@Description	Удаляет подписку
//	@Param 			user_id query string true "delete subscription by user id"
//	@Param 			service_name query string true "delete subscription by service name"
//	@Produce		json
//	@Success		200	{object}	nil
//	@Failure		500	{object}	nil
//	@Failure		400	{object}	nil
//	@Router			/subscribtions [delete]
func (s *Server) deleteSubscribtion(c *fiber.Ctx) error {
	queries := c.Queries()
	userId, okUserId := queries["user_id"]
	serviceName, okServiceName := queries["service_name"]
	if !okUserId || !okServiceName {
		return newError(fmt.Errorf("invalid params"),
			http.StatusBadRequest, fiber.Map{})
	}

	ctxPostgres, cancel := context.WithTimeout(context.Background(), s.cfg.Postgres.Timeout)
	defer cancel()

	err := s.repository.DeleteSubscribtion(ctxPostgres, userId, serviceName)
	if err != nil {
		return newError(fmt.Errorf("failed to delete subscribtion: %w", err),
			http.StatusInternalServerError, fiber.Map{})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{})
}

func validateAndReturnStruct(c *fiber.Ctx) (domain.Subscribtion, error) {
	var request domain.SubscribtionRequest
	if err := c.BodyParser(&request); err != nil {
		return domain.Subscribtion{}, newError(fmt.Errorf("failed to parse request body: %w", err),
			http.StatusBadRequest, fiber.Map{})
	}
	if request.UserId == "" || request.ServiceName == "" {
		return domain.Subscribtion{}, newError(fmt.Errorf("invalid json"),
			http.StatusBadRequest, fiber.Map{})
	}

	var start, end *domain.Date
	var err error
	if request.Start != "" {
		start, err = domain.DateFromString(request.Start)
		if err != nil {
			return domain.Subscribtion{}, newError(fmt.Errorf("invalid start date"),
				http.StatusBadRequest, fiber.Map{})
		}
	}
	if request.End != "" {
		end, err = domain.DateFromString(request.End)
		if err != nil {
			return domain.Subscribtion{}, newError(fmt.Errorf("invalid end date"),
				http.StatusBadRequest, fiber.Map{})
		}
	}

	return domain.Subscribtion{
		UserId:      request.UserId,
		ServiceName: request.ServiceName,
		Price:       request.Price,
		Start:       start,
		End:         end,
	}, nil
}
