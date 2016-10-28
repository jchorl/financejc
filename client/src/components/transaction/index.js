import React from 'react';
import classNames from 'classnames';
import { reduxForm } from 'redux-form';
import styles from './transaction.css';
import { toCurrency, toDate, toDecimal, toWhole } from '../../utils';
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
          <TransactionForm form={ transaction.get('id').toString() } transaction={ transaction } done={ this.exitEditMode } currency={ currency } />
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

function pad(n) {
  return n<10 ? '0'+n : n
}

function toRFC3339(d) {
  return d.getUTCFullYear() + '-'
    + pad(d.getUTCMonth() + 1) + '-'
    + pad(d.getUTCDate());
}

@reduxForm({
  fields: [
    'name',
    'date',
    'category',
    'amount'
  ]
},
  (state, props) => {
    let {
      transaction,
      currency
    } = props;

    if (props.transaction) {
      return {
        initialValues: {
          name: transaction.get('name'),
          date: toRFC3339(transaction.get('date')),
          category: transaction.get('category'),
          amount: toDecimal(transaction.get('amount'), currency.get('digitsAfterDecimal'))
        }
      };
    }
    return {
      initialValues: {
        name: '',
        date: toRFC3339(new Date()),
        category: '',
        amount: 0
      }
    };
  })
export class TransactionForm extends React.Component {
  static propTypes = {
    currency: React.PropTypes.object.isRequired,
    dispatch: React.PropTypes.func.isRequired,
    fields: React.PropTypes.object.isRequired,
    // either transaction (for editing) or accountId (for new transactions) should be passed
    transaction: React.PropTypes.object,
    accountId: React.PropTypes.number,
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
    }

    obj.date = new Date(obj.date);
    obj.accountId = accountId;
    obj.amount = newAmount;
    dispatch(putTransaction(obj, difference));
    done && done();
  }

  render () {
    const {
      fields: {
        name,
        date,
        category,
        amount
      },
      handleSubmit
    } = this.props;

    return (
      <div className={ styles.transaction }>
        <form onSubmit={ handleSubmit(this.submit) }>
          <div className={ styles.transactionFields }>
            <input type="text" placeholder="Name" className={ styles.transactionField } { ...name }/>
            <input type="date" placeholder={ toRFC3339(new Date()) } className={ styles.transactionField } { ...date }/>
            <input type="text" placeholder="Category" className={ styles.transactionField } { ...category }/>
            <input type="text" placeholder="0" className={ styles.transactionField } { ...amount }/>
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
