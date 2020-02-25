import Vuex from "vuex";
import { ITEMS } from "./items";
import { STORAGES } from "./storages";
import SearchResult, { PrintableResult } from "../model/SearchResult";
import { VueRouter } from "vue-router/types/router";
import { API } from "./http";
import { AxiosError, AxiosResponse } from "axios";
import { AUTH } from "./auth";
import Vue from "vue";

Vue.use(Vuex);

export const STORE = new Vuex.Store({
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
      API.put("/searches", {
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
});

interface SearchResponse {
  results: SearchResult<PrintableResult>[];
}
