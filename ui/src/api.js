
/**
 * Alice public theming and extension API
 */

const apiFunc = () => {
  let subs = [];
  let props = [];

  const apply = () => {
    let called = false;
    for (const sub of subs) {
      for (const prop of props) {
        sub(...prop);
        called = true;
      }
    }
    if (called) {
      subs = [];
    }
  }

  const call = (...params) => {
    props.push(params);
    apply();
  };

  const subscribe = (fn) => {
    subs = [...subs, fn];
    apply();
  };

  return ([
    call, subscribe
  ])
}

export const apiCallback = () => {
  const [call, subscribe] = apiFunc();
  const callback = (...args) => subscribe((fn) => fn(...args));
  return [call, callback];
}

export const [updateContent, updateContentApi] = apiFunc();
export const [onLayoutReady, onLayoutReadyApi] = apiCallback();

const Api = {
  updateContent,
  onLayoutReady,
};

export default Api;
