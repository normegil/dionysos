import Axios, { AxiosError } from "axios";
import { STORE } from "./index";

export const Server = Axios.create({
  baseURL: window.location.origin,
  timeout: 5000
});

export const API = Axios.create({
  baseURL: window.location.origin + "/api/",
  timeout: 5000
});

const errorHandler = (error: AxiosError): Promise<AxiosError> => {
  if (undefined !== error.response) {
    if (error.response.status === 401) {
      STORE.commit("auth/setShowLoginModal", true);
    }
  }
  // eslint-disable-next-line prefer-promise-reject-errors
  return Promise.reject({ error });
};

API.interceptors.response.use(
  response => response,
  error => errorHandler(error)
);
