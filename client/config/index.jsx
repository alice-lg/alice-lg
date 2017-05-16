

export const configureAxios = function(axios) {
  // Setup axios to use django xsrf token
  axios.defaults.xsrfCookieName = 'csrftoken';
  axios.defaults.xsrfHeaderName = 'X-CSRFToken';
};

