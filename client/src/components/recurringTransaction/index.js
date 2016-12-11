import React from 'react';
import classNames from 'classnames';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { connect } from 'react-redux';
import Autosuggest from 'react-autosuggest';
import styles from './recurringTransaction.css';
import { toCurrency, toDate, toDecimal, toWhole, toRFC3339, queryByFieldAndVal } from '../../utils';
import { putRecurringTransaction, deleteRecurringTransaction } from '../../actions';

function getFixedType(scheduleType) {
    switch(scheduleType) {
    case 'fixedDayWeek': return 'week';
    case 'fixedDayMonth': return 'month';
    case 'fixedDayYear': return 'year';
    }
}

export class RecurringTransaction extends React.Component {
    static propTypes = {
        recurringTransaction: React.PropTypes.object,
        currency: ImmutablePropTypes.map.isRequired
    };

    constructor(props) {
        super(props);
        this.state = { editMode: false };
    }

    enterEditMode = () => {
        this.setState({ editMode: true });
    }

    exitEditMode = () => {
        this.setState({ editMode: false });
    }

    save = () => {
        this.exitEditMode();
    }

    render () {
        const {
            recurringTransaction,
            currency
        } = this.props;

        return recurringTransaction ? (
            this.state.editMode ? (
                <RecurringTransactionForm recurringTransaction={ recurringTransaction } initialValues={ getFormInitialValues(recurringTransaction, currency) } done={ this.exitEditMode } currency={ currency } />
            ) : (
                <div className={ classNames(styles.recurringTransaction, styles.recurringTransactionFields) }>
                    <span className={ classNames(styles.recurringTransactionField, styles.nonEdit) } onClick={ this.enterEditMode }>{ recurringTransaction.get('transaction').get('name') }</span>
                    <span className={ classNames(styles.recurringTransactionField, styles.nonEdit) } onClick={ this.enterEditMode }>{ toDate(recurringTransaction.get('transaction').get('date')) }</span>
                    <span className={ classNames(styles.recurringTransactionField, styles.nonEdit) } onClick={ this.enterEditMode }>{ recurringTransaction.get('transaction').get('category') }</span>
                    <span className={ classNames(styles.recurringTransactionField, styles.nonEdit) } onClick={ this.enterEditMode }>{ toCurrency(toDecimal(recurringTransaction.get('transaction').get('amount'), currency.get('digitsAfterDecimal')), currency.get('code')) }</span>
                    <span className={ classNames(styles.recurringTransactionField, styles.nonEdit) } onClick={ this.enterEditMode }>{ recurringTransaction.get('scheduleType') === 'fixedInterval' ? 'Fixed Interval' : 'Fixed Period' }</span>
                    { recurringTransaction.get('scheduleType') === 'fixedInterval' ? (
                        <div className={ classNames(styles.recurringTransactionField, styles.nonEdit, styles.details) } onClick={ this.enterEditMode }>
                            <div className={ styles.detail }>Interval: { recurringTransaction.get('secondsBetween') / (24 * 60 * 60) } days</div>
                            <div className={ styles.detail }>Days before to post: { recurringTransaction.get('secondsBeforeToPost') / (24 * 60 * 60) } days</div>
                        </div>
                    ) : (
                        <div className={ classNames(styles.recurringTransactionField, styles.nonEdit, styles.details) } onClick={ this.enterEditMode }>
                            <div className={ styles.detail }>Transaction will be generated every { getFixedType(recurringTransaction.get('scheduleType')) } on the { recurringTransaction.get('dayOf') } day</div>
                            <div className={ styles.detail }>Days before to post: { recurringTransaction.get('secondsBeforeToPost') / (24 * 60 * 60) } days</div>
                        </div>
                    ) }
                </div>
            )
        ) : null;
    }
}

function getFormInitialValues(recurringTransaction, currency) {
    return {
        transaction: {
            name: recurringTransaction.get('transaction').get('name'),
            date: toRFC3339(recurringTransaction.get('transaction').get('date')),
            category: recurringTransaction.get('transaction').get('category'),
            amount: toDecimal(recurringTransaction.get('transaction').get('amount'), currency.get('digitsAfterDecimal'))
        },
        scheduleType: recurringTransaction.get('scheduleType'),
        secondsBetween: recurringTransaction.get('secondsBetween'),
        dayOf: recurringTransaction.get('dayOf'),
        secondsBeforeToPost: recurringTransaction.get('secondsBeforeToPost')
    };
}

function getSuggestionValue(field, suggestion) {
    return suggestion[field];
}

function renderSuggestion(field, suggestion) {
    return (
        <span>{ suggestion[field] }</span>
    );
}

