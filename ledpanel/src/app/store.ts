import { configureStore, ThunkAction, Action } from '@reduxjs/toolkit'
import counterReducer from '../features/counter/counterSlice'
import colorPaletteReducer from '../features/colorlist/colorlist.slice'

export const store = configureStore({
  reducer: {
    counter: counterReducer,
    colorPalettes: colorPaletteReducer
  }
})

export type RootState = ReturnType<typeof store.getState>;
export type AppThunk<ReturnType = void> = ThunkAction<
  ReturnType,
  RootState,
  unknown,
  Action<string>
>;
