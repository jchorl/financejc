import React, { Component } from 'react';
import PropTypes from 'prop-types';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { connect } from 'react-redux';
import { fromRFC3339, toRFC3339 } from '../../util';
import { NEW_RECURRING_TRANSACTION_ID, SCHEDULE_TYPES } from '../../constants';
import { editRecurringTransaction, deleteRecurringTransaction, newRecurringTransaction } from '../../actions/accountData';
import SuggestedInput from '../SuggestedInput';

class RecurringTransactionEdit extends Component {
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
        dispatch: PropTypes.func.isRequired,
        exitEditMode: PropTypes.func.isRequired
    }

    constructor(props) {
        super(props);

        const {
            recurringTransaction
        } = this.props;
        const transaction = recurringTransaction.get('transaction');

        this.state = {
            name: transaction.get('name'),
            date: transaction.get('date'),
            category: transaction.get('category'),
            amount: transaction.get('amount'),
            scheduleType: recurringTransaction.get('scheduleType'),
            secondsBetween: recurringTransaction.get('secondsBetween'),
            dayOf: recurringTransaction.get('dayOf'),
            secondsBeforeToPost: recurringTransaction.get('secondsBeforeToPost'),

            savedName: '',
            savedCategory: '',
            savedAmount: 0,
            isSaved: false
        }
    }

    handleChange = field => e => {
        let newState = {};
        switch (field) {
            case 'date':
                newState[field] = fromRFC3339(e.target.value);
                break;

            case 'amount': {
                if (e.target.value === '') {
                    newState['amount'] = '';
                    break;
                }

                const { currency } = this.props;

                newState[field] = parseFloat(e.target.value, 10) * Math.pow(10, currency.get('digitsAfterDecimal'));
                break;
            }

            case 'dayOf':
                newState[field] = parseInt(e.target.value, 10);
                break;

            case 'secondsBetween':
                newState[field] = parseInt(e.target.value, 10) * 60 * 60 * 24;
                break;

            case 'secondsBeforeToPost':
                newState[field] = parseInt(e.target.value, 10) * 60 * 60 * 24;
                break;

            default:
                newState[field] = e.target.value;
        }

        this.setState(newState);
    }

    deleteRecurringTransaction = () => {
        const {
            dispatch,
            recurringTransaction,
        } = this.props;

        dispatch(deleteRecurringTransaction(recurringTransaction.toJS()));
    }

    onSubmit = e => {
        const {
            dispatch,
            exitEditMode,
            recurringTransaction: originalRecurringTransaction
        } = this.props;

        const {
            name,
            date,
            category,
            amount,
            scheduleType,
            secondsBetween,
            dayOf,
            secondsBeforeToPost
        } = this.state;

        const recurringTransaction = {
            transaction: {
                name,
                date,
                category,
                amount,
                accountId: originalRecurringTransaction.getIn(['transaction', 'accountId'])
            },
            scheduleType,
            secondsBetween,
            dayOf,
            secondsBeforeToPost
        };

        if (originalRecurringTransaction.get('id') === NEW_RECURRING_TRANSACTION_ID) {
            dispatch(newRecurringTransaction(recurringTransaction))
                .then(exitEditMode);
        } else {
            recurringTransaction.id = originalRecurringTransaction.get('id');
            dispatch(editRecurringTransaction(recurringTransaction))
                .then(exitEditMode);
        }

        e.preventDefault();
    }

    selectNameSuggestion = suggestion => {
        this.setState({
            name: suggestion.get('name'),
            category: suggestion.get('category'),
            amount: suggestion.get('amount')
        });
    }

    hoverNameSuggestion = suggestion => {
        const {
            name,
            category,
            amount,
            isSaved
        } = this.state;

        let newState = {
            name: suggestion.get('name'),
            category: suggestion.get('category'),
            amount: suggestion.get('amount')
        };

        if (!isSaved) {
            newState.savedName = name;
            newState.savedCategory = category;
            newState.savedAmount = amount;
            newState.isSaved = true;
        }

        this.setState(newState);
    }

    unhoverNameSuggestion = () => {
        const {
            savedName,
            savedCategory,
            savedAmount
        } = this.state;

        this.setState({
            name: savedName,
            category: savedCategory,
            amount: savedAmount,
            isSaved: false
        });
    }

    selectCategorySuggestion = suggestion => {
        this.setState({
            category: suggestion.get('category')
        });
    }

    hoverCategorySuggestion = suggestion => {
        const {
            category,
            isSaved
        } = this.state;

        let newState = {
            category: suggestion.get('category')
        };

        if (!isSaved) {
            newState.savedCategory = category;
            newState.isSaved = true;
        }

        this.setState(newState);
    }

    unhoverCategorySuggestion = () => {
        const {
            savedCategory
        } = this.state;

        this.setState({
            category: savedCategory,
            isSaved: false
        });
    }

    render() {
        const {
            recurringTransaction,
            currency,
            exitEditMode
        } = this.props;

        return (
                <form className="recurringTransactionEdit" onSubmit={ this.onSubmit }>
                    <div className="inputWithLabel">
                        <label htmlFor="name">Name: </label>
                        <SuggestedInput
                            type="text"
                            name="name"
                            autoComplete="off"
                            value={ this.state.name }
                            onChange={ this.handleChange('name') }
                            accountId={ recurringTransaction.getIn(['transaction', 'accountId']) }
                            field="name"
                            selectSuggestion={ this.selectNameSuggestion }
                            hoverSuggestion={ this.hoverNameSuggestion }
                            unhoverSuggestion={ this.unhoverNameSuggestion }
                        />
                    </div>
                    <div className="inputWithLabel">
                        <label htmlFor="date">Next Occurs: </label>
                        <input
                            type="date"
                            name="date"
                            value={ toRFC3339(this.state.date) }
                            onChange={ this.handleChange('date') }
                        />
                    </div>
                    <div className="inputWithLabel">
                        <label htmlFor="category">Category: </label>
                        <SuggestedInput
                            type="text"
                            name="category"
                            autoComplete="off"
                            value={ this.state.category }
                            onChange={ this.handleChange('category') }
                            accountId={ recurringTransaction.getIn(['transaction', 'accountId']) }
                            field="category"
                            selectSuggestion={ this.selectCategorySuggestion }
                            hoverSuggestion={ this.hoverCategorySuggestion }
                            unhoverSuggestion={ this.unhoverCategorySuggestion }
                        />
                    </div>
                    <div className="inputWithLabel">
                        <label htmlFor="amount">Amount: </label>
                        <input
                            type="number"
                            step="0.01"
                            name="amount"
                            value={ this.state.amount / Math.pow(10, currency.get('digitsAfterDecimal')) }
                            onChange={ this.handleChange('amount') }
                        />
                    </div>
                    <div>
                        <button type="button" onClick={ exitEditMode }>Cancel</button>
                        <button type="submit">Save</button>
                    </div>
                    <div className="inputWithLabel scheduleTypeInputWithLabel">
                        <label htmlFor="fixedInterval">Schedule Type: </label>
                        <radiogroup>
                            <input
                                type="radio"
                                name="fixedDate"
                                checked={ this.state.scheduleType === SCHEDULE_TYPES.FIXED_DAY_WEEK }
                                value={ SCHEDULE_TYPES.FIXED_DAY_WEEK }
                                onChange={ this.handleChange('scheduleType') }
                            />Weekly
                            <input
                                type="radio"
                                name="fixedDate"
                                checked={ this.state.scheduleType === SCHEDULE_TYPES.FIXED_DAY_MONTH }
                                value={ SCHEDULE_TYPES.FIXED_DAY_MONTH }
                                onChange={ this.handleChange('scheduleType') }
                            />Monthly
                            <input
                                type="radio"
                                name="fixedDate"
                                checked={ this.state.scheduleType === SCHEDULE_TYPES.FIXED_DAY_YEAR }
                                value={ SCHEDULE_TYPES.FIXED_DAY_YEAR }
                                onChange={ this.handleChange('scheduleType') }
                            />Yearly
                            <input
                                type="radio"
                                name="fixedInterval"
                                checked={ this.state.scheduleType === SCHEDULE_TYPES.FIXED_INTERVAL }
                                value={ SCHEDULE_TYPES.FIXED_INTERVAL }
                                onChange={ this.handleChange('scheduleType') }
                            />Fixed Interval
                        </radiogroup>
                    </div>
                    {
                    recurringTransaction.get('id') !== NEW_RECURRING_TRANSACTION_ID
                    ? (
                    <div>
                        <button type="button" className="deleteButton" onClick={ this.deleteRecurringTransaction }>Delete</button>
                    </div>
                    ) : null
                    }
                    <div className="inputWithLabel scheduleDetailsInputWithLabel">
                        {(() => {
                        switch (this.state.scheduleType) {
                            case SCHEDULE_TYPES.FIXED_DAY_WEEK:
                            return (
                            <div>
                                <label htmlFor="dayOf">This transaction will repeat every week on: </label>
                                <select
                                    name="dayOf"
                                    onChange={ this.handleChange('dayOf') }
                                    value={ this.state.dayOf }
                                >
                                    <option value={ 0 }>Sunday</option>
                                    <option value={ 1 }>Monday</option>
                                    <option value={ 2 }>Tuesday</option>
                                    <option value={ 3 }>Wednesday</option>
                                    <option value={ 4 }>Thursday</option>
                                    <option value={ 5 }>Friday</option>
                                    <option value={ 6 }>Saturday</option>
                                </select>
                            </div>
                            );

                            case SCHEDULE_TYPES.FIXED_DAY_MONTH:
                            return (
                            <div>
                                <label htmlFor="dayOf">This transaction will repeat every month on day: </label>
                                <select
                                    name="dayOf"
                                    onChange={ this.handleChange('dayOf') }
                                    value={ this.state.dayOf }
                                >
                                    {
                                    [...Array(31).keys()].map(n => <option key={ n } value={ n + 1 }>{ n + 1 }</option>)
                                    }
                                </select>
                            </div>
                            );

                            case SCHEDULE_TYPES.FIXED_DAY_YEAR:
                            return (
                            <div className="inputWithLabel">
                                <label htmlFor="dayOf">This transaction will repeat every year on day: </label>
                                <input
                                    type="number"
                                    name="dayOf"
                                    value={ this.state.dayOf }
                                    onChange={ this.handleChange('dayOf') }
                                />
                            </div>
                            );

                            case SCHEDULE_TYPES.FIXED_INTERVAL:
                            return (
                            <div className="inputWithLabel fixedIntervalInputWithLabel">
                                <label htmlFor="secondsBetween">This transaction will repeat every </label>
                                <input
                                    type="number"
                                    name="secondsBetween"
                                    value={ this.state.secondsBetween / (60 * 60 * 24) }
                                    onChange={ this.handleChange('secondsBetween') }
                                /> days
                            </div>
                            );

                            default:
                                return null;
                        }
                        })()
                        }
                    </div>
                    <div className="inputWithLabel secondsBeforeToPostInputWithLabel">
                        <label htmlFor="secondsBeforeToPost">The transaction will be posted </label>
                        <input
                            type="number"
                            name="secondsBeforeToPost"
                            value={ this.state.secondsBeforeToPost / (60 * 60 * 24) }
                            onChange={ this.handleChange('secondsBeforeToPost') }
                        /> days in advance
                    </div>
                </form>
                );
    }
}

export default connect()(RecurringTransactionEdit);
