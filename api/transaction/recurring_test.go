// +build integration

package transaction

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
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
		require.NoError(suite.T(), err, "failed to create recurring transaction: %+v", tr)
	}

	retrieved, err := GetRecurring(suite.Ctx, acc.Id)
	require.NoError(suite.T(), err, "unable to retrieve recurring transactions")

	for _, tr := range retrieved {
		initial := recurringTransactions[tr.Transaction.Name]
		switch tr.Transaction.Name {
		case "fixed interval", "fixed interval, starting 25 days ago":
			checkRecurringTransactionEqual(suite.T(), initial, tr, true)
		case "tuesdays, starting 25 days ago":
			checkRecurringTransactionEqual(suite.T(), initial, tr, false)

			// ensure the date is a tuesday
			require.Equal(suite.T(), time.Tuesday, tr.Transaction.Date.Weekday(), "transaction should be scheduled for a tuesday")
			// it should be 19-25 days ago
			lowerBound := now.AddDate(0, 0, -25)
			upperBound := now.AddDate(0, 0, -19)
			require.True(suite.T(), tr.Transaction.Date.After(lowerBound), "initial scheduled date should be after 25 days ago, %+v, but is %+v", now.AddDate(0, 0, -25), tr.Transaction.Date)
			require.True(suite.T(), tr.Transaction.Date.Before(upperBound), "initial scheduled date should be before 19 days ago, %+v, but is %+v", now.AddDate(0, 0, -19), tr.Transaction.Date)
		}
	}
}

func (suite *RecurringTestSuite) TestGenerateFuture() {
	// create account
	acc := integration.NewAccount(suite.T(), suite.Ctx)

	now := time.Now()
	currDay := util.WeekdayToInt(now.Weekday())
	tr := RecurringTransaction{
		Transaction: Transaction{
			Name:      "future",
			Date:      now.AddDate(0, 0, 25),
			Category:  "fun1",
			Amount:    -505,
			Note:      "note1",
			AccountId: acc.Id,
		},
		ScheduleType:        constants.FIXED_DAY_WEEK,
		DayOf:               &currDay,
		SecondsBeforeToPost: 1,
	}
	_, err := NewRecurring(suite.Ctx, &tr)
	require.NoError(suite.T(), err, "failed to create recurring transaction: %+v", tr)

	tr = RecurringTransaction{
		Transaction: Transaction{
			Name:      "futureButPosted",
			Date:      now.AddDate(0, 0, 1),
			Category:  "fun2",
			Amount:    -503,
			Note:      "note2",
			AccountId: acc.Id,
		},
		ScheduleType:        constants.FIXED_DAY_WEEK,
		DayOf:               &currDay,
		SecondsBeforeToPost: 60 * 60 * 24 * 7,
	}
	_, err = NewRecurring(suite.Ctx, &tr)
	require.NoError(suite.T(), err, "failed to create recurring transaction: %+v", tr)

	err = GenRecurringTransactions(suite.Ctx)
	require.NoError(suite.T(), err, "failed to generate recurring transactions")

	retrieved, err := GetFuture(suite.Ctx, acc.Id, &now)
	require.NoError(suite.T(), err, "failed to retrieve transactions after generating")

	require.Equal(suite.T(), 1, len(retrieved), "should be one transaction generated")
	generated := retrieved[0]
	checkTransactionEqual(suite.T(), tr.Transaction, generated)

	// verify that the date got shifted forward
	recurring, err := GetRecurring(suite.Ctx, acc.Id)
	require.NoError(suite.T(), err, "unable to retrieve recurring transactions")

	for _, rt := range recurring {
		if rt.Transaction.Name == "futureButPosted" {
			checkRecurringTransactionEqual(suite.T(), tr, rt, false)
			require.Equal(suite.T(), currDay, util.WeekdayToInt(rt.Transaction.Date.Weekday()), "transaction should be scheduled for the required day")
			// it should be at least 7 days from now
			sevenDays := now.AddDate(0, 0, 7)
			require.True(suite.T(), rt.Transaction.Date.After(sevenDays), "next run should be at least 7 days from now, since earliest is tomorrow, but actual is: %+v and should be after %+v", rt.Transaction.Date, sevenDays)
		}
	}
}

func TestRecurringTestSuite(t *testing.T) {
	suite.Run(t, new(RecurringTestSuite))
}

func checkTransactionEqual(t *testing.T, expected, actual Transaction) {
	require.Equal(t, expected.Name, actual.Name, "actual name should be same as expected")
	require.Equal(t, expected.AccountId, actual.AccountId, "actual account id should be same as expected")
	require.Equal(t, expected.Amount, actual.Amount, "actual amount should be same as expected")
	require.Equal(t, expected.Category, actual.Category, "actual category should be same as expected")
	require.Equal(t, expected.Note, actual.Note, "actual note should be same as expected")
	require.Equal(t, expected.Date.Year(), actual.Date.Year(), "actual year should be same as expected")
	require.Equal(t, expected.Date.Month(), actual.Date.Month(), "actual month should be same as expected")
	require.Equal(t, expected.Date.Day(), actual.Date.Day(), "actual day should be same as expected")
}

func checkRecurringTransactionEqual(t *testing.T, expected, actual RecurringTransaction, checkDate bool) {
	require.Equal(t, expected.Transaction.AccountId, actual.Transaction.AccountId, "expected account id should be same as expected")
	require.Equal(t, expected.Transaction.Amount, actual.Transaction.Amount, "expected amount should be same as expected")
	require.Equal(t, expected.Transaction.Category, actual.Transaction.Category, "expected category should be same as expected")
	require.Equal(t, expected.Transaction.Note, actual.Transaction.Note, "expected note should be same as expected")

	require.Equal(t, expected.ScheduleType, actual.ScheduleType, "expected schedule type should be same as expected")
	require.Equal(t, expected.SecondsBetween, actual.SecondsBetween, "expected seconds between should be same as expected")
	require.Equal(t, expected.DayOf, actual.DayOf, "expected day of should be same as expected")
	require.Equal(t, expected.SecondsBeforeToPost, actual.SecondsBeforeToPost, "expected seconds before to post should be same as expected")

	if checkDate {
		require.Equal(t, expected.Transaction.Date.Year(), actual.Transaction.Date.Year(), "expected year should be same as expected")
		require.Equal(t, expected.Transaction.Date.Month(), actual.Transaction.Date.Month(), "expected month should be same as expected")
		require.Equal(t, expected.Transaction.Date.Day(), actual.Transaction.Date.Day(), "expected day should be same as expected")
	}
}
