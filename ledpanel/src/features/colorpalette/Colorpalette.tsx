// eslint-disable-next-line no-use-before-define
import React from 'react'
import { useSelector } from 'react-redux'
import {
  selectColorPalettesState,
  calculateBackgroundStyle,
  selectColorPalettesIds
} from './colorpaletteSlice'
import styles from './Colorpalette.module.css'

export function ColorPalette () {
  const palettes = useSelector(selectColorPalettesState)
  const paletteIds = useSelector(selectColorPalettesIds)

  return (
    <div className={styles.container}>
      {
        paletteIds.map(paletteId => {
          return (
            <div key={paletteId} className={styles.palette} style={calculateBackgroundStyle(palettes, paletteId)}>
              <div className={'w-25 p-3'}>
                  { paletteId }
              </div>
            </div>
          )
        })
      }
    </div>
  )
}
