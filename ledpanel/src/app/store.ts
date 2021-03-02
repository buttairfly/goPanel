import { configureStore, ThunkAction, Action } from '@reduxjs/toolkit';
import counterReducer from '../features/counter/counterSlice';
import colorPaletteReducer from '../features/colorpalette/colorpaletteSlice';

export const store = configureStore({
  reducer: {
    counter: counterReducer,
    colorPalette: colorPaletteReducer,
  },
});

export type RootState = ReturnType<typeof store.getState>;
export type AppThunk<ReturnType = void> = ThunkAction<
  ReturnType,
  RootState,
  unknown,
  Action<string>
>;
