import Vue from "vue";
import Widget from "./Widget.vue";

Vue.config.productionTip = false;

const div = document.createElement("div");
div.id = "sb-w-app";
document.body.appendChild(div);
new Vue({
  render: (h) => h(Widget),
}).$mount("#sb-w-app");
