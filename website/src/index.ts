import Vue from "vue";
import App from "./App.vue";
import "./assets/scss/index.scss";
import "typeface-lato";

new Vue({
  // eslint-disable-next-line @typescript-eslint/explicit-function-return-type
  render: h => h(App)
}).$mount("#app");
