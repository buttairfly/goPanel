import React from 'react'
import { useSelector } from 'react-redux'
import {
  selectColorPalettesState,
  calculateBackgroundStyle,
  calculateFixpointBackgroundStyle,
  selectColorPalettesIds,
  selectById
} from './colorpaletteSlice'
import styles from './colorpalette.module.css'
import { FixColor } from './colorpalette.type'

const renderFixColors = (colors: FixColor[]) => {
  return colors.map((fixColor, index) => {
    const xPos = `${fixColor.pos * 100}%`
    return (
      <div
        key={index}
        className={styles.paletteFixColor}
        data-x={xPos}
        style={{ left: xPos }}
      >
        <div className={styles.paletteFixColorBackground}/>
        <div className={styles.paletteFixColorVisual}/>
        <div className={styles.paletteFixColorBackgroundcolor} style={calculateFixpointBackgroundStyle(fixColor)}/>
        <label className={styles.paletteFixColorLabel}>
          <div className={styles.paletteFixColorLabelBackground}>
            <input className={styles.paletteFixColorLabelInput} value={xPos}/>
            &nbsp;
          </div>
        </label>
      </div>
    )
  })
}

export const ColorPalette = () => {
  const palettes = useSelector(selectColorPalettesState)
  const paletteIds = useSelector(selectColorPalettesIds)
  return (
    <div className={styles.container}>
      { paletteIds.map(paletteId => {
        const colors = selectById(palettes, paletteId).colors
        return (
            <div key={paletteId} className={styles.paletteContainer}>
              <div className={styles.palette}>
                <div
                  className={styles.paletteBackgroundcolor}
                  style={calculateBackgroundStyle(palettes, paletteId)}
                >
                  <div className={styles.paletteFixColorContainer}>
                    { renderFixColors(colors) }
                  </div>
                  { paletteId }
                </div>
              </div>
            </div>
        )
      })}
    </div>
  )
}
