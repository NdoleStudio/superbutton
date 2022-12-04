import Vue from "vue";
import Widget from "./Widget.vue";

Vue.config.productionTip = false;

new Vue({
  render: (h) => h(Widget),
}).$mount("#app");
