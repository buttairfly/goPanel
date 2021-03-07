import React, { useCallback, useEffect, useState } from 'react'
import Draggable from 'react-draggable'
import { useSelector, useDispatch } from 'react-redux'

import {
  calculateFixpointBackgroundStyle,
  selectFixColor
} from './fixcolor.calc'
import styles from './fixcolor.module.css'
import { FixColorRemovePayload, FixColorUpdatePayload } from './fixcolor.type'

import { updateFixColor, removeFixColor, selectState } from '../colorlist.slice'
import { Button } from 'react-bootstrap'

type Props = {
  parentId: string;
  fixColorIndex: number;
  parentWidth: number;
}

export const FixColorComponent = (props: Props) => {
  const { parentId, parentWidth, fixColorIndex } = props
  const dispatch = useDispatch()
  const state = useSelector(selectState)
  const fixColor = selectFixColor(state, parentId, fixColorIndex)

  const [pos, setPos] = useState(fixColor.pos)
  const [label, setLabel] = useState(`${((pos) * 100).toFixed(1)}%`)
  const [active, setActive] = useState(false)

  useEffect(() => {
    const fixColorUpdate: FixColorUpdatePayload = {
      id: parentId,
      fixColorIndex,
      fixColor: {
        pos,
        active
      }
    }
    dispatch(updateFixColor(fixColorUpdate))
  }, [pos])

  const updateFixColorPos = useCallback((e, position) => {
    const newPos = position.x / parentWidth
    setPos(newPos)
    setLabel(`${((newPos) * 100).toFixed(1)}%`)
    setActive(true)
  }, [])

  const deleteFixColor = useCallback(() => {
    const fixColor: FixColorRemovePayload = {
      id: parentId,
      fixColorIndex
    }
    dispatch(removeFixColor(fixColor))
  }, [fixColorIndex])

  const changeLabelPos = useCallback((e) => {
    const val = e.target.value
    const num1 = val.replace(/,/, '.')
    const num2 = num1.replace(/a-z%!"§%&\/\(\)=\?@;-_#\+¿*/gmi, '')
    const num3 = num2.replace(/\.$/, '')
    const newPos: number = parseFloat(num3)
    if (!isNaN(newPos)) {
      if (newPos <= 100.0 && newPos >= 0.0) {
        setPos(newPos / 100.0)
        setLabel(`${((newPos)).toFixed(1)}%`)
      }
    }
    setLabel(e.target.value)
  }, [])

  return (
      <Draggable
        position={{ x: pos * parentWidth, y: 0 }}
        axis='x'
        bounds='parent'
        onStart={updateFixColorPos}
        onDrag={updateFixColorPos}
      >
        <div className={styles.paletteFixColor}>
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
                value={label}
                onChange={changeLabelPos}
              />
              &nbsp;
            </div>
          </label>
          { fixColor?.active && (
            <Button
              className={styles.removeButton}
              onClick={deleteFixColor}
              variant="danger">X</Button>
          )}
        </div>
      </Draggable>
  )
}
