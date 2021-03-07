import { ColorPalette } from './colorpalette/colorpalette.type'

export type ColorPaletteListState = {
  palettes: ColorPaletteList;
  currentColor: string;
  currentPaletteName: string;
}

export type ColorPaletteList = {
  [name:string]:ColorPalette;
}
