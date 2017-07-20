import React, { Component } from 'react';
import PropTypes from 'prop-types';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { connect } from 'react-redux';
import { newTransaction, editTransaction, deleteTransaction } from '../../actions/accountData';
import { fromRFC3339, toRFC3339 } from '../../util';
import { NEW_TRANSACTION_ID } from '../../constants';
import SuggestedInput from '../SuggestedInput';
import './TransactionEdit.css';

class TransactionEdit extends Component {
    static propTypes = {
        transaction: ImmutablePropTypes.contains({
            id: PropTypes.number.isRequired,
            name: PropTypes.string.isRequired,
            date: PropTypes.instanceOf(Date),
            category: PropTypes.string.isRequired,
            amount: PropTypes.number.isRequired,
            note: PropTypes.string.isRequired,
            accountId: PropTypes.number.isRequired
        }).isRequired,
        currency: ImmutablePropTypes.contains({
            code: PropTypes.string.isRequired,
            digitsAfterDecimal: PropTypes.number.isRequired
        }).isRequired,
        dispatch: PropTypes.func.isRequired,
        exitEditMode: PropTypes.func.isRequired
    }

    constructor(props) {
        super(props);

        const {
            transaction
        } = this.props;

        this.state = {
            name: transaction.get('name'),
            date: transaction.get('date'),
            category: transaction.get('category'),
            amount: transaction.get('amount'),
            note: transaction.get('note'),
            savedName: '',
            savedCategory: '',
            savedAmount: 0,
            savedNote: '',
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

            default:
                newState[field] = e.target.value;
        }

        this.setState(newState);
    }

    deleteTransaction = () => {
        const {
            dispatch,
            transaction,
        } = this.props;

        dispatch(deleteTransaction(transaction.toJS()));
    }

    onSubmit = e => {
        const {
            dispatch,
            exitEditMode,
            transaction: originalTransaction
        } = this.props;

        const {
            name,
            date,
            category,
            note,
            amount
        } = this.state;

        const transaction = {
            name,
            date,
            category,
            note,
            amount: Math.round(amount),
            accountId: originalTransaction.get('accountId')
        };

        if (originalTransaction.get('id') === NEW_TRANSACTION_ID) {
            const amountDifference = amount;
            dispatch(newTransaction(transaction, amountDifference))
                .then(exitEditMode);
        } else {
            const amountDifference = amount - originalTransaction.get('amount');
            transaction.id = originalTransaction.get('id');
            dispatch(editTransaction(transaction, amountDifference))
                .then(exitEditMode);
        }

        e.preventDefault();
    }

    selectNameSuggestion = suggestion => {
        this.setState({
            name: suggestion.get('name'),
            category: suggestion.get('category'),
            amount: suggestion.get('amount'),
            note: suggestion.get('note')
        });
    }

    hoverNameSuggestion = suggestion => {
        const {
            name,
            category,
            amount,
            note,
            isSaved
        } = this.state;

        let newState = {
            name: suggestion.get('name'),
            category: suggestion.get('category'),
            amount: suggestion.get('amount'),
            note: suggestion.get('note')
        };

        if (!isSaved) {
            newState.savedName = name;
            newState.savedCategory = category;
            newState.savedAmount = amount;
            newState.savedNote = note;
            newState.isSaved = true;
        }

        this.setState(newState);
    }

    unhoverNameSuggestion = () => {
        const {
            savedName,
            savedCategory,
            savedAmount,
            savedNote
        } = this.state;

        this.setState({
            name: savedName,
            category: savedCategory,
            amount: savedAmount,
            note: savedNote,
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
            currency,
            exitEditMode,
            transaction
        } = this.props;

        return (
                <form className="transactionEdit" onSubmit={ this.onSubmit }>
                    <div className="inputWithLabel">
                        <label htmlFor="name">Name: </label>
                        <SuggestedInput
                            type="text"
                            name="name"
                            autoComplete="off"
                            value={ this.state.name }
                            onChange={ this.handleChange('name') }
                            accountId={ transaction.get('accountId') }
                            field="name"
                            selectSuggestion={ this.selectNameSuggestion }
                            hoverSuggestion={ this.hoverNameSuggestion }
                            unhoverSuggestion={ this.unhoverNameSuggestion }
                        />
                    </div>
                    <div className="inputWithLabel">
                        <label htmlFor="date">Date: </label>
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
                            accountId={ transaction.get('accountId') }
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
                    <div className="inputWithLabel noteInputWithLabel">
                        <label htmlFor="Notes">Notes: </label>
                        <input
                            type="text"
                            name="note"
                            value={ this.state.note }
                            onChange={ this.handleChange('note') }
                        />
                    </div>
                    {
                    transaction.get('id') !== NEW_TRANSACTION_ID
                    ? (
                    <div>
                        <button type="button" className="deleteButton" onClick={ this.deleteTransaction }>Delete</button>
                    </div>
                    ) : null
                    }
                </form>
                );
    }
}

export default connect()(TransactionEdit);
