package producer

import (
	"fmt"

	"github.com/anrisys/quicket/user-service/internal/mq"
	"github.com/anrisys/quicket/user-service/pkg/mq/rabbitmq"
	"github.com/rs/zerolog"
)

type UserProduser struct {
	publisher *rabbitmq.Publisher
	logger zerolog.Logger
}

func NewUserPublisher(publisher *rabbitmq.Publisher, logger zerolog.Logger) *UserProduser {
	return &UserProduser{
		publisher: publisher,
		logger: logger,
	}
}

func (usp *UserProduser) PublishUserCreation(ID uint, PublicID string) error {
	log := usp.logger.With().
		Str("producer", "user_producer").
		Str("action", "publish_user_creation").
		Str("user_public_id", PublicID).
		Logger()
	
	exchangeName := "user.exchange"

	err := usp.publisher.DeclareExchange(exchangeName, "topic")
	if err != nil {
		log.Error().Err(err).Str("exchange", exchangeName).Msg("failed to declare exchange")
		return fmt.Errorf("%w: %v", mq.ErrFailedToDeclareExchange, err)
	}

	defer usp.publisher.Channel.Close()
	body := fmt.Appendf(nil, `{"ID": %d, "PublicID": %s}`, ID, PublicID)

	err = usp.publisher.Publish(exchangeName, "users.users.created", body)
	if err != nil {
		log.Error().Err(err).Str("exchange", exchangeName).Msg("failed to publish user creation message")
		return err
	}

	log.Info().Msgf("Published user creation: %s", string(body))
	return nil
}