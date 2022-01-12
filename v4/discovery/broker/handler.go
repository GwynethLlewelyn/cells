package broker

import (
	"fmt"

	"github.com/pydio/cells/v4/common"
	"github.com/pydio/cells/v4/common/broker"
	pb "github.com/pydio/cells/v4/common/proto/broker"
)

type Handler struct {
	pb.UnimplementedBrokerServer

	broker      broker.Broker
	subscribers map[string][]string
}

func NewHandler(b broker.Broker) *Handler {
	return &Handler{
		broker: b,
	}
}

func (h *Handler) Name() string {
	return common.ServiceGrpcNamespace_ + common.ServiceBroker
}

func (h *Handler) Publish(stream pb.Broker_PublishServer) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			fmt.Println("Error is ", err)
			return err
		}

		for _, message := range req.Messages {
			if err := h.broker.PublishRaw(stream.Context(), req.Topic, message.Body, message.Header); err != nil {
				return err
			}
		}
	}
}

func (h *Handler) Subscribe(stream pb.Broker_SubscribeServer) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}

		topic := req.GetTopic()
		// queue := req.GetQueue()

		// TODO v4 - manage unsubscription
		h.broker.Subscribe(stream.Context(), topic, func(msg broker.Message) error {
			mm := &pb.SubscribeResponse{}
			var target = &pb.Message{}
			target.Header, target.Body = msg.RawData()

			// fmt.Println("And the message we've received is ? ", string(target.Body))
			// mm.Queue = queue
			mm.Id = req.Id
			mm.Messages = append(mm.Messages, target)
			if err := stream.Send(mm); err != nil {
				return err
			}

			return nil
		})
	}
}
