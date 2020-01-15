import Rowable from "./Rowable";

export default class Item implements Rowable {
  private identifier: string;
  private name: string;

  constructor(identifier: string, name: string) {
    this.identifier = identifier;
    this.name = name;
  }

  id(): string {
    return this.identifier;
  }

  toRow(): string[] {
    return [this.name];
  }
}
