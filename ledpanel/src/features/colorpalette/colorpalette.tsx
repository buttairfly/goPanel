import React from 'react'
import { useSelector } from 'react-redux'
import {
  selectColorPalettesState,
  calculateBackgroundStyle,
  selectColorPalettesIds,
  selectById
} from './colorpaletteSlice'
import styles from './colorpalette.module.css'
import { FixColorComponent } from './fixcolor/fixcolor.component'

export const ColorPalette = () => {
  const palettes = useSelector(selectColorPalettesState)
  const paletteIds = useSelector(selectColorPalettesIds)
  return (
    <div className={styles.container}>
      { paletteIds.map(paletteId => {
        const colors = selectById(palettes, paletteId).colors
        return (
          <div key={paletteId}>
            <div className={styles.paletteName}>
              { paletteId }
            </div>
            <div className={styles.paletteContainer}>
              <div className={styles.palette}>
                <div
                  className={styles.paletteBackgroundcolor}
                  style={calculateBackgroundStyle(palettes, paletteId)}
                >
                  <div className={styles.paletteFixColorContainer}>
                    <FixColorComponent fixColors={colors} />
                  </div>
                </div>
              </div>
            </div>
          </div>
        )
      })}
    </div>
  )
}
