import React, { useEffect, useReducer, useState } from 'react'
import Draggable, { DraggableEventHandler } from 'react-draggable'
import { useSelector } from 'react-redux'

import {
  calculateFixpointBackgroundStyle,
  selectFixColor
} from './fixcolor.calc'
import styles from './fixcolor.module.css'
import { FixColorUpdate } from './fixcolor.type'

import reducer, { updateFixColor, selectState } from '../colorlist.slice'

type Props = {
  id: string;
  fixColorIndex: number;
  parentWidth: number;
}

export const FixColorComponent = (props: Props) => {
  const { id, parentWidth, fixColorIndex } = props
  const width = parentWidth - 2 * 42
  const [state, dispatch] = useReducer(reducer, useSelector(selectState))
  const fixColor = selectFixColor(state, id, fixColorIndex)

  const [pos, setPos] = useState(fixColor.pos)

  useEffect(() => {
    const fixColorUpdate: FixColorUpdate = {
      id,
      fixColorIndex,
      fixColor: {
        pos
      }
    }
    console.log(width, pos, pos / width, JSON.stringify(fixColorUpdate))
    if (pos !== fixColor.pos) { dispatch(updateFixColor(fixColorUpdate)) }
  })

  const roundDecimals = (num: number): number => {
    const decimals = 1000
    return Math.round(num * decimals) / decimals
  }

  const updateFixColorPos: DraggableEventHandler = (e, position) => {
    setPos(roundDecimals(position.x / width))
  }

  const changeLabelPos: =

  return (
      <Draggable
        position={{ x: pos * width, y: 0 }}
        axis='x'
        bounds='parent'
        onDrag={updateFixColorPos}
      >
        <div
          className={styles.paletteFixColor}
        >
          <div className={styles.paletteFixColorBackground}/>
          <div className={styles.paletteFixColorVisual}/>
          <div
            className={styles.paletteFixColorBackgroundcolor}
            style={calculateFixpointBackgroundStyle(fixColor)}
          />
          <label className={styles.paletteFixColorLabel}>
            <div className={styles.paletteFixColorLabelBackground}>
              <input
                className={styles.paletteFixColorLabelInput}
                value={roundDecimals(pos)}
                onChange={changeLabelPos}
              />
              &nbsp;
            </div>
          </label>
        </div>
      </Draggable>
  )
}
