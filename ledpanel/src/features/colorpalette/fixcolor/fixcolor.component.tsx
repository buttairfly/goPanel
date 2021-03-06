import React, { Component } from 'react'
import {
  calculateFixpointBackgroundStyle
} from '../colorpaletteSlice'
import styles from './fixcolor.module.css'
import { FixColor } from './fixcolor.type'

type Props = {
  fixColors: FixColor[];
}

export class FixColorComponent extends Component<Props> {
  render = () => {
    const { fixColors } = this.props
    return fixColors.map((fixColor, index) => {
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
}
