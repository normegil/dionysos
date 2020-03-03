import Rowable from "./Rowable";

export default interface ResourceStore {
  items: Rowable[];
  resourceName: string;
  filter: string;
  total: number;
}
