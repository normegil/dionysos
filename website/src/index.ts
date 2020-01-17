import Vue from "vue";
import VueI18n from "vue-i18n";
import VueRouter from "vue-router";
import Vuex from "vuex";
import App from "./App.vue";
import { ROUTER } from "./routes";
import { STORE_OPTIONS } from "./store";
import { TRANSLATOR_OPTIONS } from "./assets/translations";
import "./assets/scss/index.scss";
import "typeface-ubuntu-mono";
import "typeface-rokkitt";
import "line-awesome/dist/line-awesome/css/line-awesome.min.css";

Vue.use(VueI18n);
Vue.use(VueRouter);
Vue.use(Vuex);

const i18n = new VueI18n(TRANSLATOR_OPTIONS);
const store = new Vuex.Store(STORE_OPTIONS);

const vm = new Vue({
  i18n: i18n,
  router: ROUTER,
  store: store,
  // eslint-disable-next-line @typescript-eslint/explicit-function-return-type
  render: h => h(App)
}).$mount("#app");

vm.$router.push("/items");
