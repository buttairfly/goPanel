import React from 'react'
import { Counter } from './features/counter/Counter'
import { ColorPalette } from './features/colorpalette/Colorpalette'
import './App.css'

function App () {
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
