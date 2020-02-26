import { Module } from "vuex";
import { RootState } from "./model";
import { Server } from "./http";
import { AxiosRequestConfig } from "axios";
import sleep from "../tools/sleep";

interface AuthState {
  showLoginModal: boolean;
  currentContext: AxiosRequestConfig | undefined;
  authentified: boolean;
}

interface LoginInformations {
  username: string;
  password: string;
}

export const AUTH: Module<AuthState, RootState> = {
  namespaced: true,
  state: {
    showLoginModal: false,
    currentContext: undefined,
    authentified: false
  },
  mutations: {
    setShowLoginModal: (state, show: boolean): void => {
      state.showLoginModal = show;
    },
    setAuthentified: (state, authentified: boolean): void => {
      state.authentified = authentified;
    }
  },
  actions: {
    signIn: async (ctx, login: LoginInformations): Promise<void> => {
      await Server.get("/auth/sign-in", { auth: login });
      ctx.commit("setAuthentified", true);
    },
    requireLogin: async (ctx): Promise<void> => {
      ctx.commit("setShowLoginModal", true);
      while (!ctx.state.authentified) {
        await sleep(200);
      }
    }
  }
};
