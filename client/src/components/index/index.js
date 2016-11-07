import React from 'react';
import { connect } from 'react-redux';
import { Link, withRouter } from 'react-router';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { logout, importData } from '../../actions';
import GoogleLoginButton from '../googleLoginButton';
import styles from './index.css';

@withRouter
@connect((state) => {
  return { auth: state.auth };
})
export default class App extends React.Component {
  static propTypes = {
    auth: ImmutablePropTypes.map.isRequired,
    dispatch: React.PropTypes.func.isRequired,
    children: React.PropTypes.element.isRequired,
    router: React.PropTypes.object.isRequired
  }

  fileChange = e => {
    const {
      dispatch,
    } = this.props;

    dispatch(importData(e.target.files));
  }

  dispatchLogout = () => {
    const {
      dispatch,
      router: { push }
    } = this.props;

    dispatch(logout(function() {
      push('/');
    }));
  }

  goToTransactions = () => {
    const {
      router: { push }
    } = this.props;

    push('/transactions');
  }

  render () {
    const {
      auth
    } = this.props;

    return (
      <div>
        <nav className={ styles.navBar }>
          <div className={ styles.logoAndNavOptions }>
            <div className={ styles.logo }><Link to={ '/' } className={ styles.unstyledLink }>FinanceJC</Link></div>
            <div className={ styles.navOption }><Link to={ '/transactions' } className={ styles.unstyledLink }>Transactions</Link></div>
            <div className={ styles.navOption }><Link to={ '/transactions/recurring' } className={ styles.unstyledLink }>Recurring Transactions</Link></div>
            <div className={ styles.navOption }><Link to={ '/transactions/templates' } className={ styles.unstyledLink }>Transaction Templates</Link></div>
          </div>
          <div className={ styles.options }>
            { auth.get('authd') ? (
              <div className={ styles.dropdown }>
                { auth.get('user').get('email') } <i className={ 'fa fa-chevron-down ' + styles.dropdownChevron } aria-hidden="true"></i>
                <div className={ styles.dropdownContent }>
                  <div className={ styles.dropdownOption }>Import <input type="file" onChange={ this.fileChange } /></div>
                  <div className={ styles.dropdownOption } onClick={ this.dispatchLogout }>Logout</div>
                </div>
              </div>
            ) : <GoogleLoginButton onLogin={ this.goToTransactions }/> }
          </div>
        </nav>
        { this.props.children }
      </div>
    );
  }
}
