import React, { Component } from 'react';
import PropTypes from 'prop-types';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { connect } from 'react-redux';
import SearchResult from './SearchResult';
import TransactionEdit from '../TransactionEdit';

class SearchResultWrapper extends Component {
    static propTypes = {
        account: ImmutablePropTypes.contains({
            items: ImmutablePropTypes.mapOf(
                           ImmutablePropTypes.contains({
                               name: PropTypes.string.isRequired,
                               currency: PropTypes.string.isRequired
                           }).isRequired
            ).isRequired
        }).isRequired,
        transaction: ImmutablePropTypes.map.isRequired,
        currency: ImmutablePropTypes.contains({
            items: ImmutablePropTypes.mapOf(
                           ImmutablePropTypes.contains({
                               digitsAfterDecimal: PropTypes.number.isRequired
                           }).isRequired
                           ).isRequired
        }).isRequired
    }

    constructor() {
        super();
        this.state = { isEditMode: false }
    }

    enterEditMode = () => {
        this.setState({ isEditMode: true });
    }

    exitEditMode = () => {
        this.setState({ isEditMode: false });
    }

    render() {
        const {
            account,
            transaction,
            currency
        } = this.props;

        const { isEditMode } = this.state;

        const acc = account.getIn(['items', transaction.get('accountId')]);
        const currencyCode = acc.get('currency');
        const curr = currency.getIn(['items', currencyCode]);

        return isEditMode
            ? (
                    <TransactionEdit transaction={ transaction } currency={ curr } exitEditMode={ this.exitEditMode } />
                    )
            : (
                    <SearchResult transaction={ transaction } currency={ curr } accountName={ acc.get('name') } enterEditMode={ this.enterEditMode } />
                    );
    }
}

export default connect(state => ({
    account: state.account,
    currency: state.currency
}))(SearchResultWrapper);
