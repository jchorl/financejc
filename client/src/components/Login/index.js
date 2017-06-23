import React, { Component } from 'react';
import { Redirect } from 'react-router-dom';
import PropTypes from 'prop-types';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { connect } from 'react-redux';
import GoogleLoginButton from './GoogleLoginButton';
import { UNAUTHD_ID } from '../../constants';
import './Login.css';

class Login extends Component {
    static propTypes = {
        user: ImmutablePropTypes.contains({
            fetched: PropTypes.bool.isRequired
        })
    }

    render() {
        const { user } = this.props;

        return user.get('id') === UNAUTHD_ID
            ? (
                    <div className="login container">
                        <h1>Login</h1>
                        <div className="divider"></div>
                        <GoogleLoginButton />
                    </div>
                    )
            : (
                    <Redirect to="/accounts" />
                    );
    }
}

export default connect(state => {
    return {
        user: state.user
    };
})(Login);
