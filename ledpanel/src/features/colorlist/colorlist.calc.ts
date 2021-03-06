
import { ColorPaletteListState } from './colorlist.type'
import { ColorPalette } from './colorpalette/colorpalette.type'
import { Id } from '../../types/id'

export const calcPaletteById = (state: ColorPaletteListState, id: Id): ColorPalette => {
  return state.palettes[`${id}`]
}
