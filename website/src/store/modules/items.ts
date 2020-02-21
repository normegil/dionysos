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
  filter: string;
}

interface ItemCollection {
  items: ItemDTO[];
  totalSize: number;
  offset: number;
  limit: number;
  filter: string;
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
    totalItems: 4,
    filter: ""
  },
  getters: {

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
    },
    setFilter: (state, filter: string): void => {
      state.filter = filter;
    }
  },
  actions: {
    load: (ctx): void => {
      let url =
        "/items?limit=" +
        ctx.state.itemsPerPage +
        ";offset=" +
        ctx.state.currentIndex;
      if (ctx.state.filter !== "") {
        url += ";filter=" + ctx.state.filter;
      }
      HTTP.get(url)
        .then((r: AxiosResponse<ItemCollection>) => {
          return ctx.dispatch("refreshItems", r.data);
        })
        .catch((err: AxiosError) => {
          console.log(err);
        });
    },
    save: (ctx, item: ItemDTO): void => {
      HTTP.put("/items", item)
        .then(() => {
          return ctx.dispatch("load");
        })
        .catch((err: AxiosError) => {
          console.log(err);
        });
    },
    delete: (ctx, id: string): void => {
      HTTP.delete("/items/" + id)
        .then(() => {
          return ctx.dispatch("load");
        })
        .catch((err: AxiosError) => {
          console.log(err);
        });
    },
    refreshItems: (ctx, data: ItemCollection): void => {
      const itemCollection = data.items.map(
        (dto): Item => new Item(dto.id, dto.name)
      );
      ctx.commit("setItems", itemCollection);
      ctx.commit("setTotalItems", data.totalSize);
      ctx.commit("setItemsPerPage", data.limit);
      ctx.commit("setFilter", data.filter);
      ctx.dispatch("setCurrentIndex", data.offset).catch((err: AxiosError) => {
        console.log(err);
      });
    }
  }
};
