import { Module } from "vuex";
import Item from "../../model/Item";
import { RootState } from "../model";
import { HTTP } from "../http";
import { AxiosError, AxiosResponse } from "axios";

export interface ItemsState {
  items: Item[];
  itemsPerPage: number;
  currentIndex: number;
  totalItems: number;
}

interface ItemCollection {
  items: ItemDTO[];
}

interface ItemDTO {
  id: string;
  name: string;
}

export const ITEMS: Module<ItemsState, RootState> = {
  namespaced: true,
  state: {
    items: [
      new Item("0", "Tomates"),
      new Item("1", "Oignons"),
      new Item("2", "Carotte"),
      new Item("3", "Spaggethi")
    ],
    currentIndex: 0,
    itemsPerPage: 10,
    totalItems: 4
  },
  mutations: {
    set: (state, items: Item[]): void => {
      console.log("set data: " + items);
      state.items = items;
    }
  },
  actions: {
    load: (ctx): void => {
      HTTP.get("/items")
        .then((r: AxiosResponse<ItemCollection>) => {
          const itemCollection = r.data.items.map(
            (dto): Item => new Item(dto.id, dto.name)
          );
          ctx.commit("set", itemCollection);
        })
        .catch((err: AxiosError) => {
          console.log(err);
        });
    }
  }
};
