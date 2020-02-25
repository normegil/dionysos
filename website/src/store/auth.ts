import { Module } from "vuex";
import { RootState } from "./model";
import { Server } from "./http";
import { AxiosError, AxiosResponse } from "axios";
import Error from "../model/Error";

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
    showLoginModal: true
  },
  mutations: {
    setShowLoginModal: (state, show: boolean): void => {
      state.showLoginModal = show;
    }
  },
  actions: {
    signIn: (ctx, login: LoginInformations): void => {
      Server.get("/auth/sign-in", {
        auth: login
      })
        .then((r: AxiosResponse) => {
          console.log("Login success");
        })
        .catch((err: AxiosError<Error>) => {
          console.error(err);
        });
    }
  }
};
