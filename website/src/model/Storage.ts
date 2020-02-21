import Rowable from "./Rowable";

export default class Storage implements Rowable {
  id: string;
  name: string;

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
