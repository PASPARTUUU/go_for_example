package tickers

import (
	"context"
	"time"

	"github.com/PASPARTUUU/go_for_example/pkg/dye"
	"github.com/PASPARTUUU/go_for_example/pkg/errpath"
	"github.com/PASPARTUUU/go_for_example/service/models"
	"github.com/sirupsen/logrus"
)

const checkUserRegPeriod = time.Second * 60
const registrationTimeLength = time.Second * 120

// startCheckUserRegTime - логирует пользователя зарегистрированного n-времени назад
func (ticker *Ticker) startCheckUserRegTime() {
	wg.Add(1)
	go ticker.checkUserRegTimeTicker(context.Background())
}

func (ticker *Ticker) checkUserRegTimeTicker(ctx context.Context) {
	defer wg.Done()

	for {
		select {
		case <-closed:
			return
		case <-time.After(checkUserRegPeriod):
			var err error
			log := logrus.WithField("event", "check user registration time")

			var users []models.User
			ticker.handler.Storage.Pg.DB.Model(&users).
				Select()
			if err != nil {
				log.WithField("reason", "error getting users").Error(errpath.Err(err).Error())
				continue
			}

			for _, u := range users {
				if u.CreatedAt.Add(registrationTimeLength).After(time.Now()) &&
					u.CreatedAt.Add(registrationTimeLength).Before(time.Now().Add(checkUserRegPeriod)) { // второе условие для единоразового вывода
					dye.New(u.Name, " long time with us").TextBlue().Print()
				}
			}

		}
	}

}
