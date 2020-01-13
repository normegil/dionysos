import Vue from "vue";
import App from "./App.vue";
import "./assets/scss/index.scss";
import "typeface-lato";
import "line-awesome/dist/line-awesome/css/line-awesome.min.css";

new Vue({
  // eslint-disable-next-line @typescript-eslint/explicit-function-return-type
  render: h => h(App)
}).$mount("#app");
