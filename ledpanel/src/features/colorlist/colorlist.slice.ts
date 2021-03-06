import { createSlice, PayloadAction } from '@reduxjs/toolkit'
import { RootState } from '../../app/store'
import { BlenderId } from '../../types/blender'
import { ColorPalette } from './colorpalette/colorpalette.type'
import { FixColor, FixColorUpdate } from './fixcolor/fixcolor.type'
import { calcPaletteById } from './colorlist.calc'
import { Id } from '../../types/id'

export interface ColorPaletteListState {
  palettes: {
    [name:string]:ColorPalette;
  }
  currentPaletteName: string;
}

const initialState: ColorPaletteListState = {
  palettes: {
    default: {
      id: 'default',
      blender: BlenderId.RGB,
      colors: [{
        color: '#f00',
        pos: 0.0
      },
      {
        color: '#ff0',
        pos: 0.5
      },
      {
        color: '#0f0',
        pos: 1
      }]
    },
    rainbow: {
      id: 'rainbow',
      blender: BlenderId.RGB,
      colors: [{
        color: '#f00',
        pos: 0.0
      }]
    }
  },
  currentPaletteName: 'default'
}

export const colorPaletteSlice = createSlice({
  name: 'colorPalettes',
  initialState,
  reducers: {
    updateById: (state, action: PayloadAction<ColorPalette>) => {
      state.palettes[`${action.payload.id}`] = action.payload
    },
    updateFixColor: (state, action: PayloadAction<FixColorUpdate>) => {
      console.log(JSON.stringify(action))
      const update = action.payload
      state.palettes[`${update.id}`].colors[update.fixColorIndex] = {
        ...state.palettes[`${update.id}`].colors[update.fixColorIndex],
        ...update.fixColor
      }
    }
  }
})

export const { updateById, updateFixColor } = colorPaletteSlice.actions

export const selectState = (state: RootState): ColorPaletteListState => state.colorPalettes

export const selectColorPalettesIds = (state: RootState) =>
  Object.keys(selectState(state).palettes)

export const selectPalette = (state: RootState, id: String): ColorPalette =>
  calcPaletteById(selectState(state), id)

export const selectCurrentPalette = (state: RootState): ColorPalette =>
  selectPalette(state, selectState(state).currentPaletteName)

export const selectFixColor = (state: RootState, id: Id, index: number): FixColor =>
  selectPalette(state, id).colors[index]

export default colorPaletteSlice.reducer