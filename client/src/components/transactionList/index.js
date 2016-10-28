import React from 'react';
import { connect } from 'react-redux';
import ImmutablePropTypes from 'react-immutable-proptypes';
import styles from './transactionList.css';
import { Transaction, TransactionForm } from '../transaction';
import { toRFC3339 } from '../../utils';

@connect((state) => {
  return {
    accountTransaction: state.accountTransaction
  }
})
export default class TransactionList extends React.Component {
  static propTypes = {
    accountId: React.PropTypes.number.isRequired,
    accountTransaction: ImmutablePropTypes.map.isRequired,
    currency: React.PropTypes.object.isRequired
  };

  constructor (props) {
    super(props);
    this.state = {
      newTransaction: false
    };
  }

  startNewTransaction = () => {
    this.setState({
      newTransaction: true
    });
  }

  exitNewTransaction = () => {
    this.setState({
      newTransaction: false
    });
  }

  render () {
    const {
      accountId,
      accountTransaction,
      currency
    } = this.props;

    let transactions = accountTransaction.get(accountId).get("transactions");

    return (
      <div>
        <div className={ styles.headings }>
          <span className={ styles.column }>Name</span>
          <span className={ styles.column }>Date</span>
          <span className={ styles.column }>Category</span>
          <span className={ styles.column }>Amount</span>
        </div>
        { !this.state.newTransaction ?
            (
              <button className={ styles.newTransaction } onClick={ this.startNewTransaction }>
                New
              </button>
            ) : (
              <TransactionForm accountId={ accountId } form='new' done={ this.exitNewTransaction } currency={ currency } initialValues={ {date: toRFC3339(new Date())} } />
            )
        }
        <div>
          { transactions.map(transaction => (<Transaction key={ transaction.get('id') } transaction={ transaction } currency={ currency }/>)).toOrderedSet().toArray() }
        </div>
      </div>
    )
  }
}
