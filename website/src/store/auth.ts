import { Module } from "vuex";
import { RootState } from "./model";
import { API, Server } from "./http";
import { AxiosResponse } from "axios";
import sleep from "../tools/sleep";
import User from "../model/User";

interface AuthState {
  showLoginModal: boolean;
  authentifiedUser: User | undefined;
}

interface LoginInformations {
  username: string;
  password: string;
}

export const AUTH: Module<AuthState, RootState> = {
  namespaced: true,
  state: {
    showLoginModal: false,
    authentifiedUser: undefined
  },
  getters: {
    isAuthenticated: (state): boolean => {
      return !(
        undefined === state.authentifiedUser ||
        state.authentifiedUser.name === "anonymous"
      );
    }
  },
  mutations: {
    setShowLoginModal: (state, show: boolean): void => {
      state.showLoginModal = show;
    },
    setAuthentified: (state, authentifiedUser: User | undefined): void => {
      let username: string | undefined;
      if (authentifiedUser === undefined) {
        username = undefined;
      } else {
        username = authentifiedUser.name;
      }
      console.log("Change user to: " + username);
      state.authentifiedUser = authentifiedUser;
    }
  },
  actions: {
    loadAuthentication: async (ctx): Promise<void> => {
      const response: AxiosResponse<User> = await API.get("/users/current");
      ctx.commit("setAuthentified", response.data);
    },
    signIn: async (ctx, login: LoginInformations): Promise<void> => {
      const response: AxiosResponse<User> = await Server.get("/auth/sign-in", {
        auth: login
      });
      console.log("Sign in as: " + response.data.name);
      ctx.commit("setAuthentified", response.data);
      ctx.commit("setShowLoginModal", false);
    },
    signOut: async (ctx): Promise<void> => {
      await Server.get("/auth/sign-out", {
        headers: {
          "X-Authentication-Action": "sign-out"
        }
      });
      ctx.commit("setAuthentified", undefined);
    },
    requireLogin: async (ctx): Promise<boolean> => {
      ctx.commit("setShowLoginModal", true);
      while (!ctx.getters.isAuthenticated && ctx.state.showLoginModal) {
        await sleep(200);
      }
      return new Promise((resolve): void => {
        resolve(ctx.getters.isAuthenticated);
      });
    }
  }
};
