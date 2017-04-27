import PropTypes from 'prop-types';
import React from 'react';
import { connect } from 'react-redux';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { fetchAccounts, fetchCurrencies, importData } from '../../actions';
import AccountList from '../accountList';
import AccountForm from '../accountForm';
import Loader from '../loader';
import styles from './accountsPage.css';

@connect((state) => {
    return {
        accounts: state.accounts,
        accountTransaction: state.accountTransaction,
        currencies: state.currencies
    };
})
class AccountsPage extends React.Component {
    static propTypes = {
        accounts: ImmutablePropTypes.map.isRequired,
        accountTransaction: ImmutablePropTypes.map.isRequired,
        currencies: ImmutablePropTypes.map.isRequired,
        dispatch: PropTypes.func.isRequired,
        children: PropTypes.element.isRequired
    }

    importButton = () => {
        this.props.dispatch(importData());
    }

    render () {
        const {
            accounts,
            currencies,
            children
        } = this.props;

        let selected = accounts.get('selected');
        let currency;
        if (selected !== -1) {
            let currencyCode = accounts.get('accounts').get(selected).get('currency');
            currency = currencies.get('currencies').get(currencyCode);
        }

        return (
            <div className={ styles.accountsPage }>
                <div className={ styles.accountList }>
                    <AccountList selected={ selected } />
                </div>
                {
                    selected !== -1 ? (
                        <div className={ styles.transactionList }>
                            {
                                React.cloneElement(children, {
                                    accountId: selected,
                                    currency: currency
                                })
                            }
                        </div>
                    ) : <AccountForm />
                }
            </div>
        );
    }
}

@connect((state) => {
    return {
        accounts: state.accounts,
        currencies: state.currencies
    };
})
export default class AccountsPageWrapper extends React.Component {
    static propTypes = {
        accounts: ImmutablePropTypes.map.isRequired,
        currencies: ImmutablePropTypes.map.isRequired,
        children: PropTypes.element.isRequired,
        dispatch: PropTypes.func.isRequired
    }

    constructor (props) {
        super(props);
        props.dispatch(fetchAccounts());
        props.dispatch(fetchCurrencies());
    }

    render () {
        const {
            accounts,
            currencies,
            children
        } = this.props;

        return (
            <Loader loading={ !accounts.get('fetched') || !currencies.get('fetched') }>
                <AccountsPage children={ children }/>
            </Loader>
        );
    }
}
