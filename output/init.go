package input

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fluent/fluent-logger-golang/fluent"
	"github.com/redpanda-data/benthos/v4/public/service"
)

const (
	outputHost        = "host"
	outputPort        = "port"
	outputTag         = "tag"
	outputMaxInFlight = "max_in_flight"
)

var (
	_ service.Output = &ForwardOutput{}
)

type ForwardOutput struct {
	logger *fluent.Fluent
	host   string
	port   int
	tag    string
}

func NewForwardOutput(host string, port int, tag string) *ForwardOutput {
	return &ForwardOutput{
		host: host,
		port: port,
		tag:  tag,
	}
}

func (o *ForwardOutput) Connect(_ context.Context) error {
	logger, err := fluent.New(fluent.Config{
		FluentHost: o.host,
		FluentPort: o.port,
	})
	if err != nil {
		return fmt.Errorf("could not create the forward client: %s", err.Error())
	}
	o.logger = logger

	return nil
}

func (o *ForwardOutput) Write(_ context.Context, msg *service.Message) error {
	byteMsg, err := msg.AsBytes()
	var data map[string]interface{}

	err = json.Unmarshal(byteMsg, &data)
	if err != nil {
		return err
	}
	err = o.logger.Post(o.tag, data)

	if err != nil {
		return err
	}
	return nil
}

func (o *ForwardOutput) Close(_ context.Context) error {
	err := o.logger.Close()
	if err != nil {
		return err
	}
	return nil
}

func init() {
	config := service.NewConfigSpec().
		Summary("Fluentd Forward output plugin").
		Fields(
			service.NewStringField(outputHost).
				Description("Fluentd server name").
				Default("localhost"),
			service.NewIntField(outputPort).
				Description("Fluentd server port").
				Default(24224),
			service.NewStringField(outputTag).
				Description("Fluentd tag in dot separated format"),
			service.NewIntField(outputMaxInFlight).
				Description("Max number of write calls that can be run in parallel").
				Default(10),
		)

	err := service.RegisterOutput("fluentd_forward", config, func(conf *service.ParsedConfig, mgr *service.Resources) (out service.Output, maxInFlight int, err error) {
		host, err := conf.FieldString(outputHost)
		if err != nil {
			panic(err)
		}
		port, err := conf.FieldInt(outputPort)
		if err != nil {
			panic(err)
		}
		tag, err := conf.FieldString(outputTag)
		if err != nil {
			panic(err)
		}
		maxInFlight, err = conf.FieldInt(outputMaxInFlight)
		if err != nil {
			panic(err)
		}

		forwardOutput := NewForwardOutput(host, port, tag)
		return forwardOutput, maxInFlight, nil
	})
	if err != nil {
		fmt.Printf("failed to register plugin: %s\n", err.Error())
		return
	}
}
