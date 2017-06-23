import React, { Component } from 'react';
import PropTypes from 'prop-types';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { connect } from 'react-redux';
import { toCurrencyString } from '../../util';

class SummaryTransaction extends Component {
    static propTypes = {
        account: ImmutablePropTypes.contains({
            items: ImmutablePropTypes.mapOf(
                           ImmutablePropTypes.contains({
                               name: PropTypes.string.isRequired,
                               currency: PropTypes.string.isRequired
                           }).isRequired
                           ).isRequired
        }).isRequired,
        currency: ImmutablePropTypes.contains({
            items: ImmutablePropTypes.mapOf(
                           ImmutablePropTypes.contains({
                               digitsAfterDecimal: PropTypes.number.isRequired
                           }).isRequired
                           ).isRequired
        }).isRequired,
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
            account,
            currency,
            transaction
        } = this.props;

        const acc = account.getIn(['items', transaction.get('accountId')]);
        const currencyCode = acc.get('currency');
        const digitsAfterDecimal = currency.getIn(['items', currencyCode, 'digitsAfterDecimal']);

        return (
                <div className="summaryTransaction">
                    <div>{ transaction.get('name') }</div>
                    <div>{ transaction.get('date').toLocaleDateString() }</div>
                    <div>{ transaction.get('category') }</div>
                    <div>{ toCurrencyString(transaction.get('amount'), currencyCode, digitsAfterDecimal) }</div>
                    <div>{ acc.get('name') }</div>
                </div>
                );
    }
}

export default connect(state => ({
    account: state.account,
    currency: state.currency
}))(SummaryTransaction);
