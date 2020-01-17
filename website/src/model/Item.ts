import Rowable from "./Rowable";

export default class Item implements Rowable {
  private id: string;
  private name: string;

  constructor(id: string, name: string) {
    this.id = id;
    this.name = name;
  }

  identifier(): string {
    return this.id;
  }

  toRow(): string[] {
    return [this.name];
  }
}
