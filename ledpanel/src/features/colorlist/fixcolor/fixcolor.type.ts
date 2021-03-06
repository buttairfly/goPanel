export type FixColor = {
  color: string;
  pos: number;
}

export type ColorUpdate = {
  color?: string;
  pos?: number;
}

export type FixColorUpdate = {
  id: string;
  fixColorIndex: number;
  fixColor: ColorUpdate;
}
