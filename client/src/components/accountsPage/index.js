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
    dispatch: React.PropTypes.func.isRequired,
    children: React.PropTypes.element.isRequired
  }

  constructor (props) {
    super(props);
    let selected = props.accounts.get('accounts').size !== 0
      ? props.accounts.get('accounts').first().get('id')
      : -1;
    this.state = {
      selected: selected
    };

    // this.fetchTransactionsIfNecessary(selected);
  }

  componentWillReceiveProps(nextProps) {
    let nextIds = nextProps.accounts.get('accounts').keySeq().toSet();
    let currIds = this.props.accounts.get('accounts').keySeq().toSet();
    if (nextIds.size === currIds.size + 1) {
      let diff = nextIds.subtract(currIds);
      this.selectAccount(diff.first());
    }
  }

  // fetchTransactionsIfNecessary = (id) => {
  //   const { accountTransaction, dispatch } = this.props;
  //   // reload transactions if necessary
  //   if (id === -1) return

  //   // no need to fetch for new accounts
  //   if (!accountTransaction.has(id)) return
  //   if (accountTransaction.get(id).get("transactions").isEmpty()) {
  //     dispatch(fetchTransactions(id));
  //   }
  // }

  selectAccount = (id) => {
    this.setState({
      selected: id
    });

    // this.fetchTransactionsIfNecessary(id);
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
    const selected = this.state.selected;

    let currency;
    if (selected !== -1) {
      let currencyCode = accounts.get('accounts').get(selected).get('currency');
      currency = currencies.get('currencies').get(currencyCode);
    }

    return (
      <div>
        {
          accounts.get('accounts').size !== 0 ? null : (
            <div>
              Place QIF files in the /import folder and click <button onClick={ this.importButton }>Import</button>
            </div>
          )
        }
        <div className={ styles.accountsPage }>
          <div className={ styles.accountList }>
            <AccountList selected={ selected } onSelect={ this.selectAccount } />
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
    children: React.PropTypes.element.isRequired,
    dispatch: React.PropTypes.func.isRequired
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
