package bus

import (
	"context"
	"log/slog"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
)

type Bus struct {
	router           *message.Router
	CommandBus       *cqrs.CommandBus
	commandProcessor *cqrs.CommandProcessor
	EventBus         *cqrs.EventBus
	EventProcessor   *cqrs.EventProcessor
}

func NewBus(l *slog.Logger) (*Bus, error) {
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
		Persistent: true,
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

	// Initialize bus operations if necessary
	return &Bus{
		router:           router,
		CommandBus:       commandBus,
		commandProcessor: commandProcessor,
		EventBus:         eventBus,
		EventProcessor:   eventProcessor,
	}, nil
}

func (b *Bus) Start(ctx context.Context) error {
	return b.router.Run(ctx)
}
