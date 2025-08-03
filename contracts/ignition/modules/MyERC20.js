const { buildModule } = require("@nomicfoundation/hardhat-ignition/modules");

const MyERC20Module = buildModule("MyERC20Module", (m) => {
  const name = m.getParameter("name", "TestToken");
  const symbol = m.getParameter("symbol", "TTK");
  const totalSupply = m.getParameter("totalSupply", 1000000);
  const owner = m.getParameter("owner");

  const token = m.contract("MyERC20", [name, symbol, totalSupply, owner]);

  return { token };
});

module.exports = MyERC20Module;