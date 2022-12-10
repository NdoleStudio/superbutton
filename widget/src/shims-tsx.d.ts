import Vue, { VNode } from "vue";

declare global {
  namespace JSX {
    interface Element extends VNode {}
    interface ElementClass extends Vue {}
    interface IntrinsicElements {
      [elem: string]: any;
    }
  }

  interface Window {
    SB_USER_ID?: string;
    SB_PROJECT_ID?: string;
  }
}
