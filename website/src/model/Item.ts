import Rowable from "./Rowable";

export default class Item implements Rowable {
  private _id: string;
  name: string;

  constructor(id: string, name: string) {
    this._id = id;
    this.name = name;
  }

  identifier(): string {
    return this._id;
  }

  get id(): string {
    return this._id;
  }

  toRow(): string[] {
    return [this.name];
  }
}
