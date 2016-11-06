import React from 'react';
import { connect } from 'react-redux';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { newAccount } from '../../actions';
// import styles from './accountForm.css';

@connect((state) => {
  return {
    currencies: state.currencies
  };
})
export default class AccountForm extends React.Component {
  static propTypes = {
    currencies: ImmutablePropTypes.map.isRequired,
    dispatch: React.PropTypes.func.isRequired
  };

  onSubmit = (e) => {
    const {
      dispatch
    } = this.props;

    let data = {
      name: e.target['name'].value,
      currency: e.target['currency'].value
    };

    dispatch(newAccount(data));
    e.preventDefault();
  }

  render () {
    const {
      currencies
    } = this.props;

    return (
      <form onSubmit={ this.onSubmit }>
        <div>
          <label htmlFor="name">Name</label>
          <input name="name" type="text"/>
        </div>
        <div>
          <label htmlFor="currency">Currency</label>
          <select name="currency">
            <option></option>
            { currencies.get('currencies').map(currency => (<option key={ currency.get('code') } value={ currency.get('code') }>{ currency.get('code') }: { currency.get('name') }</option>)).toOrderedSet().toArray() }
          </select>
        </div>
        <button type="submit">Save</button>
      </form>
    );
  }
}
