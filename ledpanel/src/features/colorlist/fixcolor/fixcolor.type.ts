export type FixColor = {
  color: string;
  pos: number;
}

export type FixColorUpdate = {
  color?: string;
  pos?: number;
}

export type FixColorRemovePayload = {
  id: string;
  fixColorIndex: number;
}

export interface FixColorUpdatePayload extends FixColorRemovePayload {
  fixColor: FixColorUpdate;
}

export interface FixColorAddPayload extends FixColorRemovePayload {
  fixColor: FixColor;
}
