package mysaga

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/PASPARTUUU/go_for_example/examples/ex-pattern-saga-implementation/mysaga/brocker"
	"github.com/PASPARTUUU/go_for_example/pkg/dye"
	"github.com/PASPARTUUU/go_for_example/pkg/tool"
	"github.com/gofrs/uuid"
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
)

type StepType string

const (
	Begining     StepType = "begining"
	Intermediate StepType = "intermediate"
	Finishing    StepType = "finishing"
)

var ErrNeedRedo = errors.New("ErrNeedRedo")
var ErrNeedUndo = errors.New("ErrNeedUndo")

type ParsingFunc func(data [][]byte) ([]interface{}, error)

// -------------------------------------------------

type msgBus interface {
	Publish(payload interface{}) error
	HandleNewEvent() <-chan []byte
}

type stepData struct {
	execParser func(data [][]byte) ([]interface{}, error)
	ExecParams []interface{}
	ExecFunc   interface{}

	compParser func(data [][]byte) ([]interface{}, error)
	CompParams []interface{}
	CompFunc   interface{}
}

type Step struct {
	TransactionName string
	Name            string
	ParentSteps     []string
	Type            StepType
	RedoCount       int
	stepData        stepData
	SagaStore       SagaStore

	bus msgBus
}

type Client struct {
	OwnerHost string
	rabbit    *brocker.Rabbit
	steps     []Step
}

func NewClient(ownerHost string, msgBus brocker.ConnCredits, queueName, consumerName string) (*Client, error) {

	rabb, err := brocker.NewConnection(msgBus)
	if err != nil {
		return nil, err
	}

	if queueName != "" {
		rabb.Queue = queueName
	}
	if consumerName != "" {
		rabb.Consumer = consumerName
	}

	rabb.Listen()

	return &Client{
		rabbit:    rabb,
		OwnerHost: ownerHost,
	}, nil
}

// -------------------------------------------------

func (c *Client) AddStep(
	transactionName string,
	name string,
	parentSteps []string,
	stepType StepType,
	execFunc interface{},
	compFunc interface{},
	redoCount int,
	store SagaStore,
) (*Step, error) {

	if transactionName == "" {
		return nil, errors.New("transaction name must not be empty")
	}

	if stepType == "" {
		stepType = Intermediate
	}

	if store == nil {
		store = &Store{
			data:         make(map[string]interface{}),
			execRes:      make(map[string]interface{}),
			redoTryCount: make(map[string]int),
		}
	}

	// TODO: проверить у execFunc и compFunc последний возвр аргумент должен быть ошибкой

	step := Step{
		TransactionName: transactionName,
		Name:            name,
		ParentSteps:     parentSteps,
		Type:            stepType,
		stepData: stepData{
			ExecFunc: execFunc,
			CompFunc: compFunc,

			execParser: defaultParse,
			compParser: defaultParse,
		},
		RedoCount: redoCount,
		SagaStore: store,
		bus:       c.rabbit,
	}
	c.steps = append(c.steps, step)

	return &step, nil
}

func (c *Client) GetStep(name string) *Step {
	for i, s := range c.steps {
		if s.Name == name {
			return &c.steps[i]
		}
	}
	return nil
}

// -------------------------------------------------

func (s *Step) ParseMe(exec, comp ParsingFunc) {
	s.stepData.execParser = exec
	s.stepData.compParser = comp
}

