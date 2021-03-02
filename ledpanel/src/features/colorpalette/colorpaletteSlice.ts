import { createSlice, PayloadAction } from '@reduxjs/toolkit'
import { RootState } from '../../app/store'
import { Id } from '../../types/id'
import { BlenderId } from '../../types/blender'

type FixColor = {
  color: string;
  pos: number;
}

type ColorPalette = {
  id: Id;
  blender: BlenderId;
  colors: FixColor[];
}

interface ColorPaletteState {
  palettes: {
    [name:string]:ColorPalette;
  }
}

const initialState: ColorPaletteState = {
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
      }]
    }
  }
}

export const colorPaletteSlice = createSlice({
  name: 'palette',
  initialState,
  reducers: {
    updateById: (state, action: PayloadAction<ColorPalette>) => {
      state.palettes[`${action.payload.id}`] = action.payload
    }
  }
})

export const { updateById } = colorPaletteSlice.actions

export const selectColorPalettesIds = (state: RootState) => Object.keys(state.colorPalette.palettes)
export const selectColorPalettesState = (state: RootState) => state.colorPalette

export const calculateBackgroundStyle = (state: ColorPaletteState, id: Id) => {
  const p = selectById(state, id)
  if (p.colors.length < 1) {
    return {
      background: '#000'
    }
  }
  if (p.colors.length < 2) {
    return {
      background: `${p.colors[0].color}`
    }
  }

  let s = ''
  p.colors.forEach(c => {
    s += `, ${c.color} ${c.pos * 100}%`
  })
  return {
    background: `linear-gradient(to right${s})`
  }
}
export const selectById = (state: ColorPaletteState, id: Id): ColorPalette => {
  return state.palettes[`${id as string}`]
}

export default colorPaletteSlice.reducer
