import axios from 'axios'
import { useDispatch } from 'react-redux'

import { ColorPaletteList } from './colorlist.type'
import config from '../../app/config'
import { getAllPalettes } from './colorlist.slice'

export const useAllPalettes = async () => {
  try {
    const palettes: ColorPaletteList = (await axios.get(`${config.baseUrl}/api/v1/panel/palette/all`)).data
    console.log(palettes)
    const dispatch = useDispatch()
    dispatch(getAllPalettes(palettes))
    return palettes
  } catch (e) {
    console.log(JSON.stringify(e))
  }
}
