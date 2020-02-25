import Axios from "axios";

export const Server = Axios.create({
  baseURL: window.location.origin,
  timeout: 5000
});

export const API = Axios.create({
  baseURL: window.location.origin + "/api/",
  timeout: 5000
});
