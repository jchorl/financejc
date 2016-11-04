import React from 'react';
import { Router, IndexRoute, Route, browserHistory } from 'react-router';
import thunkMiddleware from 'redux-thunk';
import { render } from 'react-dom';
import { Provider } from 'react-redux';
import { createStore, applyMiddleware } from 'redux';
import App from './components/index';
import Home from './components/home';
import Recurring from './components/recurring';
import reducers from './reducers';

const rootElement = document.getElementById('app')

let store = createStore(
  reducers,
  applyMiddleware(
    thunkMiddleware
  )
);

render(
  <Provider store={ store }>
    <Router history={ browserHistory }>
      <Route path="/" component={ App }>
        <IndexRoute component={ Home } />
        <Route path="recurring" component={ Recurring } />
      </Route>
    </Router>
  </Provider>,
  rootElement
)
