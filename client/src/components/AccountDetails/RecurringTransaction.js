import React, { Component } from 'react';
import PropTypes from 'prop-types';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { toCurrencyString } from '../../util';
import { SCHEDULE_TYPES } from '../../constants';

const numToWeekday = {
    0: 'Sunday',
    1: 'Monday',
    2: 'Tuesday',
    3: 'Wednesday',
    4: 'Thursday',
    5: 'Friday',
    6: 'Saturday'
}

export default class RecurringTransaction extends Component {
    static propTypes = {
        recurringTransaction: ImmutablePropTypes.contains({
            id: PropTypes.number.isRequired,
            transaction: ImmutablePropTypes.contains({
                name: PropTypes.string.isRequired,
                date: PropTypes.instanceOf(Date),
                category: PropTypes.string.isRequired,
                amount: PropTypes.number.isRequired,
                accountId: PropTypes.number.isRequired
            }).isRequired,
            scheduleType: PropTypes.oneOf(Object.values(SCHEDULE_TYPES)).isRequired,
            secondsBetween: PropTypes.number,
            dayOf: PropTypes.number,
            secondsBeforeToPost: PropTypes.number.isRequired
        }),
        currency: ImmutablePropTypes.contains({
            code: PropTypes.string.isRequired,
            digitsAfterDecimal: PropTypes.number.isRequired
        }),
        enterEditMode: PropTypes.func.isRequired
    }

    getScheduleString = () => {
        let { recurringTransaction } = this.props;
        let scheduleString = 'Generated ';
        switch(recurringTransaction.get('scheduleType')) {
            case SCHEDULE_TYPES.FIXED_DAY_WEEK:
                scheduleString += 'every ' + numToWeekday[recurringTransaction.get('dayOf')];
                break;

            case SCHEDULE_TYPES.FIXED_DAY_MONTH:
                scheduleString += 'on the ' + recurringTransaction.get('dayOf') + ' of every month';
                break;

            case SCHEDULE_TYPES.FIXED_DAY_YEAR:
                scheduleString += 'on the ' + recurringTransaction.get('dayOf') + ' of every year';
                break;

            case SCHEDULE_TYPES.FIXED_INTERVAL:
                scheduleString += 'every ' + (recurringTransaction.get('secondsBetween') / (60 * 60 * 24)) + ' days';
                break;

            default:
                break;
        }

        scheduleString += ', and will be posted ' + (recurringTransaction.get('secondsBeforeToPost') / (60 * 60 * 24)) + ' days before occuring.';
        return scheduleString;
    }

    render() {
        const {
            recurringTransaction,
            currency,
            enterEditMode
        } = this.props;

        return (
                <div className="recurringTransaction" onClick={ enterEditMode } >
                    <div>{ recurringTransaction.getIn(['transaction', 'name']) }</div>
                    <div>{ recurringTransaction.getIn(['transaction', 'date']).toLocaleDateString() }</div>
                    <div>{ recurringTransaction.getIn(['transaction', 'category']) }</div>
                    <div>{ toCurrencyString(recurringTransaction.getIn(['transaction', 'amount']), currency.get('code'), currency.get('digitsAfterDecimal')) }</div>
                    <div className="scheduleString">{ this.getScheduleString() }</div>
                </div>
                );
    }
}
