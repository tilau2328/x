// GENERATED BY MITOSIS

import { Fragment, component$, h, useStore } from "@builder.io/qwik";
type Props = {
  message: string;
};
export const Text = component$((props: Props) => {
  const state = useStore<any>({ name: "Foo" });
  return (
    <div>
      {props.message || "Hello"}
      {state.name}! I can run in React, Vue, Solid or Svelte!
    </div>
  );
});
export default Text;
