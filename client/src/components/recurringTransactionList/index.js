import React from 'react';
import { connect } from 'react-redux';
import ImmutablePropTypes from 'react-immutable-proptypes';
import classNames from 'classnames';
import { fetchRecurringTransactions } from '../../actions';
import { toRFC3339 } from '../../utils';
import styles from './recurringTransactionList.css';
import { RecurringTransaction, RecurringTransactionForm } from '../recurringTransaction';

@connect((state) => {
  return {
    accountRecurringTransaction: state.accountRecurringTransaction,
    currencies: state.currencies
  };
})
export default class RecurringTransactionList extends React.Component {
  static propTypes = {
    accountId: React.PropTypes.number.isRequired,
    accountRecurringTransaction: ImmutablePropTypes.map.isRequired,
    currency: ImmutablePropTypes.map.isRequired,
    dispatch: React.PropTypes.func.isRequired
  }

  constructor(props) {
    super(props);

    this.state = {
      newRecurringTransaction: false
    };

    if (!props.accountRecurringTransaction.get(props.accountId).get('fetched')) {
      props.dispatch(fetchRecurringTransactions(props.accountId));
    }
  }

  startNewRecurringTransaction = () => {
    this.setState({
      newRecurringTransaction: true
    });
  }

  exitNewRecurringTransaction = () => {
    this.setState({
      newRecurringTransaction: false
    });
  }

  render () {
    const {
      accountRecurringTransaction,
      accountId,
      currency
    } = this.props;

    const recurringTransactions = accountRecurringTransaction.get(accountId).get('recurringTransactions');
    let initialValues = {
      transaction: {
        name: '',
        date: toRFC3339(new Date()),
        category: '',
        amount: '0'
      },
      scheduleType: 'fixedInterval',
      secondsBetween: 24*60*60*30,
      dayOf: 1,
      secondsBeforeToPost: 24*60*60*2
    };

    return (
      <div>
        <div className={ styles.headings }>
          <span className={ styles.column }>Name</span>
          <span className={ styles.column }>Next Occurs</span>
          <span className={ styles.column }>Category</span>
          <span className={ styles.column }>Amount</span>
        </div>
        <div className={ styles.headings }>
          <span className={ styles.column }>Schedule Type</span>
          <span className={ classNames(styles.column, styles.details) }>Schedule Details</span>
        </div>
        { !this.state.newRecurringTransaction ?
            (
              <button className={ styles.newRecurringTransaction } onClick={ this.startNewRecurringTransaction }>
                New
              </button>
            ) : (
              <RecurringTransactionForm accountId={ accountId } form='new' done={ this.exitNewRecurringTransaction } currency={ currency } initialValues={ initialValues } />
            )
        }
        { recurringTransactions.map(recurringTransaction => (<RecurringTransaction key={ recurringTransaction.get('id') } recurringTransaction={ recurringTransaction } currency={ currency }/>)).toOrderedSet().toArray() }
      </div>
    );
  }
}