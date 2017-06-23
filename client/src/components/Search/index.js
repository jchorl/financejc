import React, { Component } from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { fetchAccounts } from '../../actions/account';
import { fetchCurrencies } from '../../actions/currency';
import SearchResultWrapper from './SearchResultWrapper';
import './Search.css';

class Search extends Component {
    static propTypes = {
        account: ImmutablePropTypes.contains({
            fetched: PropTypes.bool.isRequired
        }).isRequired,
        currency: ImmutablePropTypes.contains({
            fetched: PropTypes.bool.isRequired
        }).isRequired,
        search: ImmutablePropTypes.contains({
            items: ImmutablePropTypes.list.isRequired
        }).isRequired
    }

    componentWillMount() {
        const {
            account,
            currency,
            dispatch
        } = this.props;

        if (!account.get('fetched')) {
            dispatch(fetchAccounts())
                .then(this.dataReceived);
        }

        if (!currency.get('fetched')) {
            dispatch(fetchCurrencies())
                .then(this.dataReceived);
        }
    }

    render() {
        const { search } = this.props;

        return (
                <div className="search container">
                    <h1>Search Results</h1>
                    <div className="divider"></div>
                    <div className="searchResultsList">
                        { search.get('items').map(t => <SearchResultWrapper key={ t.get('id') } transaction={ t } />) }
                    </div>
                </div>
                );
    }
}

export default connect(state => ({
    account: state.account,
    currency: state.currency,
    search: state.search
}))(Search);
