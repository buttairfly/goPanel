import { ColorPalette } from './colorpalette.type'
import { Id } from '../../../types/id'
import { FixColor } from '../fixcolor/fixcolor.type'

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
  const colors: FixColor[] = JSON.parse(JSON.stringify(p.colors))
  const sortedColors = colors.sort((a, b) => a.pos - b.pos)
  sortedColors.forEach(c => {
    s += `, ${c.color} ${c.pos * 100}%`
  })
  return {
    background: `linear-gradient(to right${s})`
  }
}
