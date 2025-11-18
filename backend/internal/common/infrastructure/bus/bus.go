package bus

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/SirNacou/weeate/backend/internal/common/domain"
	"github.com/ThreeDotsLabs/watermill"
	wsql "github.com/ThreeDotsLabs/watermill-sql/v4/pkg/sql"
	"github.com/samber/lo"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/components/forwarder"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
)

type Bus struct {
	router           *message.Router
	logger           watermill.LoggerAdapter
	CommandBus       *cqrs.CommandBus
	commandProcessor *cqrs.CommandProcessor
	EventBus         *cqrs.EventBus
	EventProcessor   *cqrs.EventProcessor
}

func NewBus(conn *sql.Conn, l *slog.Logger) (*Bus, error) {
	logger := watermill.NewSlogLoggerWithLevelMapping(l, map[slog.Level]slog.Level{
		slog.LevelDebug: slog.LevelDebug,
		slog.LevelInfo:  slog.LevelInfo,
		slog.LevelWarn:  slog.LevelWarn,
		slog.LevelError: slog.LevelError,
	})

	cqrsMarshaler := cqrs.ProtoMarshaler{
		GenerateName: cqrs.StructName,
	}

	generateEventsTopic := func(eventName string) string {
		return "events." + eventName
	}

	generateCommandsTopic := func(commandName string) string {
		return "commands." + commandName
	}

	goChannel := gochannel.NewGoChannel(gochannel.Config{
		Persistent:                     true,
		PreserveContext:                true,
		BlockPublishUntilSubscriberAck: true,
	}, logger)

	// CQRS is built on messages router. Detailed documentation: https://watermill.io/docs/messages-router/
	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		return nil, err
	}

	router.AddMiddleware(middleware.Recoverer)

	commandBus, err := cqrs.NewCommandBusWithConfig(goChannel, cqrs.CommandBusConfig{
		GeneratePublishTopic: func(cbgptp cqrs.CommandBusGeneratePublishTopicParams) (string, error) {
			return generateCommandsTopic(cbgptp.CommandName), nil
		},
		OnSend: func(params cqrs.CommandBusOnSendParams) error {
			logger.Info("Sending command", watermill.LogFields{
				"command_name": params.CommandName,
			})

			params.Message.Metadata.Set("sent_at", time.Now().String())

			return nil
		},
		Marshaler: cqrsMarshaler,
		Logger:    logger,
	})
	if err != nil {
		return nil, err
	}

	commandProcessor, err := cqrs.NewCommandProcessorWithConfig(
		router,
		cqrs.CommandProcessorConfig{
			GenerateSubscribeTopic: func(params cqrs.CommandProcessorGenerateSubscribeTopicParams) (string, error) {
				return generateCommandsTopic(params.CommandName), nil
			},
			SubscriberConstructor: func(params cqrs.CommandProcessorSubscriberConstructorParams) (message.Subscriber, error) {
				// we can reuse subscriber, because all commands have separated topics
				return goChannel, nil
			},
			OnHandle: func(params cqrs.CommandProcessorOnHandleParams) error {
				start := time.Now()

				err := params.Handler.Handle(params.Message.Context(), params.Command)

				logger.Info("Command handled", watermill.LogFields{
					"command_name": params.CommandName,
					"duration":     time.Since(start),
					"err":          err,
				})

				return err
			},
			Marshaler: cqrsMarshaler,
			Logger:    logger,
		},
	)
	if err != nil {
		return nil, err
	}

	eventBus, err := cqrs.NewEventBusWithConfig(goChannel, cqrs.EventBusConfig{
		GeneratePublishTopic: func(params cqrs.GenerateEventPublishTopicParams) (string, error) {
			return generateEventsTopic(params.EventName), nil
		},

		OnPublish: func(params cqrs.OnEventSendParams) error {
			logger.Info("Publishing event", watermill.LogFields{
				"event_name": params.EventName,
			})

			params.Message.Metadata.Set("published_at", time.Now().String())

			return nil
		},

		Marshaler: cqrsMarshaler,
		Logger:    logger,
	})
	if err != nil {
		return nil, err
	}

	eventProcessor, err := cqrs.NewEventProcessorWithConfig(
		router,
		cqrs.EventProcessorConfig{
			GenerateSubscribeTopic: func(params cqrs.EventProcessorGenerateSubscribeTopicParams) (string, error) {
				return generateEventsTopic(params.EventName), nil
			},
			SubscriberConstructor: func(params cqrs.EventProcessorSubscriberConstructorParams) (message.Subscriber, error) {
				return goChannel, nil
			},

			OnHandle: func(params cqrs.EventProcessorOnHandleParams) error {
				start := time.Now()

				err := params.Handler.Handle(params.Message.Context(), params.Event)

				logger.Info("Event handled", watermill.LogFields{
					"event_name": params.EventName,
					"duration":   time.Since(start),
					"err":        err,
				})

				return err
			},

			Marshaler: cqrsMarshaler,
			Logger:    logger,
		},
	)
	if err != nil {
		return nil, err
	}

	forwarderConfig := wsql.SubscriberConfig{
		InitializeSchema: true,
		SchemaAdapter:    wsql.DefaultPostgreSQLSchema{},
		OffsetsAdapter:   wsql.DefaultPostgreSQLOffsetsAdapter{},
	}
	sqlSubcriber, err := wsql.NewSubscriber(wsql.BeginnerFromStdSQL(conn), forwarderConfig, logger)
	if err != nil {
		return nil, err
	}

	_, err = forwarder.NewForwarder(sqlSubcriber, goChannel, logger, forwarder.Config{
		Router: router,
	})
	if err != nil {
		return nil, err
	}

	// Initialize bus operations if necessary
	return &Bus{
		router:           router,
		logger:           logger,
		CommandBus:       commandBus,
		commandProcessor: commandProcessor,
		EventBus:         eventBus,
		EventProcessor:   eventProcessor,
	}, nil
}

func (b *Bus) Start(ctx context.Context) error {
	return b.router.Run(ctx)
}

type SqlPublisher struct {
	pub    message.Publisher
	logger watermill.LoggerAdapter
}

func (b *Bus) NewSqlPublisher(tx *sql.Tx) (*SqlPublisher, error) {
	var publisher message.Publisher
	publisher, err := wsql.NewPublisher(
		wsql.TxFromStdSQL(tx),
		wsql.PublisherConfig{
			SchemaAdapter: wsql.DefaultPostgreSQLSchema{},
		},
		b.logger,
	)
	if err != nil {
		return nil, err
	}

	return &SqlPublisher{
		pub:    publisher,
		logger: b.logger,
	}, nil
}

func (p *SqlPublisher) Publish(ctx context.Context, events ...domain.Event) error {
	// Decorate publisher so it wraps an event in an envelope understood by the Forwarder component.
	publisher := forwarder.NewPublisher(p.pub, forwarder.PublisherConfig{
		ForwarderTopic: "eventsToForward",
	})

	// Publish an event announcing the lottery winner. Please note we're publishing to a Google Cloud topic here,
	// while using decorated MySQL publisher.
	groupedEvents := lo.GroupBy(events, func(e domain.Event) string { return e.Name() })
	for eventName, eventList := range groupedEvents {
		messages := []*message.Message{}
		for _, event := range eventList {
			payload, err := json.Marshal(event)
			if err != nil {
				return err
			}
			messages = append(messages, message.NewMessageWithContext(ctx, watermill.NewULID(), payload))
		}
		err := publisher.Publish("sql_events."+eventName,
			messages...,
		)
		if err != nil {
			return err
		}
	}
	return nil
}
