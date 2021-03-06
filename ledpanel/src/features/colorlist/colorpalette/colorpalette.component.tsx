import React, { MouseEventHandler } from 'react'
import { useDispatch } from 'react-redux'

import { ColorPalette } from './colorpalette.type'
import { calculateBackgroundStyle } from './colorpalette.calc'
import styles from './colorpalette.module.css'
import { FixColorComponent } from '../fixcolor/fixcolor.component'

import { addFixColor } from '../colorlist.slice'
import { FixColorAddPayload } from '../fixcolor/fixcolor.type'

type Props = {
  id: string;
  paletteState: ColorPalette;
  parentWidth: number;
}

export const ColorPaletteComponent = (props: Props) => {
  const { paletteState, id, parentWidth } = props
  const fixColors = paletteState.colors

  const dispatch = useDispatch()
  const addFixColorClick: MouseEventHandler = (e) => {
    const addFixColorPayload: FixColorAddPayload = {
      id,
      fixColorIndex: paletteState.colors.length,
      fixColor: {
        color: '#ff0',
        pos: e.clientX / parentWidth
      }
    }
    e.stopPropagation()
    dispatch(addFixColor(addFixColorPayload))
  }
  return (
    <div key={id}>
      <div className={styles.paletteName}>
        { id }
      </div>
      <div className={styles.paletteContainer}>
        <div className={styles.palette}>
          <div
            className={styles.paletteBackgroundcolor}
            onClick={addFixColorClick}
            style={calculateBackgroundStyle(paletteState, id)}
          >
            <div className={styles.paletteFixColorContainer}>
              { fixColors.map((_, fixColorIndex) => (
              <FixColorComponent
                key={fixColorIndex}
                fixColorIndex={fixColorIndex}
                parentWidth={parentWidth}
                parentId={id}
              />
              ))}
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
