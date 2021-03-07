import React, { MouseEventHandler } from 'react'
import { useDispatch, useSelector } from 'react-redux'

import { ColorPalette } from './colorpalette.type'
import { calculateBackgroundStyle } from './colorpalette.calc'
import styles from './colorpalette.module.css'
import { FixColorComponent } from '../fixcolor/fixcolor.component'

import { addFixColor, selectIsDragging } from '../colorlist.slice'
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
  const isDragging = useSelector(selectIsDragging)
  const addFixColorClick: MouseEventHandler = (e) => {
    const addFixColorPayload: FixColorAddPayload = {
      id,
      fixColorIndex: paletteState.colors.length,
      fixColor: {
        color: '#ff0',
        pos: e.clientX / parentWidth
      }
    }
    console.log(isDragging)
    if (isDragging === false) {
      dispatch(addFixColor(addFixColorPayload))
    }
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
