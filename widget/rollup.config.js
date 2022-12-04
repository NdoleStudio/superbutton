import vue from "rollup-plugin-vue";

export default [
  // Browser build.
  {
    input: "src/main.ts",
    output: {
      format: "iife",
      file: "dist/widget.js",
    },
    plugins: [vue()],
  },
];
