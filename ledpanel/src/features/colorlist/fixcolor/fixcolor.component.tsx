import React, { ChangeEventHandler, useEffect, useState } from 'react'
import Draggable, { DraggableEventHandler } from 'react-draggable'
import { useSelector, useDispatch } from 'react-redux'

import {
  calculateFixpointBackgroundStyle,
  selectFixColor
} from './fixcolor.calc'
import styles from './fixcolor.module.css'
import { FixColorUpdatePayload } from './fixcolor.type'

import { updateFixColor, selectState, toggleDragging } from '../colorlist.slice'

type Props = {
  parentId: string;
  fixColorIndex: number;
  parentWidth: number;
}

export const FixColorComponent = (props: Props) => {
  const { parentId, parentWidth, fixColorIndex } = props
  const width = parentWidth - 2 * 42
  const dispatch = useDispatch()
  const state = useSelector(selectState)
  const fixColor = selectFixColor(state, parentId, fixColorIndex)

  const [pos, setPos] = useState(fixColor.pos)

  useEffect(() => {
    const fixColorUpdate: FixColorUpdatePayload = {
      id: parentId,
      fixColorIndex,
      fixColor: {
        pos
      }
    }
    dispatch(updateFixColor(fixColorUpdate))
  }, [pos])

  const toggleDrag: DraggableEventHandler = (e, position) => {
    dispatch(toggleDragging())
  }

  const updateFixColorPos: DraggableEventHandler = (e, position) => {
    setPos(position.x / width)
  }

  const changeLabelPos: ChangeEventHandler<any> = (e) => {
    const val = e.target.value
    const num = val.replace(/a-z%!"§%&\/\(\)=\?@;-_#\+¿*/gmi, '')
    const num2 = num.replace(/,/, '.')
    const newPos: number = parseFloat(num2)
    if (!isNaN(newPos)) {
      if (newPos > 100) {
        setPos(1)
      } else {
        setPos(newPos / 100.0)
      }
    }
    console.log(val)
  }
  return (
      <Draggable
        position={{ x: pos * width, y: 0 }}
        axis='x'
        bounds='parent'
        onStart={toggleDrag}
        onStop={toggleDrag}
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
                value={`${((pos) * 100).toFixed(1)}%`}
                onChange={changeLabelPos}
              />
              &nbsp;
            </div>
          </label>
        </div>
      </Draggable>
  )
}
