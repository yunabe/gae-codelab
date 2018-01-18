import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import App from './App';
import registerServiceWorker from './registerServiceWorker';
import { spaApp } from './reducers';
import { Provider } from 'react-redux'
import { createStore } from 'redux'

let store = createStore(spaApp)

ReactDOM.render(<Provider store={store}><App /></Provider>, document.getElementById('root'));
registerServiceWorker();

// Load a message from /api/hello.
setTimeout(() => {
  let xhr = new XMLHttpRequest();
  xhr.addEventListener('load', ()=>{
    store.dispatch(
      {type: 'SET_TOP_PAGE_MESSAGE',
       message: xhr.responseText});
  });
  xhr.open('GET', '/api/hello')
  xhr.send();
}, 1000);
