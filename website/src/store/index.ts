import { StoreOptions } from "vuex";
import { ITEMS } from "./items";
import { STORAGES } from "./storages";
import { RootState } from "./model";

export const STORE_OPTIONS: StoreOptions<RootState> = {
  modules: {
    items: ITEMS,
    storages: STORAGES
  }
};
