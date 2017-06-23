import React, { Component } from 'react';
import PropTypes from 'prop-types';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { Route } from 'react-router-dom';
import { connect } from 'react-redux';
import SelectedAccount from '../SelectedAccount';
import { fetchAccounts } from '../../actions/account';
import { fetchCurrencies } from '../../actions/currency';
import './Accounts.css';

class Accounts extends Component {
    static propTypes = {
        account: ImmutablePropTypes.contains({
            fetched: PropTypes.bool.isRequired
        }),
        currency: ImmutablePropTypes.contains({
            fetched: PropTypes.bool.isRequired,
        }),
        dispatch: PropTypes.func.isRequired
    }

    redirectIfNecessary = _ => {
        const {
            account,
            currency,
            location,
            history
        } = this.props;

        if (location.pathname === '/accounts'
                && account.get('fetched')
                && currency.get('fetched')) {
            if (!account.get('items').isEmpty()) {
                history.push('/accounts/summary');
            } else {
                history.push('/accounts/create');
            }
        }
    }

    componentWillMount() {
        const {
            account,
            currency,
            dispatch
        } = this.props;

        let needToFetch = false;
        if (!account.get('fetched')) {
            needToFetch = true;
            dispatch(fetchAccounts())
                .then(this.redirectIfNecessary);
        }

        if (!currency.get('fetched')) {
            needToFetch = true;
            dispatch(fetchCurrencies())
                .then(this.redirectIfNecessary);
        }

        if (!needToFetch) {
            this.redirectIfNecessary();
        }
    }

    render() {
        const { currency } = this.props;

        return currency.get('fetched')
            ? (
                    <div className="accounts">
                        <Route path="/accounts/:id" component={ SelectedAccount } />
                    </div>
                    )
            : null;
    }
}

export default connect(state => ({
    account: state.account,
    currency: state.currency
}))(Accounts);
