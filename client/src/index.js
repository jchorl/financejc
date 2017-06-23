import React from 'react';
import ReactDOM from 'react-dom';
import { BrowserRouter } from 'react-router-dom';
import { Provider } from 'react-redux';
import { createStore, applyMiddleware } from 'redux';
import thunk from 'redux-thunk';
import App from './components/App';
import reducers from './reducers';
import registerServiceWorker from './registerServiceWorker';
import './index.css';

let store = createStore(
        reducers,
        applyMiddleware(thunk)
        );

    ReactDOM.render(
            <Provider store={ store }>
                <BrowserRouter>
                    <App />
                </BrowserRouter>
            </Provider>,
            document.getElementById('root')
            );
    registerServiceWorker();
