import PropTypes from 'prop-types';
import React from 'react';
import { connect } from 'react-redux';

import { login } from '../../actions';
import styles from './googleLoginButton.css';

function attachToButton(dispatcher, callback) {
    return function() {
        gapi.load('auth2', function(){
            let auth2 = gapi.auth2.init({
                client_id: '900762983843-0ih1hv6b4mf4ql847ini51hhfc4svqoc.apps.googleusercontent.com',
                scope: 'email'
            });
            let btnEl = document.getElementById('googleBtn');
            if (btnEl) {
                auth2.attachClickHandler(btnEl, {},
          function(googleUser) {
              dispatcher(login(googleUser, callback));
          });
            }
        });
    };
}

@connect()
export default class GoogleLoginButton extends React.Component {
    static propTypes = {
        dispatch: PropTypes.func.isRequired,
        onLogin: PropTypes.func
    }

    componentDidMount() {
        const {
      dispatch,
      onLogin
    } = this.props;

        if (document.readyState === 'complete') {
            attachToButton(dispatch, onLogin)();
        } else {
            window.onload = attachToButton(dispatch, onLogin);
        }
    }

    componentWillUnmount() {
        window.onload = undefined;
    }

    render() {
        return (
      <div id="googleBtn" className={ styles.loginButton } >
        Login
      </div>
        );
    }
}
