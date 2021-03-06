import React from 'react'
import { Counter } from './features/counter/Counter'
import { ColorPalette } from './features/colorpalette/colorpalette'
import './App.css'

const App = () => {
  return (
    <div className="App">
      <header className="App-header">
        <Counter />
        <ColorPalette />
      </header>
    </div>
  )
}

export default App
