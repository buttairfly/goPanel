import React from 'react'
import { useSelector } from 'react-redux'
import sizeMe, { SizeMeProps } from 'react-sizeme'

import {
  selectState,
  selectColorPalettesIds
} from './colorlist.slice'
import {
  calcPaletteById
} from './colorlist.calc'
import styles from './colorlist.module.css'
import { ColorPaletteComponent } from './colorpalette/colorpalette.component'

interface StateFromProps {}

interface DispatchFromProps {}

type Props = StateFromProps & DispatchFromProps & SizeMeProps

const ColorListComponent = (props: Props) => {
  const { width } = props.size

  const palettesState = useSelector(selectState)
  const paletteIds = useSelector(selectColorPalettesIds)
  return (
    <div className={styles.container}>
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
