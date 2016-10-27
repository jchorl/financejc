import React from 'react';
import { render } from 'react-dom';
import { connect } from 'react-redux';
import classNames from 'classnames';
import { toCurrency, toDecimal } from '../../utils';
import styles from './accountList.css';

@connect((state) => {
  return {
    accounts: state.accounts,
    currencies: state.currencies
  }
})
export default class AccountList extends React.Component {
  static propTypes = {
    accounts: React.PropTypes.object.isRequired,
    currencies: React.PropTypes.object.isRequired,
    onSelect: React.PropTypes.func,
    selected: React.PropTypes.number
  }

  render () {
    const {
      accounts,
      currencies,
      onSelect,
      selected
    } = this.props;

    return (
      <div>
        <h3 className={ styles.accountsTitle }>Accounts</h3>
        <div>
          { accounts.get('accounts').map(account => {
            let selectedClass = {};
            selectedClass[styles.selected] = selected === account.get('id');
            return (
              <button key={ account.get('id') } className={ classNames(styles.accountButton, selectedClass) } onClick={ onSelect.bind(this, account.get('id')) }>
                <div className={ styles.accountName }>
                  { account.get('name') }
                </div>
                <div className={ styles.accountInfo }>
                  Balance: { toCurrency(toDecimal(account.get('futureValue'), currencies.get('currencies').get(account.get('currency')).get('digitsAfterDecimal')), account.get('currency')) }
                </div>
              </button>
            )
          }).valueSeq().toArray() }
        </div>
      </div>
    )
  }
}
