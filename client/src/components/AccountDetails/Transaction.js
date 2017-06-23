import React, { Component } from 'react';
import PropTypes from 'prop-types';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { toCurrencyString } from '../../util';

export default class Transaction extends Component {
    static propTypes = {
        transaction: ImmutablePropTypes.contains({
            name: PropTypes.string.isRequired,
            date: PropTypes.instanceOf(Date),
            category: PropTypes.string.isRequired,
            amount: PropTypes.number.isRequired,
            balance: PropTypes.number.isRequired
        }).isRequired,
        currency: ImmutablePropTypes.contains({
            code: PropTypes.string.isRequired,
            digitsAfterDecimal: PropTypes.number.isRequired
        }).isRequired,
        enterEditMode: PropTypes.func.isRequired
    }

    render() {
        const {
            transaction,
            currency,
            enterEditMode
        } = this.props;

        return (
                <div className="transaction" onClick={ enterEditMode } >
                    <div>{ transaction.get('name') }</div>
                    <div>{ transaction.get('date').toLocaleDateString() }</div>
                    <div>{ transaction.get('category') }</div>
                    <div>{ toCurrencyString(transaction.get('amount'), currency.get('code'), currency.get('digitsAfterDecimal')) }</div>
                    <div>{ toCurrencyString(transaction.get('balance'), currency.get('code'), currency.get('digitsAfterDecimal')) }</div>
                </div>
                );
    }
}