@connect()
export class RecurringTransactionForm extends React.Component {
    static propTypes = {
        currency: ImmutablePropTypes.map.isRequired,
        dispatch: React.PropTypes.func.isRequired,
        initialValues: React.PropTypes.object.isRequired,
        // either recurringTransaction (for editing) or accountId (for new transactions) should be passed
        accountId: React.PropTypes.number,
        recurringTransaction: ImmutablePropTypes.map,
        done: React.PropTypes.func
    };

    constructor(props) {
        super(props);

        this.state = {
            values: props.initialValues,
            suggestions: {
                name: [],
                category: []
            },
            fixedInterval: props.initialValues.scheduleType === 'fixedInterval'
        };

        this.lastRequestId = null;
    }

    componentWillReceiveProps(nextProps) {
        this.setState({
            values: nextProps.initialValues
        });
    }

    loadSuggestions = (field, value) => {
        const {
            accountId,
            recurringTransaction
        } = this.props;

        let resolvedAccountId = accountId;
        if (recurringTransaction) {
            resolvedAccountId = recurringTransaction.get('transaction').get('accountId');
        }

        let id = Math.random();
        this.setState({
            lastRequestId: id
        });

        let that = this;

        // ideally requests are made from actions, buuuuut it is much easier and faster to skip redux
        queryByFieldAndVal(resolvedAccountId, field, value).then(transactions => {
            // stale query
            if (id !== that.state.lastRequestId) {
                return;
            }

            let newState = { suggestions: {} };
            Object.assign(newState.suggestions, that.state.suggestions);
            newState.suggestions[field] = transactions;
            that.setState(newState);
        });
    }

    onChange = field => {
        let that = this;
        return function(event, { newValue }) {
            let values = Object.assign({}, that.state.values);
            values.transaction = Object.assign({}, values.transaction);
            values.transaction[field] = newValue;
            that.setState({
                values
            });
        };
    };

    onSuggestionsFetchRequested = (field, { value }) => {
        this.loadSuggestions(field, value);
    };

    onSuggestionsClearRequested = () => {
        this.setState({
            suggestions: {
                name: [],
                category: []
            }
        });
    };

    onNameSuggestionSelected = (event, { suggestion }) => {
        const { currency } = this.props;

        this.setState({
            values: {
                transaction: {
                    name: suggestion.name,
                    // don't update the date
                    date: this.state.values.transaction.date,
                    category: suggestion.category,
                    amount: toDecimal(suggestion.amount, currency.get('digitsAfterDecimal'))
                },
                scheduleType: this.state.values.scheduleType,
                secondsBetween: this.state.values.secondsBetween,
                dayOf: this.state.values.dayOf,
                secondsBeforeToPost: this.state.values.secondsBeforeToPost
            }
        });
    }

    fieldChange = name => e => {
        let newState = Object.assign({}, this.state.values);
        let value = e.target.value;
        if (name === 'secondsBetween' || name === 'secondsBeforeToPost') {
            value = value * 24 * 60 * 60;
        }
        newState[name] = value;
        this.setState({
            values: newState
        });
    }

    deleteRecurringTransaction = (e) => {
        const {
            dispatch,
            done,
            recurringTransaction
        } = this.props;

        dispatch(deleteRecurringTransaction(recurringTransaction.get('id'), recurringTransaction.get('transaction').get('accountId')));
        done && done();
        e.preventDefault();
    }

    transactionFieldChange = name => e => {
        let newState = Object.assign({}, this.state.values);
        newState.transaction = Object.assign({}, newState.transaction);
        newState.transaction[name] = e.target.value;
        this.setState({
            values: newState
        });
    }

    changeRecurring = event => {
        this.setState({
            recurring: event.target['value'] === 'true'
        });
    }

    changeSchedule = event => {
        this.setState({
            fixedInterval: event.target['value'] === 'true'
        });
    }

    submit = (e) => {
        const {
            accountId,
            currency,
            dispatch,
            done,
            recurringTransaction
        } = this.props;

        const { values } = this.state;

        e.preventDefault();

        let scheduleType = 'fixedInterval',
            dayOf: null,
            secondsBetween: null;
        if (e.target['fixedInterval'].value !== 'true') {
            dayOf = values.dayOf;
            scheduleType = values.scheduleType;
        } else {
            secondsBetween = parseInt(e.target['interval'].value) * 24 * 60 * 60;
        }

        let obj = {
            transaction: {
                name: e.target['name'].value,
                date: new Date(e.target['date'].value),
                category: e.target['category'].value,
                amount: toWhole(parseFloat(e.target['amount'].value), currency.get('digitsAfterDecimal')),
                accountId
            },
            scheduleType,
            secondsBetween,
            dayOf: parseInt(dayOf),
            secondsBeforeToPost: parseInt(e.target['daysBeforeToPost'].value) * 24 * 60 * 60
        };

        if (recurringTransaction) {
            obj = Object.assign(recurringTransaction.toObject(), obj);
            obj.transaction.accountId = recurringTransaction.get('transaction').get('accountId');
        }

        dispatch(putRecurringTransaction(obj));
        done && done();
    }

