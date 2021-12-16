
let uniqueId = 0;
export default function generateUniqueId() {
  return `with-g2-${uniqueId++}`;
}