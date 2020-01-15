import Vue from "vue";
import VueI18n from "vue-i18n";
import App from "./App.vue";
import { TRANSLATIONS } from "./assets/translations/all";
import "./assets/scss/index.scss";
import "typeface-ubuntu-mono";
import "typeface-rokkitt";
import "line-awesome/dist/line-awesome/css/line-awesome.min.css";

Vue.use(VueI18n);

const i18n = new VueI18n({
  locale: "en",
  fallbackLocale: "en",
  messages: TRANSLATIONS
});

new Vue({
  i18n: i18n,
  // eslint-disable-next-line @typescript-eslint/explicit-function-return-type
  render: h => h(App)
}).$mount("#app");
