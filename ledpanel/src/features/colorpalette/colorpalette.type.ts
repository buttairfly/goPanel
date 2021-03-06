import { Id } from '../../types/id'
import { BlenderId } from '../../types/blender'

export type FixColor = {
  color: string;
  pos: number;
}

export type ColorPalette = {
  id: Id;
  blender: BlenderId;
  colors: FixColor[];
}
