import React from 'react';
import { Router, IndexRoute, Route, browserHistory } from 'react-router';
import thunkMiddleware from 'redux-thunk';
import { render } from 'react-dom';
import { Provider } from 'react-redux';
import { createStore, applyMiddleware } from 'redux';
import { fetchUser } from './actions';
import App from './components/index';
import Home from './components/home';
import Recurring from './components/recurring';
import AccountsPage from './components/accountsPage';
import reducers from './reducers';

const rootElement = document.getElementById('app')

let store = createStore(
  reducers,
  applyMiddleware(
    thunkMiddleware
  )
);

function fetchAuth(nextState, replace, callback) {
  if (!store.getState().auth.get('fetched')) {
    store.dispatch(fetchUser(callback));
  } else {
    callback();
  }
}

function checkAuth(nextState, replace, callback) {
  if (!store.getState().auth.get('authd')) {
    replace({
      pathname: '/'
    });
  }
  callback();
}

function goToAccounts(nextState, replace, callback) {
  if (store.getState().auth.get('authd')) {
    replace({
      pathname: '/accounts'
    });
  }
  callback();
}

render(
  <Provider store={ store }>
    <Router history={ browserHistory }>
      <Route path="/" component={ App } onEnter={ fetchAuth }>
        <IndexRoute component={ Home } onEnter={ goToAccounts }/>
        <Route path="accounts" component={ AccountsPage } onEnter={ checkAuth }/>
        <Route path="recurring" component={ Recurring } onEnter={ checkAuth }/>
      </Route>
    </Router>
  </Provider>,
  rootElement
)
