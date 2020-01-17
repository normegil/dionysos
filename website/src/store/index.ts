import { StoreOptions } from "vuex";
import { ITEMS } from "./modules/items";
import { RootState } from "./model";

export const STORE_OPTIONS: StoreOptions<RootState> = {
  modules: {
    items: ITEMS
  }
};
