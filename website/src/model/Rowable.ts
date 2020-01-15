export default interface Rowable {
  id(): string;
  toRow(): string[];
}
