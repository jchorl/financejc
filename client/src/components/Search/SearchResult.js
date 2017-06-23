import React, { Component } from 'react';
import PropTypes from 'prop-types';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { toCurrencyString } from '../../util';

export default class SearchResult extends Component {
    static propTypes = {
        accountName: PropTypes.string.isRequired,
        currency: ImmutablePropTypes.contains({
            code: PropTypes.string.isRequired,
            digitsAfterDecimal: PropTypes.number.isRequired
        }).isRequired,
        enterEditMode: PropTypes.func.isRequired,
        transaction: ImmutablePropTypes.contains({
            name: PropTypes.string.isRequired,
            date: PropTypes.instanceOf(Date),
            category: PropTypes.string.isRequired,
            amount: PropTypes.number.isRequired,
            accountId: PropTypes.number.isRequired
        }).isRequired
    }

    render() {
        const {
            accountName,
            currency,
            enterEditMode,
            transaction
        } = this.props;

        return (
                <div className="searchTransaction" onClick={ enterEditMode }>
                    <div>{ transaction.get('name') }</div>
                    <div>{ transaction.get('date').toLocaleDateString() }</div>
                    <div>{ transaction.get('category') }</div>
                    <div>{ toCurrencyString(transaction.get('amount'), currency.get('code'), currency.get('digitsAfterDecimal')) }</div>
                    <div>{ accountName }</div>
                </div>
                );
    }
}
