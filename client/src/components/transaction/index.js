import React from 'react';
import classNames from 'classnames';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { Field, reduxForm } from 'redux-form';
import styles from './transaction.css';
import { toCurrency, toDate, toDecimal, toWhole, toRFC3339 } from '../../utils';
import { putTransaction } from '../../actions';

export class Transaction extends React.Component {
  static propTypes = {
    transaction: React.PropTypes.object,
    currency: React.PropTypes.object.isRequired,
  };

  constructor(props) {
    super(props);
    this.state = {editMode: false};
  }

  enterEditMode = () => {
    this.setState({editMode: true});
  }

  exitEditMode = () => {
    this.setState({editMode: false});
  }

  save = () => {
    this.exitEditMode();
  }

  render () {
    const {
      transaction,
      currency
    } = this.props;

    return transaction ? (
      <div className={ styles.transaction }>
        { this.state.editMode ? (
          <TransactionForm form={ transaction.get('id').toString() } transaction={ transaction } initialValues={ getFormInitialValues(transaction, currency) } done={ this.exitEditMode } currency={ currency } />
        ) : (
          <div className={ styles.transactionFields }>
            <span className={ classNames(styles.transactionField, styles.nonEdit) } onClick={ this.enterEditMode }>{ transaction.get('name') }</span>
            <span className={ classNames(styles.transactionField, styles.nonEdit) } onClick={ this.enterEditMode }>{ toDate(transaction.get('date')) }</span>
            <span className={ classNames(styles.transactionField, styles.nonEdit) } onClick={ this.enterEditMode }>{ transaction.get('category') }</span>
            <span className={ classNames(styles.transactionField, styles.nonEdit) } onClick={ this.enterEditMode }>{ toCurrency(toDecimal(transaction.get('amount'), currency.get('digitsAfterDecimal')), currency.get('code')) }</span>
          </div>
        ) }
      </div>
    ) : null
  }
}

function getFormInitialValues(transaction, currency) {
  return {
    name: transaction.get('name'),
    date: toRFC3339(transaction.get('date')),
    category: transaction.get('category'),
    amount: toDecimal(transaction.get('amount'), currency.get('digitsAfterDecimal'))
  }
}

@reduxForm()
export class TransactionForm extends React.Component {
  static propTypes = {
    currency: ImmutablePropTypes.map.isRequired,
    dispatch: React.PropTypes.func.isRequired,
    handleSubmit: React.PropTypes.func.isRequired,
    initialValues: React.PropTypes.object,
    // either transaction (for editing) or accountId (for new transactions) should be passed
    accountId: React.PropTypes.number,
    transaction: ImmutablePropTypes.map,
    done: React.PropTypes.func
  };

  submit = (data) => {
    const {
      accountId,
      currency,
      dispatch,
      done,
      transaction
    } = this.props;

    let obj = data;
    let newAmount = toWhole(parseFloat(data.amount), currency.get('digitsAfterDecimal'));
    let difference = newAmount;

    if (transaction) {
      obj = Object.assign(transaction.toObject(), data);
      difference = difference - transaction.get('amount');
      obj.accountId = transaction.get('accountId');
    } else {
      obj.accountId = accountId;
    }

    obj.date = new Date(obj.date);
    obj.amount = newAmount;
    dispatch(putTransaction(obj, difference));
    done && done();
  }

  render () {
    const {
      handleSubmit
    } = this.props;

    return (
      <div className={ styles.transaction }>
        <form onSubmit={ handleSubmit(this.submit) }>
          <div className={ styles.transactionFields }>
            <Field type="text" name="name" placeholder="Name" component="input" className={ styles.transactionField } />
            <Field type="date" name="date" placeholder={ toRFC3339(new Date()) } component="input" className={ styles.transactionField } />
            <Field type="text" name="category" placeholder="Category" component="input" className={ styles.transactionField } />
            <Field type="text" name="amount" placeholder="0" component="input" className={ styles.transactionField } />
          </div>
          <div className={ styles.saveExit }>
              <button type="button" onClick={ this.props.done }>Cancel</button>
              <button type="submit">Save</button>
          </div>
        </form>
      </div>
    );
  }
}
