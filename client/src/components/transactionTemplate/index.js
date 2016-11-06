import React from 'react';
import classNames from 'classnames';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { connect } from 'react-redux';
import Autosuggest from 'react-autosuggest';
import styles from './transactionTemplate.css';
import { toCurrency, toDecimal, toWhole, queryByFieldAndVal } from '../../utils';
import { putTransactionTemplate } from '../../actions';

export class TransactionTemplate extends React.Component {
  static propTypes = {
    transactionTemplate: React.PropTypes.object,
    currency: ImmutablePropTypes.map.isRequired,
  };

  constructor(props) {
    super(props);
    this.state = { editMode: false };
  }

  enterEditMode = () => {
    this.setState({ editMode: true });
  }

  exitEditMode = () => {
    this.setState({ editMode: false });
  }

  save = () => {
    this.exitEditMode();
  }

  render () {
    const {
      transactionTemplate,
      currency
    } = this.props;

    return transactionTemplate ? (
      this.state.editMode ? (
        <TransactionTemplateForm transactionTemplate={ transactionTemplate } initialValues={ getFormInitialValues(transactionTemplate, currency) } done={ this.exitEditMode } currency={ currency } />
      ) : (
        <div className={ classNames(styles.transactionTemplate, styles.transactionTemplateFields) }>
          <span className={ classNames(styles.transactionTemplateField, styles.nonEdit) } onClick={ this.enterEditMode }>{ transactionTemplate.get('templateName') }</span>
          <span className={ classNames(styles.transactionTemplateField, styles.nonEdit) } onClick={ this.enterEditMode }>{ transactionTemplate.get('name') }</span>
          <span className={ classNames(styles.transactionTemplateField, styles.nonEdit) } onClick={ this.enterEditMode }>{ transactionTemplate.get('category') }</span>
          <span className={ classNames(styles.transactionTemplateField, styles.nonEdit) } onClick={ this.enterEditMode }>{ toCurrency(toDecimal(transactionTemplate.get('amount'), currency.get('digitsAfterDecimal')), currency.get('code')) }</span>
        </div>
      )
    ) : null;
  }
}

function getFormInitialValues(transactionTemplate, currency) {
  return {
    templateName: transactionTemplate.get('templateName'),
    name: transactionTemplate.get('name'),
    category: transactionTemplate.get('category'),
    amount: toDecimal(transactionTemplate.get('amount'), currency.get('digitsAfterDecimal'))
  };
}

function getSuggestionValue(field, suggestion) {
  return suggestion[field];
}

function renderSuggestion(field, suggestion) {
  return (
    <span>{ suggestion[field] }</span>
  );
}

@connect()
export class TransactionTemplateForm extends React.Component {
  static propTypes = {
    currency: ImmutablePropTypes.map.isRequired,
    dispatch: React.PropTypes.func.isRequired,
    initialValues: React.PropTypes.object.isRequired,
    // either transactionTemplate (for editing) or accountId (for new transactionTemplates) should be passed
    accountId: React.PropTypes.number,
    transactionTemplate: ImmutablePropTypes.map,
    done: React.PropTypes.func
  };

  constructor(props) {
    super(props);

    this.state = {
      values: props.initialValues,
      suggestions: {
        name: [],
        category: []
      }
    };

    this.lastRequestId = null;
  }

  componentWillReceiveProps(nextProps) {
    this.setState({
      values: nextProps.initialValues
    });
  }

  loadSuggestions = (field, value) => {
    const {
      accountId,
      transactionTemplate
    } = this.props;

    let resolvedAccountId = accountId;
    if (transactionTemplate) {
      resolvedAccountId = transactionTemplate.get('accountId');
    }

    let id = Math.random();
    this.setState({
      lastRequestId: id
    });

    let that = this;

    // ideally requests are made from actions, buuuuut it is much easier and faster to skip redux
    queryByFieldAndVal(resolvedAccountId, field, value).then(transactionTemplates => {
      // stale query
      if (id !== that.state.lastRequestId) {
        return;
      }

      let newState = { suggestions: {} };
      Object.assign(newState.suggestions, that.state.suggestions);
      newState.suggestions[field] = transactionTemplates;
      that.setState(newState);
    });
  }

  onChange = field => {
    let that = this;
    return function(event, { newValue }) {
      let values = Object.assign({}, that.state.values);
      values[field] = newValue;
      that.setState({
        values
      });
    };
  };

  onSuggestionsFetchRequested = (field, { value }) => {
    this.loadSuggestions(field, value);
  };

  onSuggestionsClearRequested = () => {
    this.setState({
      suggestions: {
        name: [],
        category: []
      }
    });
  };

  onNameSuggestionSelected = (event, { suggestion }) => {
    const { currency } = this.props;

    this.setState({
      values: {
        name: suggestion.name,
        category: suggestion.category,
        amount: toDecimal(suggestion.amount, currency.get('digitsAfterDecimal'))
      }
    });
  }

  fieldChange = name => e => {
    let newState = Object.assign({}, this.state.values);
    newState[name] = e.target.value;
    this.setState({
      values: newState
    });
  }

  submit = (e) => {
    const {
      accountId,
      currency,
      dispatch,
      done,
      transactionTemplate
    } = this.props;

    let obj = {
      templateName: e.target['templateName'].value,
      name: e.target['name'].value,
      category: e.target['category'].value,
      amount: toWhole(parseFloat(e.target['amount'].value), currency.get('digitsAfterDecimal')),
      accountId
    };

    if (transactionTemplate) {
      obj = Object.assign(transactionTemplate.toObject(), obj);
      obj.accountId = transactionTemplate.get('accountId');
    }

    dispatch(putTransactionTemplate(obj));
    done && done();
    e.preventDefault();
  }

  render () {
    const {
      suggestions,
      values
    } = this.state;

    const nameInputProps = {
      name: 'name',
      placeholder: 'Name',
      value: values.name,
      onChange: this.onChange('name')
    };

    const categoryInputProps = {
      name: 'category',
      placeholder: 'Category',
      value: values.category,
      onChange: this.onChange('category')
    };

    return (
      <div className={ styles.transactionTemplate }>
        <form onSubmit={ this.submit }>
          <div className={ styles.transactionTemplateFields }>
            <input type="text" name="templateName" value={ values.templateName } onChange={ this.fieldChange('templateName') } placeholder="Template Name" className={ styles.transactionTemplateField } />
            <Autosuggest
              id="name"
              suggestions={ suggestions.name }
              onSuggestionsFetchRequested={ this.onSuggestionsFetchRequested.bind(this, 'name') }
              onSuggestionsClearRequested={ this.onSuggestionsClearRequested }
              onSuggestionSelected={ this.onNameSuggestionSelected }
              getSuggestionValue={ getSuggestionValue.bind(undefined, 'name') }
              renderSuggestion={ renderSuggestion.bind(undefined, 'name') }
              inputProps={ nameInputProps }
              theme={ styles } />
            <Autosuggest
              id="category"
              suggestions={ suggestions.category }
              onSuggestionsFetchRequested={ this.onSuggestionsFetchRequested.bind(this, 'category') }
              onSuggestionsClearRequested={ this.onSuggestionsClearRequested }
              getSuggestionValue={ getSuggestionValue.bind(undefined, 'category') }
              renderSuggestion={ renderSuggestion.bind(undefined, 'category') }
              inputProps={ categoryInputProps }
              theme={ styles } />
            <input type="text" name="amount" value={ values.amount } onChange={ this.fieldChange('amount') } placeholder="0" className={ styles.transactionTemplateField } />
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
