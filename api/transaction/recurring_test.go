// +build integration

package transaction

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/jchorl/financejc/api/util"
	"github.com/jchorl/financejc/constants"
	"github.com/jchorl/financejc/integration"
)

type RecurringTestSuite struct {
	suite.Suite
	Ctx context.Context
}

func (suite *RecurringTestSuite) SetupTest() {
	db := integration.FreshDb(suite.T())
	ctx := integration.ContextWithUserAndDB(0, db)
	uid := integration.NewUser(suite.T(), ctx)
	suite.Ctx = integration.ContextWithUserAndDB(uid, db)
}

func (suite *RecurringTestSuite) TestNewRecurring() {
	// create account
	acc := integration.NewAccount(suite.T(), suite.Ctx)

	// 12 days
	secondsBetween := int((time.Second * time.Duration(60*60*24*12)).Seconds())
	tuesday := util.WeekdayToInt(time.Tuesday)
	now := time.Now()
	recurringTransactions := map[string]RecurringTransaction{
		"fixed interval": RecurringTransaction{
			Transaction: Transaction{
				Name:      "fixed interval",
				Date:      now,
				Category:  "fun",
				Amount:    -500,
				Note:      "note",
				AccountId: acc.Id,
			},
			ScheduleType:        constants.FIXED_INTERVAL,
			SecondsBetween:      &secondsBetween,
			SecondsBeforeToPost: 0,
		},
		"tuesdays, starting 25 days ago": RecurringTransaction{
			Transaction: Transaction{
				Name:      "tuesdays, starting 25 days ago",
				Date:      now.AddDate(0, 0, -25),
				Category:  "fun1",
				Amount:    -505,
				Note:      "note1",
				AccountId: acc.Id,
			},
			ScheduleType:        constants.FIXED_DAY_WEEK,
			DayOf:               &tuesday,
			SecondsBeforeToPost: 0,
		},
		"fixed interval, starting 25 days ago": RecurringTransaction{
			Transaction: Transaction{
				Name:      "fixed interval, starting 25 days ago",
				Date:      now.AddDate(0, 0, -25),
				Category:  "fun2",
				Amount:    -502,
				Note:      "note2",
				AccountId: acc.Id,
			},
			ScheduleType:        constants.FIXED_INTERVAL,
			SecondsBetween:      &secondsBetween,
			SecondsBeforeToPost: 0,
		},
	}

	for _, tr := range recurringTransactions {
		_, err := NewRecurring(suite.Ctx, &tr)
		assert.NoError(suite.T(), err, "failed to create recurring transaction: %+v", tr)
	}

	retrieved, err := GetRecurring(suite.Ctx, acc.Id)
	assert.NoError(suite.T(), err, "unable to retrieve recurring transactions")

	for _, tr := range retrieved {
		initial := recurringTransactions[tr.Transaction.Name]
		assert.Equal(suite.T(), initial.Transaction.AccountId, tr.Transaction.AccountId, "initial account id should be same as retrieved recurring transaction")
		assert.Equal(suite.T(), initial.Transaction.Amount, tr.Transaction.Amount, "initial amount should be same as retrieved recurring transaction")
		assert.Equal(suite.T(), initial.Transaction.Category, tr.Transaction.Category, "initial category should be same as retrieved recurring transaction")
		assert.Equal(suite.T(), initial.Transaction.Note, tr.Transaction.Note, "initial note should be same as retrieved recurring transaction")

		assert.Equal(suite.T(), initial.ScheduleType, tr.ScheduleType, "initial schedule type should be same as retrieved recurring transaction")
		assert.Equal(suite.T(), initial.SecondsBetween, tr.SecondsBetween, "initial seconds between should be same as retrieved recurring transaction")
		assert.Equal(suite.T(), initial.DayOf, tr.DayOf, "initial day of should be same as retrieved recurring transaction")
		assert.Equal(suite.T(), initial.SecondsBeforeToPost, tr.SecondsBeforeToPost, "initial seconds before to post should be same as retrieved recurring transaction")

		switch tr.Transaction.Name {
		case "fixed interval", "fixed interval, starting 25 days ago":
			assert.Equal(suite.T(), initial.Transaction.Date.Year(), tr.Transaction.Date.Year(), "fixed interval transaction date should not be changed on save")
			assert.Equal(suite.T(), initial.Transaction.Date.Month(), tr.Transaction.Date.Month(), "fixed interval transaction date should not be changed on save")
			assert.Equal(suite.T(), initial.Transaction.Date.Day(), tr.Transaction.Date.Day(), "fixed interval transaction date should not be changed on save")
		case "tuesdays, starting 25 days ago":
			// ensure the date is a tuesday
			assert.Equal(suite.T(), time.Tuesday, tr.Transaction.Date.Weekday(), "transaction should be scheduled for a tuesday")
			// it should be 19-25 days ago
			lowerBound := now.AddDate(0, 0, -25)
			upperBound := now.AddDate(0, 0, -19)
			assert.True(suite.T(), tr.Transaction.Date.After(lowerBound), "initial scheduled date should be after 25 days ago, %+v, but is %+v", now.AddDate(0, 0, -25), tr.Transaction.Date)
			assert.True(suite.T(), tr.Transaction.Date.Before(upperBound), "initial scheduled date should be before 19 days ago, %+v, but is %+v", now.AddDate(0, 0, -19), tr.Transaction.Date)
		}
	}
}

func TestRecurringTestSuite(t *testing.T) {
	suite.Run(t, new(RecurringTestSuite))
}
