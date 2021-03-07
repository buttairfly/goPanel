import React from 'react'
import { useDispatch, useSelector } from 'react-redux'
import sizeMe, { SizeMeProps } from 'react-sizeme'

import {
  selectState,
  selectColorPalettesIds,
  getAllPalettesAsync
} from './colorlist.slice'
import {
  calcPaletteById
} from './colorlist.calc'
import styles from './colorlist.module.css'
import { ColorPaletteComponent } from './colorpalette/colorpalette.component'
import { Button } from 'react-bootstrap'

interface StateFromProps {}

interface DispatchFromProps {}

type Props = StateFromProps & DispatchFromProps & SizeMeProps

const ColorListComponent = (props: Props) => {
  const { width } = props.size
  const palettesState = useSelector(selectState)
  const paletteIds = useSelector(selectColorPalettesIds)

  const dispatch = useDispatch()
  return (
    <div className={styles.container}>
      <Button
          className={styles.asyncButton}
          onClick={() => dispatch(getAllPalettesAsync())}
          variant={'primary'}
        >
          Get Palettes
        </Button>
      { paletteIds.map(paletteId => {
        const palette = calcPaletteById(palettesState, paletteId)
        return (<ColorPaletteComponent
          key={paletteId}
          id={paletteId}
          parentWidth={ width || 100}
          paletteState={palette}
        />
        )
      })}
    </div>
  )
}

export const ColorList = sizeMe()(ColorListComponent)
