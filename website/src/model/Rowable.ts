export default interface Rowable {
  identifier(): string;
  toRow(): string[];
}
