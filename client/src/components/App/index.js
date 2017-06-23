import React, { Component } from 'react';
import { Route, withRouter } from 'react-router-dom';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';
import { fetchUser } from '../../actions/user';
import Nav from '../Nav';
import Accounts from '../Accounts';
import Login from '../Login';
import Search from '../Search';
import './App.css';

class App extends Component {
    static propTypes = {
        dispatch: PropTypes.func.isRequired,
        location: PropTypes.shape({
            pathname: PropTypes.string.isRequired
        }).isRequired
    }

    componentWillMount() {
        const {
            dispatch,
            history,
            location,
            user
        } = this.props;

        if (!user.get('fetched')) {
            dispatch(fetchUser())
                .then(u => {
                    if (u && location.pathname === '/') {
                        history.push('/accounts');
                    } else if (!u && location.pathname !== '/') {
                        history.push('/');
                    }
                });
        }
    }

    render() {
        const {
            user
        } = this.props;

        return (
                <div id="app">
                    <Nav />
                    {
                    user.get('fetched')
                    ? (
                    <div className="content">
                        <Route exact path="/" component={ Login } />
                        <Route path="/accounts" component={ Accounts } />
                        <Route path="/search" component={ Search } />
                    </div>
                    )
                    : null
                    }
                </div>
                );
    }
}

// need withRouter so that url updates cause a rerender
// see: https://reacttraining.com/react-router/web/guides/redux-integration
export default withRouter(connect(state => {
    return {
        user: state.user
    }
})(App));
