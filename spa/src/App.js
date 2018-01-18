import React, { Component } from 'react';
import logo from './logo.svg';
import './App.css';
import { connect } from 'react-redux'

class App extends Component {
  render() {
    return (
      <div className="App">
        <header className="App-header">
          <img src={logo} className="App-logo" alt="logo" />
          <h1 className="App-title">Welcome to React</h1>
        </header>
        <p className="App-intro">{this.props.message}</p>
      </div>
    );
  }
}

App = connect((state) => {
  let message = state.topPage.message;
  if (!message) {
    message = 'Loading...';
  }
  return {message};
})(App);

export default App;
