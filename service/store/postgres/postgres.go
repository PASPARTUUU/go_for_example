package postgres

import (
	"fmt"
	"time"

	"github.com/go-pg/pg"
	"github.com/sirupsen/logrus"
	"go.uber.org/atomic"

	"github.com/PASPARTUUU/go_for_example/pkg/errpath"
	"github.com/PASPARTUUU/go_for_example/service/config"
	"github.com/PASPARTUUU/go_for_example/service/store/repo"
)

// KeepAlivePollPeriod - is a Pg keepalive check time period
const keepAlivePollPeriod = time.Second * 60

// Pg -
type Pg struct {
	DB *pg.DB
	//---
	User         repo.User
	PostgresUser PostgresUser
	//---
	cfg config.Postgres
}

// Tracer -
var Tracer *DBQueryTraceHook

// -------------------------------------------------
// -------------------------------------------------
// -------------------------------------------------

// NewConnect -
func NewConnect(cfg config.Postgres) (*Pg, error) {

	// Connect to the db and remember to close it
	db := pg.Connect(&pg.Options{
		Addr:     fmt.Sprintf("%s:%v", cfg.Host, cfg.Port),
		User:     cfg.User,
		Password: cfg.Password,
		Database: cfg.DBName,
	})

	// Test connection
	var ping int
	_, err := db.QueryOne(pg.Scan(&ping), "SELECT 1")
	if err != nil {
		return nil, errpath.Err(err, "failed to connect to the db")
	}

	Tracer = InitDebugSQLQueryHook(db)

	pg := Pg{
		DB:           db,
		User:         NewUserRepo(db),
		PostgresUser: NewUserRepo(db),
	}
	if db != nil {
		go pg.keepAlive()
	}

	return &pg, nil
}

// -------------------------------------------------

// keepAlive - makes sure PostgreSQL is alive and reconnects if needed
func (pg *Pg) keepAlive() {
	log := logrus.WithField("event", "KeepAlivePg")
	var err error
	for {
		// Check if PostgreSQL is alive every 'KeepAlivePollPeriod' seconds
		time.Sleep(keepAlivePollPeriod)
		lostConnect := false
		if pg == nil {
			lostConnect = true
		} else if _, err = pg.DB.Exec("SELECT 1"); err != nil {
			lostConnect = true
		}
		if !lostConnect {
			continue
		}
		log.Warnln(errpath.Infof("Lost PostgreSQL connection. Restoring..."))

		pg, err = NewConnect(pg.cfg)
		if err != nil {
			log.Errorln(errpath.Err(err))
			continue
		}
		log.Infoln(errpath.Infof("PostgreSQL reconnected"))
	}
}

// -------------------------------------------------
// -------------------------------------------------
// -------------------------------------------------

// DBQueryTraceHook -
type DBQueryTraceHook struct {
	enableCounter *atomic.Int32
}

// InitDebugSQLQueryHook -
func InitDebugSQLQueryHook(conn *pg.DB) *DBQueryTraceHook {
	hook := DBQueryTraceHook{
		enableCounter: atomic.NewInt32(0),
	}
	conn.AddQueryHook(&hook)
	return &hook
}

// BeforeQuery -
func (db *DBQueryTraceHook) BeforeQuery(q *pg.QueryEvent) {
	if db.enableCounterVal() > 0 {
		query, err := q.FormattedQuery()
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		}
		fmt.Printf(errpath.InfofWithFuncCaller(12, "\x1b[35m %s \x1b[0m\n\n", query))

		db.enableCounterDec()
	}
}

// AfterQuery -
func (db *DBQueryTraceHook) AfterQuery(q *pg.QueryEvent) {}

// DebugQueryHook - активирует логирование SQL запроса
type DebugQueryHook interface {
	// StartTrace - начать логировать SQL
	StartTrace()
	enableCounterInc()
	enableCounterDec()
	enableCounterVal() int
}

// StartTrace - начать логировать SQL
func (db *DBQueryTraceHook) StartTrace() { db.enableCounterInc() }

func (db *DBQueryTraceHook) enableCounterInc()     { db.enableCounter.Inc() }
func (db *DBQueryTraceHook) enableCounterDec()     { db.enableCounter.Dec() }
func (db *DBQueryTraceHook) enableCounterVal() int { return int(db.enableCounter.Load()) }
