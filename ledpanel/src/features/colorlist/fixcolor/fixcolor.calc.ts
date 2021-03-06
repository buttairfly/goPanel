import { ColorPaletteListState } from '../colorlist.type'
import { FixColor } from './fixcolor.type'
import { calcPaletteById } from '../colorlist.calc'
import { Id } from '../../../types/id'

export const calculateFixpointBackgroundStyle = (fixColor: FixColor) => {
  return {
    background: `${fixColor.color}`
  }
}

export const selectFixColor = (state: ColorPaletteListState, id: Id, index: number): FixColor =>
  calcPaletteById(state, id).colors[index]
