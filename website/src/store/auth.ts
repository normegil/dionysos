import { Module } from "vuex";
import { RootState } from "./model";
import { Server } from "./http";
import { AxiosResponse } from "axios";

interface AuthState {
  showLoginModal: boolean;
}

interface LoginInformations {
  username: string;
  password: string;
}

export const AUTH: Module<AuthState, RootState> = {
  namespaced: true,
  state: {
    showLoginModal: false
  },
  mutations: {
    setShowLoginModal: (state, show: boolean): void => {
      state.showLoginModal = show;
    }
  },
  actions: {
    signIn: (ctx, login: LoginInformations): Promise<void> => {
      return Server.get("/auth/sign-in", {
        auth: login
      }).then((r: AxiosResponse) => {
        console.log("Login success");
      });
    }
  }
};
