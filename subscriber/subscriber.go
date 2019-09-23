package subscriber

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/dairaga/log"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/sawtooth-sdk-go/messaging"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/client_event_pb2"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/events_pb2"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/validator_pb2"

	zmq "github.com/pebbe/zmq4"
)

// EvtTypBlockID sawtooth block id event
//const EvtTypBlockID = "sawtooth/block-commit"

// BlockIDFile save block id
// const BlockIDFile = ".sub_block_id"

// Handler is a sawtooth subscriber handler.
// Return true if need to pass next handler, or return false.
type Handler func(string, *events_pb2.Event) bool

// Subscriber is a sawtooth event subscriber.
type Subscriber struct {
	endpoint string                          // validator endpoint.
	conn     *messaging.ZmqConnection        // zero mq connection.
	signal   chan os.Signal                  // system kill or interrupt signal.
	callback chan bool                       // wait for unsubscribing.
	running  bool                            // running flag.
	handlers map[string][]Handler            // event handlers.
	events   []*events_pb2.EventSubscription // subscribing events.
	wait     time.Duration                   // duration of waiting for unsubscribing event.

	RecordBlockID bool // flag to record last block id into file.

	OnClose        func()
	OnSubcribed    func(string, *client_event_pb2.ClientEventsSubscribeResponse)
	OnUnsubscribed func(string, *client_event_pb2.ClientEventsUnsubscribeResponse)
}

// New returns a subcriber for some sawtooth events.
func New(endpoint string, wait time.Duration) (*Subscriber, error) {
	ctx, err := zmq.NewContext()
	if err != nil {
		return nil, err
	}

	sub := new(Subscriber)

	sub.wait = wait
	sub.endpoint = endpoint

	sub.conn, err = messaging.NewConnection(ctx, zmq.DEALER, endpoint, false)
	if err != nil {
		return nil, err
	}

	sub.callback = make(chan bool)
	sub.signal = make(chan os.Signal)
	signal.Notify(sub.signal, os.Interrupt, os.Kill)

	sub.handlers = make(map[string][]Handler)

	return sub, nil
}

// routing is to call handler(s) for some event.
func (s *Subscriber) routing(id string, msg *validator_pb2.Message) {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("panic: %s %s %v", msg.CorrelationId, msg.MessageType.String(), r)
		}
	}()

	switch msg.MessageType {
	case validator_pb2.Message_CLIENT_EVENTS_SUBSCRIBE_RESPONSE: // subscribe response

		if s.OnSubcribed != nil {
			resp := new(client_event_pb2.ClientEventsSubscribeResponse)
			if err := proto.Unmarshal(msg.Content, resp); err == nil {
				s.OnSubcribed(id, resp)
			} else {
				log.Warn("client events subscribe response ummarshal: ", err)
			}
		}

	case validator_pb2.Message_CLIENT_EVENTS: // events

		evts := new(events_pb2.EventList)

		if err := proto.Unmarshal(msg.Content, evts); err == nil {
			for _, evt := range evts.Events {

				handlers, ok := s.handlers[evt.EventType]

				if ok && handlers != nil {
					for i, h := range handlers {
						if !h(id, evt) {
							log.Errorf("%s handler (%d) stop", evt.EventType, i)
							break
						}
					}
				} else {
					log.Warn("no handler for ", evt.EventType)
				}
			}
		} else {
			log.Warn("event list unmarshal: ", err)
		}

	case validator_pb2.Message_CLIENT_EVENTS_UNSUBSCRIBE_RESPONSE: // unsubscribe

		if s.OnUnsubscribed != nil {
			resp := new(client_event_pb2.ClientEventsUnsubscribeResponse)
			if err := proto.Unmarshal(msg.Content, resp); err == nil {
				s.OnUnsubscribed(id, resp)
			} else {
				log.Warn("client event unsubscribe response unmarshal: ", err)
			}
		}
		s.callback <- true

	default:
		log.Warn("unknown message type: ", msg.MessageType)
	}
}

// unsubscribe unsubscribe from chain.
func (s *Subscriber) unsubscribe() {
	unsub := new(client_event_pb2.ClientEventsUnsubscribeRequest)
	unsubBytes, err := proto.Marshal(unsub)
	if err != nil {
		log.Warn("client events unsubscribe request unmarshal: ", err)
		return
	}

	if _, err := s.conn.SendNewMsg(validator_pb2.Message_CLIENT_EVENTS_UNSUBSCRIBE_REQUEST, unsubBytes); err != nil {
		log.Error("send unsubscribe request: ", err)
	}
}

