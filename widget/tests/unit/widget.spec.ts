import { shallowMount } from "@vue/test-utils";
import Widget from "@/Widget.vue";

describe("Widget.vue", () => {
  it("renders empty ext on mount", () => {
    const wrapper = shallowMount(Widget, {});
    expect(wrapper.text()).toMatch("");
  });
});
