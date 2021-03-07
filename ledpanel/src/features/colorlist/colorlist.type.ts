import { ColorPalette } from './colorpalette/colorpalette.type'

export type ColorPaletteListState = {
  palettes: ColorPaletteList
  isDragging: boolean;
  currentPaletteName: string;
}

export type ColorPaletteList = {
  [name:string]:ColorPalette;
}
