import {Module} from "vuex";
import Item from "../../model/Item";
import {RootState} from "../model";
import {HTTP} from "../http";
import {AxiosError, AxiosResponse} from "axios";

export interface ItemsState {
  items: Item[];
  itemsPerPage: number;
  currentIndex: number;
  totalItems: number;
}

export enum PageDirection {
  FIRST,
  PREVIOUS,
  NEXT,
  LAST
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
  getters: {
    isPageFullyContained: (state): boolean => {
      return state.totalItems % state.itemsPerPage === 0;
    },
    numberOfPages: (state, getters): number => {
      let numberOfPages = Math.floor(state.totalItems / state.itemsPerPage);
      if (!getters.isPageFullyContained) {
        numberOfPages += 1;
      }
      return numberOfPages;
    },
    getPageFirstIndex: state => (pageNb: number): number => {
      return state.itemsPerPage * pageNb;
    },
    lastPageFirstIndex: (state, getters): number => {
      const nbPage = getters.numberOfPages;
      const lastPage = nbPage - 1; // pages start at 0
      return getters.getPageFirstIndex(lastPage);
    }
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
          return ctx.dispatch("setCurrentIndex", r.data.offset);
        })
        .catch((err: AxiosError) => {
          console.log(err);
        });
    },
    changePage: (ctx, direction: PageDirection): void => {
      let newIndex = 0;
      switch (direction) {
        case PageDirection.FIRST:
          newIndex = 0;
          break;
        case PageDirection.PREVIOUS:
          newIndex = ctx.state.currentIndex - ctx.state.itemsPerPage;
          break;
        case PageDirection.NEXT:
          newIndex = ctx.state.currentIndex + ctx.state.itemsPerPage;
          break;
        case PageDirection.LAST:
          newIndex = ctx.getters.lastPageFirstIndex;
          break;
      }
      ctx
        .dispatch("setCurrentIndex", newIndex)
        .then(() => {
          return ctx.dispatch("load");
        })
        .catch((err: AxiosError) => {
          console.log(err);
        });
    },
    setCurrentIndex: (ctx, currentIndex: number): void => {
      let toSet = currentIndex;
      const i = ctx.getters.lastPageFirstIndex;
      if (i < currentIndex) {
        toSet = i;
      } else if (currentIndex < 0) {
        toSet = 0;
      }
      ctx.commit("setCurrentIndex", toSet);
    }
  }
};
