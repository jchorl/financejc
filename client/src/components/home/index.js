import React from 'react';
import { connect } from 'react-redux';
import ImmutablePropTypes from 'react-immutable-proptypes';
import AccountsPage from '../accountsPage';
import Loader from '../loader';

@connect((state) => {
  return { auth: state.auth }
})
export default class Home extends React.Component {
  static propTypes = {
    auth: ImmutablePropTypes.map.isRequired,
  }

  render () {
    const {
      auth
    } = this.props;

    return (
      <Loader loading={ !auth.get('fetched') }>
        { auth.get('authd') ? <AccountsPage /> : <h1>Welcome to FinanceJC</h1> }
      </Loader>
    )
  }
}

