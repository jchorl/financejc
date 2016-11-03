import React from 'react';
import classNames from 'classnames';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { connect } from 'react-redux';
import Autosuggest from 'react-autosuggest';
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
          <TransactionForm transaction={ transaction } initialValues={ getFormInitialValues(transaction, currency) } done={ this.exitEditMode } currency={ currency } />
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

function getNameSuggestionValue(suggestion) {
  return suggestion.name;
}

function renderSuggestion(suggestion) {
  return (
    <span>{suggestion.name}</span>
  );
}

function queryByFieldAndVal(accountId, field, val) {
  return fetch(`/api/account/${accountId}/transactions/query?field=${field}&value=${val}`, {
    credentials: 'include'
  })
    .then(response => response.json());
}

@connect()
export class TransactionForm extends React.Component {
  static propTypes = {
    currency: ImmutablePropTypes.map.isRequired,
    dispatch: React.PropTypes.func.isRequired,
    initialValues: React.PropTypes.object.isRequired,
    // either transaction (for editing) or accountId (for new transactions) should be passed
    accountId: React.PropTypes.number,
    transaction: ImmutablePropTypes.map,
    done: React.PropTypes.func
  };

  constructor(props) {
    super(props);

    this.state = {
      value: props.initialValues.name,
      suggestions: [],
      isLoading: false
    };

    this.lastRequestId = null;
  }

  loadSuggestions = (field, value) => {
    const {
      accountId
    } = this.props;

    let id = Math.random();
    this.setState({
      isLoading: true,
      lastRequestId: id
    });

    let that = this;

    // ideally requests are made from actions, buuuuut it is much easier and faster to skip redux
    queryByFieldAndVal(accountId, field, value).then(transactions => {
      if (id !== that.state.lastRequestId) {
        return;
      }

      that.setState({
        isLoading: false,
        suggestions: transactions
      });
    });
  }

  onChange = (event, { newValue }) => {
    this.setState({
      value: newValue
    });
  };

  onSuggestionsFetchRequested = (field, { value }) => {
    this.loadSuggestions(field, value);
  };

  onSuggestionsClearRequested = () => {
    this.setState({
      suggestions: []
    });
  };

  submit = (e) => {
    const {
      accountId,
      currency,
      dispatch,
      done,
      transaction
    } = this.props;

    let obj = {
      name: e.target['name'].value,
      date: new Date(e.target['date'].value),
      category: e.target['category'].value,
      amount: toWhole(parseFloat(e.target['amount'].value), currency.get('digitsAfterDecimal')),
      accountId
    };
    let difference = obj.amount;

    if (transaction) {
      obj = Object.assign(transaction.toObject(), obj);
      difference = difference - transaction.get('amount');
      obj.accountId = transaction.get('accountId');
    }

    dispatch(putTransaction(obj, difference));
    done && done();
    e.preventDefault();
  }

  render () {
    const {
      value,
      suggestions
    } = this.state;

    const {
      initialValues
    } = this.props;

    const nameInputProps = {
      name: 'name',
      placeholder: 'Name',
      value,
      onChange: this.onChange
    };

    return (
      <div className={ styles.transaction }>
        <form onSubmit={ this.submit }>
          <div className={ styles.transactionFields }>
            <Autosuggest
              suggestions={ suggestions }
              onSuggestionsFetchRequested={ this.onSuggestionsFetchRequested.bind(this, 'name') }
              onSuggestionsClearRequested={ this.onSuggestionsClearRequested }
              getSuggestionValue={ getNameSuggestionValue }
              renderSuggestion={ renderSuggestion }
              inputProps={ nameInputProps }
              theme={ styles } />
            <input type="date" name="date" defaultValue={ initialValues.date } className={ styles.transactionField } />
            <input type="text" name="category" defaultValue={ initialValues.category } placeholder="Category" className={ styles.transactionField } />
            <input type="text" name="amount" defaultValue={ initialValues.amount } placeholder="0" className={ styles.transactionField } />
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
