import React from 'react';
import classNames from 'classnames';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { connect } from 'react-redux';
import Autosuggest from 'react-autosuggest';
import styles from './transaction.css';
import { formDateValueToDate, toCurrency, toDate, toDecimal, toRFC3339, toWhole, queryByFieldAndVal } from '../../utils';
import { deleteTransaction, putTransaction } from '../../actions';

export class Transaction extends React.Component {
    static propTypes = {
        // TODO convert to shape
        transaction: React.PropTypes.object,
        currency: ImmutablePropTypes.map.isRequired,
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
            transaction,
            currency
        } = this.props;

        return transaction ? (
            this.state.editMode ? (
                <TransactionForm transaction={ transaction } initialValues={ getFormInitialValues(transaction, currency) } done={ this.exitEditMode } currency={ currency } />
            ) : (
                <div className={ classNames(styles.transaction, styles.transactionFields) }>
                    <span className={ classNames(styles.transactionField, styles.nonEdit) } onClick={ this.enterEditMode }>{ transaction.get('name') }</span>
                    <span className={ classNames(styles.transactionField, styles.nonEdit) } onClick={ this.enterEditMode }>{ toDate(transaction.get('date')) }</span>
                    <span className={ classNames(styles.transactionField, styles.nonEdit) } onClick={ this.enterEditMode }>{ transaction.get('category') }</span>
                    <span className={ classNames(styles.transactionField, styles.nonEdit) } onClick={ this.enterEditMode }>{ toCurrency(toDecimal(transaction.get('amount'), currency.get('digitsAfterDecimal')), currency.get('code')) }</span>
                </div>
            )
        ) : null;
    }
}

function getFormInitialValues(transaction, currency) {
    return {
        name: transaction.get('name'),
        date: transaction.get('date'),
        category: transaction.get('category'),
        amount: toDecimal(transaction.get('amount'), currency.get('digitsAfterDecimal'))
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

// getting suggestions async example: http://codepen.io/moroshko/pen/EPZpev
@connect()
export class TransactionForm extends React.Component {
    static propTypes = {
        currency: ImmutablePropTypes.map.isRequired,
        dispatch: React.PropTypes.func.isRequired,
        initialValues: React.PropTypes.shape({
            name: React.PropTypes.string.isRequired,
            category: React.PropTypes.string.isRequired,
            date: React.PropTypes.instanceOf(Date),
            amount: React.PropTypes.number
        }).isRequired,
        // either transaction (for editing) or accountId (for new transactions) should be passed
        accountId: React.PropTypes.number,
        // TODO put fields of map
        transaction: ImmutablePropTypes.map,
        done: React.PropTypes.func
    };

    constructor(props) {
        super(props);

        this.state = {
            values: props.initialValues,
            suggestions: {
                name: [],
                category: []
            }
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
            transaction
        } = this.props;

        let resolvedAccountId = accountId;
        if (transaction) {
            resolvedAccountId = transaction.get('accountId');
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
            values[field] = newValue;
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
                name: suggestion.name,
                // don't update the date
                date: this.state.values.date,
                category: suggestion.category,
                amount: toDecimal(suggestion.amount, currency.get('digitsAfterDecimal'))
            }
        });
    }

    fieldChange = name => e => {
        let newState = Object.assign({}, this.state.values);
        let newValue = e.target.value;
        if (name === 'date') {
            newValue = formDateValueToDate(newValue);
        }
        newState[name] = newValue;
        this.setState({
            values: newState
        });
    }

    deleteTransaction = (e) => {
        const {
            dispatch,
            done,
            transaction
        } = this.props;

        dispatch(deleteTransaction(transaction.get('id'), transaction.get('accountId'), transaction.get('amount')));
        done && done();
        e.preventDefault();
    }

    submit = (e) => {
        const {
            accountId,
            currency,
            dispatch,
            done,
            transaction
        } = this.props;

        let obj = {
            name: e.target['name'].value,
            date: formDateValueToDate(e.target['date'].value),
            category: e.target['category'].value,
            amount: toWhole(parseFloat(e.target['amount'].value), currency.get('digitsAfterDecimal')),
            accountId
        };
        let difference = obj.amount;

        if (transaction) {
            obj = Object.assign(transaction.toObject(), obj);
            difference = difference - transaction.get('amount');
            obj.accountId = transaction.get('accountId');
        }

        dispatch(putTransaction(obj, difference));
        done && done();
        e.preventDefault();
    }

    render () {
        const {
            transaction
        } = this.props;

        const {
            suggestions,
            values
        } = this.state;

        const nameInputProps = {
            name: 'name',
            placeholder: 'Name',
            value: values.name,
            onChange: this.onChange('name')
        };

        const categoryInputProps = {
            name: 'category',
            placeholder: 'Category',
            value: values.category,
            onChange: this.onChange('category')
        };

        return (
            <div className={ styles.transaction }>
                <form onSubmit={ this.submit }>
                    <div className={ styles.transactionFields }>
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
                        <input type="date" name="date" value={ toRFC3339(values.date) } onChange={ this.fieldChange('date') } className={ styles.transactionField } />
                        <Autosuggest
                            id="category"
                            suggestions={ suggestions.category }
                            onSuggestionsFetchRequested={ this.onSuggestionsFetchRequested.bind(this, 'category') }
                            onSuggestionsClearRequested={ this.onSuggestionsClearRequested }
                            getSuggestionValue={ getSuggestionValue.bind(undefined, 'category') }
                            renderSuggestion={ renderSuggestion.bind(undefined, 'category') }
                            inputProps={ categoryInputProps }
                            theme={ styles } />
                        <input type="text" name="amount" value={ values.amount } onChange={ this.fieldChange('amount') } placeholder="0" className={ styles.transactionField } />
                    </div>
                    <div className={ styles.saveExit }>
                        <button type="button" onClick={ this.props.done }>Cancel</button>
                        {
                            transaction ? <button type="button" onClick={ this.deleteTransaction }>Delete</button> : null
                        }
                        <button type="submit">Save</button>
                    </div>
                </form>
            </div>
        );
    }
}
