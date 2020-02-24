// eslint-disable-next-line @typescript-eslint/no-empty-interface
import SearchResult from "../model/SearchResult";
import Namable from "../model/Namable";

export interface RootState {
  search: string;
  searchResults: SearchResult<Namable>[];
}
