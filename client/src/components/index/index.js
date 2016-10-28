import React from 'react';
import { connect } from 'react-redux';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { fetchAuth, logout } from '../../actions';
import AccountsPage from '../accountsPage';
import GoogleLoginButton from '../googleLoginButton';
import Loader from '../loader';
import globalStyle from './globalStyle.css';
import style from './index.css';

@connect((state) => {
  return { auth: state.auth }
})
export default class App extends React.Component {
  static propTypes = {
    auth: ImmutablePropTypes.map.isRequired,
    dispatch: React.PropTypes.func.isRequired
  }

  constructor(props) {
    super(props);
    props.dispatch(fetchAuth());
  }

  dispatchLogout = () => {
    this.props.dispatch(logout());
  }

  render () {
    const {
      auth
    } = this.props;

    return (
      <div>
        <nav className={ globalStyle.navBar }>
          <div>
            FinanceJC
          </div>
          <div>
            { auth.get('authd') ? <span className={ style.logout } onClick={ this.dispatchLogout } >Logout</span> : <GoogleLoginButton /> }
          </div>
        </nav>
        <Loader loading={ !auth.get('fetched') }>
          { auth.get('authd') ? <AccountsPage /> : <h1>Welcome to FinanceJC</h1> }
        </Loader>
      </div>
    )
  }
}
