import Rowable from "./Rowable";
import Namable from "./Namable";
import Identifiable from "./Identifiable";

export default class Item implements Rowable, Identifiable, Namable {
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
