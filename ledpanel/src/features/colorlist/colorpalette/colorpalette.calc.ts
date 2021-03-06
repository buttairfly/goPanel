import { ColorPalette } from './colorpalette.type'
import { Id } from '../../../types/id'

export const calculateBackgroundStyle = (p: ColorPalette, id: Id) => {
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
