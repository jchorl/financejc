/*global gapi*/
import React, { Component } from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';
import { googleLogin } from '../../actions/user';
import './googleLoginClient';

class GoogleLoginButton extends Component {
    static propTypes = {
        dispatch: PropTypes.func.isRequired
    }

    componentDidMount() {
        const { dispatch } = this.props;

        gapi.load('auth2', function(){
            const auth2 = gapi.auth2.init({
                client_id: '900762983843-0ih1hv6b4mf4ql847ini51hhfc4svqoc.apps.googleusercontent.com',
                cookiepolicy: 'single_host_origin',
                scope: 'email'
            });
            let element = document.getElementById('googleLoginButton');
            auth2.attachClickHandler(element, {},
                    function(googleUser) {
                        dispatch(googleLogin(googleUser));
                    }, function(error) {
                        console.error(JSON.stringify(error, undefined, 2));
                    });
        });
    }

    render() {
        return (
                <button id="googleLoginButton" className="loginButton">
                    <i className="fa fa-google"></i><span className="service">Google</span>
                </button>
                );
    }
}

export default connect()(GoogleLoginButton);
