package proc

import (
	"context"

	"github.com/wahed-tech/banking-ops-svc/conf"
	"github.com/wahed-tech/go-lib/v2/constant"
	"github.com/wahed-tech/go-lib/v2/errors"
	"github.com/wahed-tech/go-lib/v2/skafka"
	"github.com/wahed-tech/go-lib/v2/uuid"
	"github.com/wahed-tech/protobuf-definitions/protobuf"
	"github.com/wahed-tech/protobuf-definitions/protobuf/example"
)

type CreateExample struct {
	producer skafka.Producer
}

func NewCreateExample(producer skafka.Producer) *CreateExample {
	return &CreateExample{producer: producer}
}

func (c *CreateExample) CanHandle(headers map[string]string) bool {
	return example.EventCreateExample == headers[constant.EventTypeHeader]
}

func (c *CreateExample) Handle(ctx context.Context, message *protobuf.Message, headers map[string]string) error {
	payload := protobuf.GetExamplePayload(message).(*example.ExamplePayload_CreateExample).CreateExample
	id := upsertExample(payload)

	headers[constant.EventTypeHeader] = example.EventExampleCreated
	message.Payload = protobuf.CreateExampleEvent(
		&example.ExamplePayload_ExampleCreated{ExampleCreated: &example.ExampleCreated{
			Id: id,
		}})

	if err := c.producer.SendEvent(conf.R.Kafka.Topics.ExampleEvents, message, headers); err != nil {
		return errors.Wrap(err, errors.ErrKafkaSend)
	}

	return nil
}

func upsertExample(_ *example.CreateExample) string {
	return uuid.NewString()
}
