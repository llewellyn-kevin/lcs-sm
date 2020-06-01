import React from 'react';
import {StockTicker} from './components/StockTicker.js'
import './App.css';

function App() {
  return (
    <div className="App">
      <header className="App-header">
        <h2>League Stock Market</h2>
      </header>
      <StockTicker />
    </div>
  );
}

export default App;