func defaultParse(b [][]byte) ([]interface{}, error) {
	var res []interface{}

	err := json.Unmarshal(b[len(b)-1], &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func parseAll(b [][]byte) ([][]interface{}, error) {
	var res [][]interface{}

	for _, bb := range b {
		var i []interface{}
		err := json.Unmarshal(bb, &i)
		if err != nil {
			return nil, err
		}

		res = append(res, i)
	}

	return res, nil
}

// -------------------------------------------------

// Play - тригерит запуск саги у одного из начальных шагов
func (s *Step) Play(args ...interface{}) error {

	b, err := json.Marshal(args)
	if err != nil {
		return err
	}

	var wp = make(WantPayload, 0)
	wp = append(wp, b)
	b, err = json.Marshal(wp)
	if err != nil {
		return err
	}

	data := DataFormat{
		Event:   fmt.Sprint(s.TransactionName, ".", s.Name, ".", ActionDone, ".", uuid.Must(uuid.NewV4()).String()),
		Payload: b,
	}

	outData, err := s.doScene(data)
	if err != nil {
		return err
	}

	s.bus.Publish(outData)
	if err != nil {
		return err
	}

	return nil
}

func (s *Step) Listen() {
	go func() {
		for {
			err := s.listen()
			if err != nil {
				log.Error(err)
			}
		}
	}()
}

func (s *Step) listen() error {

	fmt.Println("\nmsg waiting...")

	eventChannel := s.bus.HandleNewEvent()
	chanData := <-eventChannel
	var data DataFormat
	err := json.Unmarshal(chanData, &data)
	if err != nil {
		dye.Err(err)
		return err
	}
	fmt.Println("get data")
	fmt.Printf("%+v\n", data)

	if getEventClaim(data.Event, ClaimName) != s.TransactionName { // если шаг не принадлежит саге
		return nil
	}

	switch {
	case tool.SliceContainsString(s.ParentSteps, getEventClaim(data.Event, ClaimStep)) &&
		getEventClaim(data.Event, ClaimAction) == string(ActionDone): // неправильное условие тут ИЛИ а должно быть И (логика с счетчиком)

		_, err := s.SagaStore.GetIncomingData(getEventClaim(data.Event, ClaimKey))
		if err != nil && err != ErrNotFound {
			return err
		}

		outData, err := s.doScene(data)
		if err != nil {
			return err
		}

		if err := s.bus.Publish(outData); err != nil {
			return err
		}

	case s.Name == getEventClaim(data.Event, ClaimStep) &&
		getEventClaim(data.Event, ClaimAction) == string(ActionRedo):

		outData, err := s.doScene(data)
		if err != nil {
			return err
		}

		if err := s.bus.Publish(outData); err != nil {
			return err
		}

	case s.Name == getEventClaim(data.Event, ClaimStep) &&
		getEventClaim(data.Event, ClaimAction) == string(ActionUndo):

		if err = s.undoScene(data); err != nil {
			return err
		}

		for _, step := range s.ParentSteps {
			data.Event = renameEventClaim(data.Event, ClaimStep, step)

			if err := s.bus.Publish(data); err != nil {
				return err
			}
		}

	case data.Event == string(HealthCheck):
	case getEventClaim(data.Event, ClaimName) == string(ClaimConsistancy):
	default:
		fmt.Println("switch case not found")
	}

	return nil
}

func (s *Step) doScene(data DataFormat) (*DataFormat, error) {

	var wp WantPayload
	idpKey := getEventClaim(data.Event, ClaimKey)

	if err := json.Unmarshal(data.Payload, &wp); err != nil {
		return nil, err
	}

	paramsForStorage, err := parseAll(wp)
	if err != nil {
		return nil, err
	}
	s.SagaStore.SetIncomingData(idpKey, paramsForStorage)

	rawParams, err := s.stepData.execParser(wp)
	if err != nil {
		return nil, err
	}

	funcValue := reflect.ValueOf(s.stepData.ExecFunc)
	funcType := reflect.TypeOf(s.stepData.ExecFunc)

	types := make([]reflect.Type, 0, funcType.NumIn())
	for i := 0; i < funcType.NumIn(); i++ {
		types = append(types, funcType.In(i))
	}

	b, err := json.Marshal(rawParams)
	if err != nil {
		return nil, err
	}
	prms, err := unmarshalParams(types, b)
	if err != nil {
		return nil, err
	}

	vals := funcValue.Call(prms)
	if callerr := isReturnError(vals); callerr != nil {
		redoTry, err := s.SagaStore.GetRedoCount(idpKey)
		if err != nil {
			return nil, err
		}

		if redoTry < s.RedoCount {
			if err := s.SagaStore.IncRedoCount(idpKey); err != nil {
				return nil, err
			}
			dye.Next(ErrNeedRedo)
			// отправить реду
			data.Event = renameEventClaim(data.Event, ClaimAction, string(ActionRedo))
			if err := s.bus.Publish(data); err != nil {
				return nil, err
			}
			return nil, errors.Wrap(callerr, ErrNeedRedo.Error())
		}

		// отправить анду
		dye.Next(ErrNeedUndo)
		data.Event = renameEventClaim(data.Event, ClaimAction, string(ActionUndo))
		data.Event = renameEventClaim(data.Event, ClaimStep, s.ParentSteps[0])
		if err := s.bus.Publish(data); err != nil {
			return nil, err
		}
		return nil, errors.Wrap(callerr, ErrNeedUndo.Error())
	}

	var resSaving = make([]interface{}, 0)
	for i, v := range vals {
		if i == len(vals)-1 { // дабы ошибка не отправлялась
			break
		}
		resSaving = append(resSaving, v.Interface())
	}
	bRes, err := json.Marshal(resSaving)
	if err != nil {
		return nil, err
	}

	wp = append(wp, bRes)
	s.SagaStore.SetOutgoingData(idpKey, resSaving)

	outData, err := s.makeOutgoingData(data.Event, wp)
	if err != nil {
		return nil, err
	}

	return outData, nil
}

func (s *Step) undoScene(data DataFormat) error {

	var wp WantPayload
	idpKey := getEventClaim(data.Event, ClaimKey)

	pl, err := s.SagaStore.GetOutgoingData(idpKey)
	if err != nil {
		return err
	}

	bpl, err := json.Marshal(pl)
	if err != nil {
		return err
	}

	wp = append(wp, bpl)

	paramsFromStorage, err := s.stepData.execParser(wp)
	if err != nil {
		return err
	}

	funcValue := reflect.ValueOf(s.stepData.CompFunc)
	funcType := reflect.TypeOf(s.stepData.CompFunc)

	types := make([]reflect.Type, 0, funcType.NumIn())
	for i := 0; i < funcType.NumIn(); i++ {
		types = append(types, funcType.In(i))
	}

	b, err := json.Marshal(paramsFromStorage)
	if err != nil {
		return err
	}
	prms, err := unmarshalParams(types, b)
	if err != nil {
		return err
	}

	vals := funcValue.Call(prms)
	if callerr := isReturnError(vals); callerr != nil {
		return callerr
	}

	/*
			var resSaving = make([]interface{}, 0)
			for i, v := range vals {
				if i == len(vals)-1 { // дабы ошибка не отправлялась
					break
				}
				resSaving = append(resSaving, v.Interface())
			}
			bRes, err := json.Marshal(resSaving)
			if err != nil {
				return nil, err
			}

			wp = append(wp, bRes)
			s.SagaStore.SetOutgoingData(idpKey, wp)

		outData, err := s.makeOutgoingData(data.Event, wp)
		if err != nil {
			return nil, err
		}
	*/
	return nil
}

func (s *Step) makeOutgoingData(event string, payload interface{}) (*DataFormat, error) {
	plRes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return &DataFormat{
		Event:   renameEventClaim(event, ClaimStep, s.Name),
		Payload: plRes,
	}, nil
}

// -------------------------------------------------

func (s *Step) SetCompParams(args ...interface{}) {
	var params = make([]interface{}, 0, len(args))
	for _, a := range args {
		params = append(params, a)
	}
	s.stepData.ExecParams = params
}

// -------------------------------------------------

func unmarshalParams(types []reflect.Type, payload []byte) ([]reflect.Value, error) {
	rawVals := make([]interface{}, 0, len(types))
	for _, typ := range types {
		rawVals = append(rawVals, reflect.New(typ).Interface())
	}

	json.Unmarshal(payload, &rawVals)
	res := make([]reflect.Value, 0, len(types))

	for i := 0; i < len(rawVals); i++ {
		objV := reflect.ValueOf(rawVals[i])

		if rawVals[i] == nil {
			objV = reflect.Zero(types[i])
		} else if reflect.TypeOf(rawVals[i]).Kind() == reflect.Ptr && objV.Type() != types[i] {
			objV = objV.Elem()
		}

		res = append(res, objV)
	}
	return res, nil
}

func isReturnError(result []reflect.Value) error {
	if len(result) > 0 && !result[len(result)-1].IsNil() {
		return result[len(result)-1].Interface().(error)
	}
	return nil
}

// -------------------------------------------------
