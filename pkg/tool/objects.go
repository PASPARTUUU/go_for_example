package tool

import (
	"time"

	"github.com/PASPARTUUU/go_for_example/pkg/errpath"
	"github.com/go-pg/pg/urlvalues"
)

// Объекты общего назначения

const (
	// Pi -
	Pi float64 = 3.14
)

// Constant - constant
type Constant string

// ConstantF - constant
// meta: go должен инлайнить (не вызвать функцию, а подставлять сразу результат) возвращаемое значение
type ConstantF func() string

// URL - url адрес
type URL string

// Coordinates -
type Coordinates struct {
	Long float64 `json:"long"`
	Lat  float64 `json:"lat"`
}

// PictureSize -
type PictureSize struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// -------------------------------------------------
// -------------------------------------------------
// -------------------------------------------------

// Paginator -
type Paginator struct {
	Page  int             `json:"page"`
	Limit int             `json:"limit"`
	Pager urlvalues.Pager `json:"-"`
}

// PaginatorResult -
type PaginatorResult struct {
	List         []interface{} `json:"list"`
	RecordsCount int           `json:"records_count"`
}

// Fasten -
func (p *Paginator) Fasten() error {
	if p.Page == 0 || p.Limit == 0 {
		return errpath.Errorf("empty pager params")
	}
	p.Pager.Limit = p.Limit
	p.Pager.SetPage(p.Page)
	return nil
}

// -------------------------------------------------
// -------------------------------------------------
// -------------------------------------------------

// DataBaseConfigParams - модель для хранения конфиг-параметров в базе
type DataBaseConfigParams struct {
	UUID        string    `sql:"uuid"`
	Key         string    `sql:"key"`         // уникальное имя конфига
	Value       string    `sql:"value"`       // параметры в json формате
	description string    `sql:"description"` // описание конфига
	CreatedAt   time.Time `sql:"created_at"`
	UpdatedAt   time.Time `sql:"updated_at"`
}

// -------------------------------------------------
// -------------------------------------------------
// -------------------------------------------------
