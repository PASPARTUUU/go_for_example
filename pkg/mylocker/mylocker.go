package mylocker

import (
	"sync"
	"time"
)

// locker -
// rus: обертка над mutex с временем взаимодействия с ним
// use: NewLocker()
type locker struct {
	sync.RWMutex
	takenAt time.Time
}

func newLocker() *locker {
	return &locker{
		RWMutex: sync.RWMutex{},
		takenAt: time.Now(),
	}
}

func (lc *locker) lock() {
	lc.RWMutex.Lock()
	lc.takenAt = time.Now()
}

func (lc *locker) unlock() {
	lc.takenAt = time.Now()
	lc.RWMutex.Unlock()
}

func (lc *locker) rlock() {
	lc.RWMutex.RLock()
	lc.takenAt = time.Now()
}

func (lc *locker) runlock() {
	lc.takenAt = time.Now()
	lc.RWMutex.RUnlock()
}

func (lc *locker) execute(job func() error) error {
	lc.RWMutex.Lock()
	defer lc.RWMutex.Unlock()
	lc.takenAt = time.Now()

	return job()
}

// -------------------------------------------------
// -------------------------------------------------
// -------------------------------------------------

// Entity -
// use: NewEntity()
type Entity struct {
	Monitor monitor
}

// monitor -
// rus: имеет методы для получения мета инфы о локкерах сущности
type monitor struct {
	lockers map[string]*locker
	m       *sync.Mutex
}

// NewEntity -
func NewEntity() *Entity {
	return &Entity{
		Monitor: monitor{
			lockers: make(map[string]*locker),
			m:       &sync.Mutex{},
		},
	}
}

// getOrCreateLocker -
func (e *Entity) getOrCreateLocker(name string) *locker {
	e.Monitor.m.Lock()
	defer e.Monitor.m.Unlock()

	if l, found := e.Monitor.lockers[name]; found {
		return l
	}

	e.Monitor.lockers[name] = newLocker()

	return e.Monitor.lockers[name]
}

// RemoveLockers -
func (e *Entity) RemoveLockers(names ...string) {
	e.Monitor.m.Lock()
	defer e.Monitor.m.Unlock()

	e.removeLockers(names...)
}

func (e *Entity) removeLockers(names ...string) {
	for _, n := range names {
		delete(e.Monitor.lockers, n)
	}
}

// RemoveOverdue -
// rus: удалить просроченные
func (e *Entity) RemoveOverdue(expirationTime time.Time) {
	e.Monitor.m.Lock()
	defer e.Monitor.m.Unlock()

	var forRemove []string
	for id, val := range e.Monitor.lockers {
		if val.takenAt.After(expirationTime) {
			forRemove = append(forRemove, id)
		}
	}

	e.removeLockers(forRemove...)
}

// Lock -
func (e *Entity) Lock(key string) {
	lc := e.getOrCreateLocker(key)
	lc.lock()
}

// Unlock -
func (e *Entity) Unlock(key string) {
	lc := e.getOrCreateLocker(key)
	lc.unlock()
}

// RLock -
func (e *Entity) RLock(key string) {
	lc := e.getOrCreateLocker(key)
	lc.rlock()
}

// RUnlock -
func (e *Entity) RUnlock(key string) {
	lc := e.getOrCreateLocker(key)
	lc.runlock()
}

// Execute -
func (e *Entity) Execute(key string, job func() error) error {
	return e.getOrCreateLocker(key).execute(job)
}

// -------------------------------------------------

// GetLockersList -
func (e *monitor) GetLockersList() []string {
	e.m.Lock()
	defer e.m.Unlock()

	var res []string
	for k := range e.lockers {
		res = append(res, k)
	}
	return res
}

// Len -
func (e *monitor) Len() int {
	e.m.Lock()
	defer e.m.Unlock()

	return len(e.lockers)
}

// -------------------------------------------------
// Example

/*

var statusEntity *mylocker.Entity

func init() {
	statusEntity = mylocker.NewEntity()
}

type Status struct {
	UUID  string
	State string
}

func (st *Status) SetState(newState string) {
	statusEntity.Lock(st.UUID)
	defer statusEntity.Unlock(st.UUID)

	st.State = newState
}

func (st *Status) GetState(newState string) {
	statusEntity.RLock(st.UUID)
	defer statusEntity.RUnlock(st.UUID)

	st.State = newState
}

*/
