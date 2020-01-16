import Vue from "vue";
import VueI18n from "vue-i18n";
import VueRouter from "vue-router";
import App from "./App.vue";
import { ROUTES } from "./routes/routes";
import { TRANSLATIONS } from "./assets/translations/all";
import "./assets/scss/index.scss";
import "typeface-ubuntu-mono";
import "typeface-rokkitt";
import "line-awesome/dist/line-awesome/css/line-awesome.min.css";

Vue.use(VueI18n);
Vue.use(VueRouter);

const i18n = new VueI18n({
  locale: "en",
  fallbackLocale: "en",
  messages: TRANSLATIONS
});

const router = new VueRouter({
  routes: ROUTES
});

let vm = new Vue({
  i18n: i18n,
  router: router,
  // eslint-disable-next-line @typescript-eslint/explicit-function-return-type
  render: h => h(App)
}).$mount("#app");

vm.$router.push("/items");