    render () {
        const {
            recurringTransaction
        } = this.props;

        const {
            fixedInterval,
            suggestions,
            values
        } = this.state;

        const nameInputProps = {
            name: 'name',
            placeholder: 'Name',
            value: values.transaction.name,
            onChange: this.onChange('name')
        };

        const categoryInputProps = {
            name: 'category',
            placeholder: 'Category',
            value: values.transaction.category,
            onChange: this.onChange('category')
        };

        return (
            <div className={ styles.recurringTransaction }>
                <form onSubmit={ this.submit }>
                    <div className={ styles.recurringTransactionFields }>
                        <Autosuggest
                            id="name"
                            suggestions={ suggestions.name }
                            onSuggestionsFetchRequested={ this.onSuggestionsFetchRequested.bind(this, 'name') }
                            onSuggestionsClearRequested={ this.onSuggestionsClearRequested }
                            onSuggestionSelected={ this.onNameSuggestionSelected }
                            getSuggestionValue={ getSuggestionValue.bind(undefined, 'name') }
                            renderSuggestion={ renderSuggestion.bind(undefined, 'name') }
                            inputProps={ nameInputProps }
                            theme={ styles } />
                        <input type="date" name="date" value={ values.transaction.date } onChange={ this.transactionFieldChange('date') } className={ styles.recurringTransactionField } />
                        <Autosuggest
                            id="category"
                            suggestions={ suggestions.category }
                            onSuggestionsFetchRequested={ this.onSuggestionsFetchRequested.bind(this, 'category') }
                            onSuggestionsClearRequested={ this.onSuggestionsClearRequested }
                            getSuggestionValue={ getSuggestionValue.bind(undefined, 'category') }
                            renderSuggestion={ renderSuggestion.bind(undefined, 'category') }
                            inputProps={ categoryInputProps }
                            theme={ styles } />
                        <input type="text" name="amount" value={ values.transaction.amount } onChange={ this.transactionFieldChange('amount') } placeholder="0" className={ styles.recurringTransactionField } />
                        <div className={ styles.recurringTransactionField }>
                            <label>
                                <input
                                    type="radio"
                                    name="fixedInterval"
                                    value="true"
                                    defaultChecked={ fixedInterval }
                                    onChange={ this.changeSchedule } />
                                Fixed interval
                            </label>
                            <label>
                                <input
                                    type="radio"
                                    name="fixedInterval"
                                    value="false"
                                    defaultChecked={ !fixedInterval }
                                    onChange={ this.changeSchedule } />
                                Fixed day
                            </label>
                        </div>
                        { fixedInterval ? (
                            <div className={ classNames(styles.recurringTransactionField, styles.details) }>
                                <div className={ styles.detail }>
                                    Interval: <input type="number" name="interval" value={ values.secondsBetween / (24 * 60 * 60) } onChange={ this.fieldChange('secondsBetween') }></input> days
                                </div>
                                <div className={ styles.detail }>
                                    Days before to post: <input type="number" name="daysBeforeToPost" value={ values.secondsBeforeToPost / (24 * 60 * 60) } onChange={ this.fieldChange('secondsBeforeToPost') }></input> days
                                </div>
                            </div>
                        ) : (
                            <div className={ classNames(styles.recurringTransactionField, styles.details) }>
                                <div className={ styles.detail }>Transaction will be generated every <select name="period" value={ values.scheduleType } onChange={ this.fieldChange('scheduleType') }>
                                        <option value="fixedDayWeek">week</option>
                                        <option value="fixedDayMonth">month</option>
                                        <option value="fixedDayYear">year</option>
                                </select> on the <input type="number" name="dayOf" onChange={ this.fieldChange('dayOf') } value={ values.dayOf } /> day</div>
                                <div className={ styles.detail }>
                                    Days before to post: <input type="number" name="daysBeforeToPost" value={ values.secondsBeforeToPost / (24 * 60 * 60) } onChange={ this.fieldChange('secondsBeforeToPost') }></input> days
                                </div>
                            </div>
                        ) }
                    </div>
                    <div className={ styles.saveExit }>
                        <button type="button" onClick={ this.props.done }>Cancel</button>
                        {
                            recurringTransaction ? <button type="button" onClick={ this.deleteRecurringTransaction }>Delete</button> : null
                        }
                        <button type="submit">Save</button>
                    </div>
                </form>
            </div>
        );
    }
}
