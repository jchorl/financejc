import React from 'react';
import { Field, reduxForm } from 'redux-form';
import { connect } from 'react-redux';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { newAccount } from '../../actions';
// import styles from './accountForm.css';

@reduxForm({
  form: "newAccount"
})
@connect((state) => {
  return {
    currencies: state.currencies
  }
})
export default class AccountForm extends React.Component {
  static propTypes = {
    handleSubmit: React.PropTypes.func.isRequired,
    currencies: ImmutablePropTypes.map.isRequired,
    dispatch: React.PropTypes.func.isRequired
  };

  onSubmit = (data) => {
    const {
      dispatch
    } = this.props;

    dispatch(newAccount(data));
  }

  render () {
    const {
      handleSubmit,
      currencies
    } = this.props;

    return (
      <form onSubmit={ handleSubmit(this.onSubmit) }>
        <div>
          <label htmlFor="name">Name</label>
          <Field name="name" component="input" type="text"/>
        </div>
        <div>
          <label htmlFor="currency">Currency</label>
          <Field name="currency" component="select">
            <option></option>
            { currencies.get('currencies').map(currency => (<option key={ currency.get('code') } value={ currency.get('code') }>{ currency.get('code') }: { currency.get('name') }</option>)).toOrderedSet().toArray() }
          </Field>
        </div>
        <button type="submit">Save</button>
      </form>
    );
  }
}
