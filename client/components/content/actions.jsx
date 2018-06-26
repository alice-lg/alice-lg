
export const CONTENT_UPDATE = "@content/CONTENT_UPDATE";

export function contentUpdate(content) {
  return {
    type: CONTENT_UPDATE,
    payload: content
  }
}