// HandleFunc appends an handler for some event.
func (s *Subscriber) HandleFunc(eventType string, h Handler) {
	s.handlers[eventType] = append(s.handlers[eventType], h)
}

// Subscribe records an event subscription with filters and handler for it.
func (s *Subscriber) Subscribe(eventType string, h Handler, filters ...*events_pb2.EventFilter) {
	s.events = append(s.events, &events_pb2.EventSubscription{
		EventType: eventType,
		Filters:   filters,
	})
	s.handlers[eventType] = append(s.handlers[eventType], h)
}

// Run subscriber run
func (s *Subscriber) run() {

	s.running = true

	go func() {
		log.Info("run start")
		for {
			id, msg, err := s.conn.RecvMsg()

			if err != nil {
				log.Error("recv: ", err)
				break
			}
			s.routing(id, msg)
		}
		log.Info("run end")
	}()
}

/*
// readLastBlockID returns last block id in file if exists.
func (s *Subscriber) readLastBlockID() (string, error) {
	dataBytes, err := ioutil.ReadFile(BlockIDFile)
	if err != nil {
		return "", err
	}

	var data []*events_pb2.Event_Attribute

	if err := json.Unmarshal(dataBytes, &data); err != nil {
		return "", err
	}

	for _, x := range data {
		if x.Key == "previous_block_id" {
			return x.Value, nil
		}
	}

	return "", fmt.Errorf("previous block id not found")
}
*/

// sendSubscribe sends subscription after running.
func (s *Subscriber) sendSubscribe(lastBlockIDs []string) (string, error) {

	/*if s.RecordBlockID {
		s.Subscribe(EvtTypBlockID, s.handleBlock)
		if prevID, err := s.readLastBlockID(); err != nil {
			log.Warnf("read last block id: %v", err)
		} else {
			lastBlockIDs = append(lastBlockIDs, prevID)
		}
	}*/

	req := &client_event_pb2.ClientEventsSubscribeRequest{
		LastKnownBlockIds: lastBlockIDs, // adds last block id
		Subscriptions:     s.events,     // events to subscribe
	}

	reqBytes, err := proto.Marshal(req)
	if err != nil {
		log.Error("ClientEventsSubscribeRequest marshal:", err)
		return "", err
	}

	corID, err := s.conn.SendNewMsg(validator_pb2.Message_CLIENT_EVENTS_SUBSCRIBE_REQUEST, reqBytes)
	if err != nil {
		log.Error("send subscribe request:", err)
		return "", err
	}
	log.Debug("send sub: ", req)
	return corID, nil
}

/*
// handleBlock to record block data into file.
func (s *Subscriber) handleBlock(id string, evt *events_pb2.Event) bool {

	dataBytes, err := json.Marshal(evt.Attributes)
	if err != nil {
		log.Errorf("json unmarshal: %v", err)
		return true // force to user's handler
	}
	//log.Debugf("record block id: %s", string(dataBytes))

	if err := ioutil.WriteFile(BlockIDFile, dataBytes, 0644); err != nil {
		log.Warnf("write block id: %v", err)
	}

	for i, attr := range evt.Attributes {
		if attr.Key == "block_id" {
			log.Debugf("%02d: key: %s, value: %s", i, attr.Key, attr.Value)
			break
		}
	}

	return true
}
*/

// WaitForShutdown is to start a subscriber and wait for shutdown.
func (s *Subscriber) WaitForShutdown(lastBlockIDs ...string) {
	defer func() {
		if r := recover(); r != nil {
			log.Error("run panic: ", r)
		}

		if s.OnClose != nil {
			s.OnClose()
		}

		close(s.signal)
		close(s.callback)
		s.conn.Close()
		log.Info("end")
	}()

	s.run()

	id, err := s.sendSubscribe(lastBlockIDs)
	if err != nil {
		log.Error("send subscribe: ", err)
		return
	}

	log.Debug("send id:", id)

	<-s.signal
	log.Info("got interrupted")

	s.running = false
	s.unsubscribe()

	ctxTimeout, cancelTimeOut := context.WithTimeout(context.Background(), s.wait)
	defer cancelTimeOut()

	select {
	case <-s.callback:
		log.Info("get unsubscribe and end")
	case <-ctxTimeout.Done():
		log.Warn("timeout and end")
	}
}
