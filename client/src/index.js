import React from 'react'
import thunkMiddleware from 'redux-thunk'
import { render } from 'react-dom'
import { Provider } from 'react-redux'
import { createStore, applyMiddleware } from 'redux'
import App from './components/index'
import reducers from './reducers'

const rootElement = document.getElementById('app')

let store = createStore(
  reducers,
  applyMiddleware(
    thunkMiddleware
  )
);

render(
  <Provider store={store}>
    <App/>
  </Provider>,
  rootElement
)
