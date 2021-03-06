import { ColorPalette } from './colorpalette/colorpalette.type'

export type ColorPaletteListState = {
  palettes: ColorPaletteList
  currentPaletteName: string;
}

export type ColorPaletteList = {
  [name:string]:ColorPalette;
}
