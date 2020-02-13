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
  totalSize: number;
  offset: number;
  limit: number;
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
    itemsPerPage: 20,
    totalItems: 4
  },
  mutations: {
    setItems: (state, items: Item[]): void => {
      state.items = items;
    },
    setTotalItems: (state, nbItems: number): void => {
      state.totalItems = nbItems;
    },
    setItemsPerPage: (state, itemsPerPage: number): void => {
      state.itemsPerPage = itemsPerPage;
    },
    setCurrentIndex: (state, currentIndex: number): void => {
      state.currentIndex = currentIndex;
    }
  },
  actions: {
    load: (ctx): void => {
      HTTP.get(
        "/items?limit=" +
          ctx.state.itemsPerPage +
          ";offset=" +
          ctx.state.currentIndex
      )
        .then((r: AxiosResponse<ItemCollection>) => {
          const itemCollection = r.data.items.map(
            (dto): Item => new Item(dto.id, dto.name)
          );
          ctx.commit("setItems", itemCollection);
          ctx.commit("setTotalItems", r.data.totalSize);
          ctx.commit("setItemsPerPage", r.data.limit);
          ctx.commit("setCurrentIndex", r.data.offset);
        })
        .catch((err: AxiosError) => {
          console.log(err);
        });
    }
  }
};
