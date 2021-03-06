import type { Id } from '../../types/id'
import type { BlenderId } from '../../types/blender'
import type { FixColor } from './fixcolor/fixcolor.type'

export type ColorPalette = {
  id: Id;
  blender: BlenderId;
  colors: FixColor[];
}
