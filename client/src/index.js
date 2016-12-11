import React from 'react';
import { Router, IndexRoute, Route, browserHistory } from 'react-router';
import thunkMiddleware from 'redux-thunk';
import { render } from 'react-dom';
import { Provider } from 'react-redux';
import { createStore, applyMiddleware } from 'redux';
import { fetchUser } from './actions';
import App from './components/index';
import Home from './components/home';
import TransactionList from './components/transactionList';
import RecurringTransactionList from './components/recurringTransactionList';
import TransactionTemplateList from './components/transactionTemplateList';
import AccountsPage from './components/accountsPage';
import reducers from './reducers';
import './index.css';

const rootElement = document.getElementById('app');

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

function goToTransactions(nextState, replace, callback) {
    if (store.getState().auth.get('authd')) {
        replace({
            pathname: '/transactions'
        });
    }
    callback();
}

render(
  <Provider store={ store }>
    <Router history={ browserHistory }>
      <Route path="/" component={ App } onEnter={ fetchAuth }>
        <IndexRoute component={ Home } onEnter={ goToTransactions }/>
        <Route path="transactions" component={ AccountsPage } onEnter={ checkAuth }>
          <IndexRoute component={ TransactionList } onEnter={ checkAuth }/>
          <Route path="recurring" component={ RecurringTransactionList } onEnter={ checkAuth }/>
          <Route path="templates" component={ TransactionTemplateList } onEnter={ checkAuth }/>
        </Route>
      </Route>
    </Router>
  </Provider>,
  rootElement
);
