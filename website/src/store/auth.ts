import { Module } from "vuex";
import { RootState } from "./model";

interface AuthState {
  showLoginModal: boolean;
}

export const AUTH: Module<AuthState, RootState> = {
  state: {
    showLoginModal: true
  },
  mutations: {
    setShowLoginModal: (state, show: boolean): void => {
      state.showLoginModal = show;
    }
  }
};
