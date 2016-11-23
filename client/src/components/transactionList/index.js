import React from 'react';
import { connect } from 'react-redux';
import Immutable from 'immutable';
import ImmutablePropTypes from 'react-immutable-proptypes';
import Infinite from 'react-infinite';
import classNames from 'classnames';
import { fetchTransactions, fetchTransactionTemplates } from '../../actions';
import styles from './transactionList.css';
import { Transaction, TransactionForm } from '../transaction';
import { toRFC3339, toDecimal } from '../../utils';

function getEmptyTemplate() {
  return Immutable.Map({
    name: '',
    date: toRFC3339(new Date()),
    category: '',
    amount: 0
  });
}

@connect((state) => {
  return {
    accountTransaction: state.accountTransaction,
    accountTransactionTemplate: state.accountTransactionTemplate
  };
})
export default class TransactionList extends React.Component {
  static propTypes = {
    accountId: React.PropTypes.number.isRequired,
    accountTransaction: ImmutablePropTypes.map.isRequired,
    accountTransactionTemplate: ImmutablePropTypes.map.isRequired,
    currency: ImmutablePropTypes.map.isRequired,
    dispatch: React.PropTypes.func.isRequired
  };

  constructor (props) {
    super(props);

    this.fetchTransactionTemplatesIfNecessary(props);

    this.state = {
      newTransaction: false,
      isInfiniteLoading: false,
      values: getEmptyTemplate()
    };
  }

  componentWillReceiveProps(nextProps) {
    if (nextProps.accountTransaction.get(nextProps.accountId).get('transactions').size > this.props.accountTransaction.get(nextProps.accountId).get('transactions').size && this.state.isInfiniteLoading) {
      this.setState({ isInfiniteLoading: false });
    }
    this.fetchTransactionTemplatesIfNecessary(nextProps);
  }

  fetchTransactionTemplatesIfNecessary = (props) => {
    if (!props.accountTransactionTemplate.get(props.accountId).get('fetched')) {
      props.dispatch(fetchTransactionTemplates(props.accountId));
    }
  }

  startNewTransaction = template => () => {
    this.setState({
      newTransaction: true,
      values: template.set('date', toRFC3339(new Date())).update('amount', am => toDecimal(am, this.props.currency.get('digitsAfterDecimal')))
    });
  }

  exitNewTransaction = () => {
    this.setState({
      newTransaction: false
    });
  }

  loadNextPage = (accountId) => {
    let that = this;

    return function() {
      const { accountTransaction, dispatch } = that.props;

      let next = accountTransaction.get(accountId).get('next');
      if (that.state.isInfiniteLoading) {
        return;
      }
      that.setState({
        isInfiniteLoading: true
      });
      dispatch(fetchTransactions(accountId, next));
    };
  }

  render () {
    const {
      accountId,
      accountTransaction,
      accountTransactionTemplate,
      currency
    } = this.props;

    const {
      values
    } = this.state;

    let transactions = accountTransaction.get(accountId).get('transactions');
    let transactionTemplates = accountTransactionTemplate.get(accountId).get('transactionTemplates');

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
              <div className={ styles.newTransactionRow }>
                <div className={ styles.newBlock }>
                  New:
                </div>
                <div className={ classNames(styles.newBlock, styles.button) } onClick={ this.startNewTransaction(getEmptyTemplate()) }>
                  Empty
                </div>
                { transactionTemplates.map(transactionTemplate => (<div key={ transactionTemplate.get('id') } className={ classNames(styles.newBlock, styles.button) } onClick={ this.startNewTransaction(transactionTemplate) }>{ transactionTemplate.get('templateName') }</div>)).toOrderedSet().toArray() }
              </div>
            ) : (
              <TransactionForm accountId={ accountId } form='new' done={ this.exitNewTransaction } currency={ currency } initialValues={ values.toObject() } />
            )
        }
        <Infinite useWindowAsScrollContainer elementHeight={ 42 } onInfiniteLoad={ this.loadNextPage(accountId) } infiniteLoadBeginEdgeOffset={ 100 } >
          { transactions.map(transaction => (<Transaction key={ transaction.get('id') } transaction={ transaction } currency={ currency }/>)).toOrderedSet().toArray() }
        </Infinite>
      </div>
    );
  }
}
