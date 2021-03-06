import React from 'react'

import { ColorPalette } from './colorpalette.type'
import { calculateBackgroundStyle } from './colorpalette.calc'
import styles from './colorpalette.module.css'
import { FixColorComponent } from '../fixcolor/fixcolor.component'

type Props = {
  id: string;
  paletteState: ColorPalette;
  parentWidth: number;
}

export const ColorPaletteComponent = (props: Props) => {
  const { paletteState, id, parentWidth } = props
  const fixColors = paletteState.colors
  return (
    <div key={id}>
      <div className={styles.paletteName}>
        { id }
      </div>
      <div className={styles.paletteContainer}>
        <div className={styles.palette}>
          <div
            className={styles.paletteBackgroundcolor}
            style={calculateBackgroundStyle(paletteState, id)}
          >
            <div className={styles.paletteFixColorContainer}>
              { fixColors.map((fixColor, fixColorIndex) => (
              <FixColorComponent
                key={fixColorIndex}
                id={id}
                fixColorIndex={fixColorIndex}
                parentWidth={parentWidth}
              />
              ))}
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
