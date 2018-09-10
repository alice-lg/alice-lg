
/*
 * Helper: Get info from api error
 */

export const infoFromError = function(error) {
    if (error.response && error.response.data && error.response.data.code) {
      return error.response.data;
    }
    return null;
}

