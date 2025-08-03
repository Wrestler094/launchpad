const { buildModule } = require("@nomicfoundation/hardhat-ignition/modules");

const PresaleModule = buildModule("PresaleModule", (m) => {
  const token = m.getParameter("token");
  const rate = m.getParameter("rate", 100); // 100 tokens per ETH
  const hardCap = m.getParameter("hardCap"); // in ETH
  const softCap = m.getParameter("softCap"); // in ETH
  const durationInDays = m.getParameter("durationInDays", 30);
  const owner = m.getParameter("owner");

  const presale = m.contract("Presale", [
    token,
    rate,
    hardCap,
    softCap,
    durationInDays,
    owner
  ]);

  return { presale };
});

module.exports = PresaleModule;