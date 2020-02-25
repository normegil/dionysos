import { StoreOptions } from "vuex";
import { ITEMS } from "./items";
import { STORAGES } from "./storages";
import { RootState } from "./model";
import SearchResult, { PrintableResult } from "../model/SearchResult";
import { VueRouter } from "vue-router/types/router";
import { HTTP } from "./http";
import { AxiosError, AxiosResponse } from "axios";
import { AUTH } from "./auth";

export const STORE_OPTIONS: StoreOptions<RootState> = {
  state: {
    search: "",
    searchResults: []
  },
  mutations: {
    setSearch: (state, searched: string): void => {
      state.search = searched;
    },
    setSearchResults: (
      state,
      searchResults: SearchResult<PrintableResult>[]
    ): void => {
      state.searchResults = searchResults;
    }
  },
  actions: {
    search: (ctx, router: VueRouter): void => {
      console.log("Searching: " + ctx.state.search);
      if (ctx.state.search === "") {
        console.error("Search query is empty");
        return;
      }
      HTTP.put("/searches", {
        search: ctx.state.search
      })
        .then((r: AxiosResponse<SearchResponse>) => {
          ctx.commit("setSearchResults", r.data.results);
          const pathToReach = "/search";
          if (router.currentRoute.path !== pathToReach) {
            return router.push(pathToReach);
          }
        })
        .catch((err: AxiosError) => {
          console.error(err);
        });
    }
  },
  modules: {
    items: ITEMS,
    storages: STORAGES,
    auth: AUTH
  }
};

interface SearchResponse {
  results: SearchResult<PrintableResult>[];
}
