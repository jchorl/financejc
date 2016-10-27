import React from 'react';
import Immutable from 'immutable';
import { render } from 'react-dom';
import { connect } from 'react-redux';
import { fetchAccounts, fetchCurrencies, fetchTransactions, importData } from '../../actions';
import AccountList from '../accountList';
import TransactionList from '../transactionList';
import Loader from '../loader';
import styles from './accountsPage.css';

@connect((state) => {
  return {
    accounts: state.accounts,
    accountTransaction: state.accountTransaction,
    currencies: state.currencies
  }
})
class AccountsPage extends React.Component {
  static propTypes = {
    accounts: React.PropTypes.object.isRequired,
    accountTransaction: React.PropTypes.object.isRequired,
    currencies: React.PropTypes.object.isRequired
  }

  constructor (props) {
    super(props);
    let selected = props.accounts.get('accounts').size !== 0
      ? props.accounts.get('accounts').first().get('id')
      : -1;
    this.state = {
      selected: selected
    };

    this.fetchTransactionsIfNecessary(selected);
  }

  fetchTransactionsIfNecessary = (id) => {
    // reload transactions if necessary
    if (this.state.selected === -1) return
    if (this.props.accountTransaction.get(id).get("transactions").isEmpty()) {
      this.props.dispatch(fetchTransactions(id));
    }
  }

  selectAccount = (id) => {
    this.setState({
      selected: id
    });

    this.fetchTransactionsIfNecessary(id);
  }

  importButton = () => {
    this.props.dispatch(importData());
  }

  render () {
    const {
      accounts,
      currencies
    } = this.props;
    const selected = this.state.selected;

    let currency;
    if (selected !== -1) {
      let currencyCode = accounts.get('accounts').get(selected).get('currency');
      currency = currencies.get('currencies').get(currencyCode);
    }

    return (
      <div className={ styles.accountsPage }>
        {
          accounts.get('accounts').size !== 0 ? (
            <div className={ styles.accountList }>
              <AccountList selected={ selected } onSelect={ this.selectAccount } />
            </div>
          ) : (
            <div>
              Place QIF files in the /import folder and click <button onClick={ this.importButton }>Import</button>
            </div>
          )
        }
        {
          selected !== -1 ? (
            <div className={ styles.transactionList }>
              <TransactionList accountId={ selected } currency={ currency } />
            </div>
          ) : null
        }
      </div>
    )
  }
}

@connect((state) => {
  return {
    accounts: state.accounts,
    currencies: state.currencies
  }
})
export default class AccountsPageWrapper extends React.Component {
  static propTypes = {
    accounts: React.PropTypes.object.isRequired
  }

  constructor (props) {
    super(props);
    props.dispatch(fetchAccounts());
    props.dispatch(fetchCurrencies());
  }

  render () {
    return (
      <Loader loading={ !this.props.accounts.get('fetched') || !this.props.currencies.get('fetched') }>
        <AccountsPage />
      </Loader>
    )
  }
}
