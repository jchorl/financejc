import React from 'react';
import { connect } from 'react-redux';
import { Link } from 'react-router';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { fetchUser, logout } from '../../actions';
import GoogleLoginButton from '../googleLoginButton';
import styles from './index.css';

@connect((state) => {
  return { auth: state.auth }
})
export default class App extends React.Component {
  static propTypes = {
    auth: ImmutablePropTypes.map.isRequired,
    dispatch: React.PropTypes.func.isRequired,
    children: React.PropTypes.element.isRequired
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
        <nav className={ styles.navBar }>
          <div className={ styles.logo }>
            <Link to={ '/' } className={ styles.unstyledLink }>FinanceJC</Link>
          </div>
          <div className={ styles.options }>
            { auth.get('authd') ? (
              <div className={ styles.dropdown }>
                { auth.get('user').get('email') } <i className={ 'fa fa-chevron-down ' + styles.dropdownChevron } aria-hidden="true"></i>
                <div className={ styles.dropdownContent }>
                  <div className={ styles.dropdownOption }><Link to={ 'recurring' } className={ styles.unstyledLink }>Recurring Transactions</Link></div>
                  <div className={ styles.dropdownOption } onClick={ this.dispatchLogout }>Logout</div>
                </div>
              </div>
            ) : <GoogleLoginButton /> }
          </div>
        </nav>
        { this.props.children }
      </div>
    )
  }
}
