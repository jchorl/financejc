import React from 'react';
import { connect } from 'react-redux';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { fetchUser, logout } from '../../actions';
import AccountsPage from '../accountsPage';
import GoogleLoginButton from '../googleLoginButton';
import Loader from '../loader';
import styles from './index.css';

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
    props.dispatch(fetchUser());
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
        <div className={ styles.navBar }>
          <div className={ styles.logo }>
            FinanceJC
          </div>
          <div className={ styles.options }>
            { auth.get('authd') ? (
              <div className={ styles.dropdown }>
                { auth.get('user').get('email') } <i className={ "fa fa-chevron-down " + styles.dropdownChevron } aria-hidden="true"></i>
                <div className={ styles.dropdownContent }>
                  <div className={ styles.dropdownOption }>Recurring Transactions</div>
                  <div className={ styles.dropdownOption } onClick={ this.dispatchLogout } > Logout</div>
                </div>
              </div>
            ) : <GoogleLoginButton /> }
          </div>
        </div>
        <Loader loading={ !auth.get('fetched') }>
          { auth.get('authd') ? <AccountsPage /> : <h1>Welcome to FinanceJC</h1> }
        </Loader>
      </div>
    )
  }
}
