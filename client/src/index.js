import React from 'react'
import thunkMiddleware from 'redux-thunk'
import createLogger from 'redux-logger'
import { render } from 'react-dom'
import { Provider } from 'react-redux'
import { createStore, applyMiddleware } from 'redux'
import { fetchAccounts } from './actions'
import App from './components/index'
import reducers from './reducers'

const initialState = window.__INITIAL_STATE__
const rootElement = document.getElementById('app')
const loggerMiddleware = createLogger()

let store = createStore(reducers,
						applyMiddleware(
							thunkMiddleware,
							loggerMiddleware
						)
);

render(
	<Provider store={store}>
		<App/>
	</Provider>,
	rootElement
)