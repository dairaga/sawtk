package tp

import (
	"fmt"
	"os"
	"reflect"

	"github.com/dairaga/log"
	"github.com/golang/protobuf/proto"

	"github.com/hyperledger/sawtooth-sdk-go/processor"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/processor_pb2"
)

// -----------------------------------------------------------------------------------

// unmarshal bytes to protobuf message.
func unmarshal(buf []byte, pb proto.Message) *processor.InvalidTransactionError {
	if err := proto.Unmarshal(buf, pb); err != nil {
		return Unmarshal.TxErrore(err)
	}

	return nil
}

// -----------------------------------------------------------------------------------

// HandlerFunc SawTK TP handler function.
type HandlerFunc func(ctx *Context, req *TPRequest) *processor.InvalidTransactionError

// RequestValidator is for MakeHandlerFunc to check req is ok or not.
type RequestValidator interface {
	proto.Message
	Validate() *processor.InvalidTransactionError
}

var (
	_ctxw = reflect.TypeOf((*Context)(nil))
	_req  = reflect.TypeOf((*processor_pb2.TpProcessRequest)(nil))
	_err  = reflect.TypeOf((*processor.InvalidTransactionError)(nil))
	_rv   = reflect.TypeOf((*RequestValidator)(nil)).Elem()
	_hf   = reflect.TypeOf((HandlerFunc)(nil))
)

// MakeHandlerFunc returns a HandlerFunc.
// It will unmarshal payload to specific datatype defined in f automatically.
func MakeHandlerFunc(f interface{}) HandlerFunc {
	ftype := reflect.TypeOf(f)
	if ftype.Kind() != reflect.Func {
		panic("f is not a function")
	}

	// check func must be f(*Context, ...) processor.InvalidTransactionError. and data must be a ptr, proto.Message, and RequestValidator.
	if ftype.NumOut() != 1 || ftype.NumIn() < 1 ||
		!ftype.Out(0).ConvertibleTo(_err) ||
		!ftype.In(0).ConvertibleTo(_ctxw) {

		panic(fmt.Sprintf("bad func def: %T", f))
	}

	if ftype.NumIn() > 1 {
		if ftype.NumIn() != 2 {
			panic(fmt.Sprintf("bad func def: %T", f))
		}

		if !ftype.In(1).Implements(_rv) {
			panic(fmt.Sprintf("bad func def: %T", f))
		}
	}

	x := reflect.MakeFunc(_hf, func(args []reflect.Value) []reflect.Value {
		input := make([]reflect.Value, 0, ftype.NumIn())
		input = append(input, args[0])

		if ftype.NumIn() == 2 {
			req := args[1].Interface().(*TPRequest)
			if req.Payload == nil {
				return []reflect.Value{reflect.ValueOf(BadParameters.TxErrorf("payload is nil"))}
			}

			newReq := reflect.New(ftype.In(1).Elem())

			if err := unmarshal(req.Payload, newReq.Interface().(proto.Message)); err != nil {
				return []reflect.Value{reflect.ValueOf(err)}
			}

			if err := newReq.Interface().(RequestValidator).Validate(); err != nil {
				return []reflect.Value{reflect.ValueOf(err)}
			}

			input = append(input, newReq)
		}

		fvalue := reflect.ValueOf(f)
		return fvalue.Call(input)
	})

	return x.Interface().(HandlerFunc)

}

// Handler is a handler for Sawtooth Transaction Processor.
type Handler struct {
	*Family
	router map[int32]HandlerFunc
	debug  bool
}

// Apply implements Apply function of processor.TransactionHandler.
// Decode payload in req (*processor_pb2.TpProcessRequest), and change to user's payload.
func (h *Handler) Apply(req *processor_pb2.TpProcessRequest, ctx *processor.Context) error {
	r := new(TPRequest)

	// 加一個環境變數，來設定是否要進入 debug mode,
	// 如果是 debug mode 則取消 recover.
	if !h.debug {
		defer func(r *TPRequest) {
			if r := recover(); r != nil {
				log.Fatal(r)
			}
		}(r)
	}

	if err := proto.Unmarshal(req.Payload, r); err != nil {
		return Unmarshal.TxErrore(err)
	}

	log.Debugf("got CMD (%d)", r.Cmd)
	hfunc, ok := h.router[r.Cmd]
	if !ok || hfunc == nil {
		return &processor.InvalidTransactionError{
			Msg: fmt.Sprintf("unknow cmd: %d", r.Cmd),
		}
	}

	ctxw := &Context{
		ref:    ctx,
		signer: req.Header.SignerPublicKey,
		cmd:    r.Cmd,
	}

	if err := hfunc(ctxw, r); err != nil {
		log.Debugf("cmd %d handler: %v", r.Cmd, err)
		return err
	}
	return nil
}

// Add an handler function for some command.
func (h *Handler) Add(cmd int32, hf HandlerFunc) {
	h.router[cmd] = hf
}

// NewHandler returns a SawTK handler.
func NewHandler(family *Family) *Handler {
	return &Handler{
		family,
		make(map[int32]HandlerFunc),
		os.Getenv("TP_DEBUG") == "true",
	}
}
