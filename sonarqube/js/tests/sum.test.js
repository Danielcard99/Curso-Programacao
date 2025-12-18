const sum = require("../src/sum");

test("add 1 + 2 to be equal 3", () => {
  expect(sum(1, 2)).toBe(3);
});
