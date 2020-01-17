import Axios from "axios";

export const HTTP = Axios.create({
  baseURL: window.location.origin + "/api/",
  timeout: 1000
});
