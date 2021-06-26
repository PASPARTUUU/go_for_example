package mysaga

import (
	"strings"
)

// -------------------------------------------------

type ExecuteFunc func(params interface{}) (interface{}, error)
type CompensateFunc func(params interface{}) (interface{}, error)

type DataFormat struct {
	Event   string
	Payload []byte
}

type WantPayload [][]byte // [step][data]

// -------------------------------------------------
// структура ивента:
// название_транзакции.название_ивента-шага.действие.ключу_идемпотентности
// name.step.action.key
// billing.money-reserved.done.kfir12-mvkfrn233-kj4jh5vnjfj

type Claim uint

const (
	ClaimName Claim = iota
	ClaimStep
	ClaimAction
	ClaimKey
)

type ReservedEvents string

const (
	HealthCheck ReservedEvents = "health_check"
)

type ReservedEventsFClaim string // зарезервированные клеймы имею

const (
	ClaimStart  ReservedEventsFClaim = "start"
	ClaimFinish ReservedEventsFClaim = "finish"

	ClaimConsistancy ReservedEventsFClaim = "consistancy" // ивент запрос подтверждения наличия записи по ключу идемпотентности
)

type Action string

const (
	ActionDone Action = "done" // исполнился (с учетом имеющейся записи с клюм идемпотентности)
	ActionUndo Action = "undo" // откатить
	ActionRedo Action = "redo" // повторить (без учета имеющейся записи с клюм идемпотентности)
)

// -------------------------------------------------

func getEventClaim(event string, num Claim) string {
	claims := strings.Split(event, ".")

	if len(claims) > int(num) {
		return claims[num]
	}

	return ""
}

func renameEventClaim(event string, num Claim, newName string) string {
	claims := strings.Split(event, ".")

	if len(claims) > int(num) {
		claims[num] = newName
		return strings.Join(claims, ".")
	}

	return event
}
