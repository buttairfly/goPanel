import React, { MouseEventHandler } from 'react'
import { useDispatch, useSelector } from 'react-redux'

import { ColorPalette } from './colorpalette.type'
import { calculateBackgroundStyle } from './colorpalette.calc'
import styles from './colorpalette.module.css'
import { FixColorComponent } from '../fixcolor/fixcolor.component'

import { addFixColor, selectCurrentColor } from '../colorlist.slice'
import { FixColorAddPayload } from '../fixcolor/fixcolor.type'

type Props = {
  id: string;
  paletteState: ColorPalette;
  parentWidth: number;
}

export const ColorPaletteComponent = (props: Props) => {
  const { paletteState, id, parentWidth } = props
  const fixColors = paletteState.colors

  const width = parentWidth - 2 * 42
  const dispatch = useDispatch()
  const currentColor = useSelector(selectCurrentColor)
  const addFixColorClick: MouseEventHandler = (e) => {
    const addFixColorPayload: FixColorAddPayload = {
      id,
      fixColorIndex: paletteState.colors.length,
      fixColor: {
        color: currentColor,
        pos: e.clientX / width,
        active: true
      }
    }
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
            onDoubleClick={addFixColorClick}
            style={calculateBackgroundStyle(paletteState, id)}
          >
            <div className={styles.paletteFixColorContainer}>
              { fixColors.map((fixColor, fixColorIndex) => (
              <FixColorComponent
                key={`${fixColorIndex}${fixColor.color}`}
                fixColorIndex={fixColorIndex}
                parentWidth={width}
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
