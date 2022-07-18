
/**
 * Check if state is up or established
 */
export const isUpState = (s) => {
    if (!s) { return false; }
    s = s.toLowerCase();
    return (s.includes("up") || s.includes("established"));
}
