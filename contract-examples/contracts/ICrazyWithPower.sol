// (c) 2022-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// SPDX-License-Identifier: MIT

pragma solidity >=0.8.0;

interface ICrazyWithPower {
  // steal(address enemy) wipes away all funds from an [enemy]and gives them to the caller address,
  // but if their protection works caller address funds will be given to [enemy].
  function steal(address enemy) external;

  // uncertainFate() doubles your funds or nothing. This function can only be used once.
  function uncertainFate() external;

  // setProtection(uint256 protection) allows one to set protection from other opposition trying to steal funds.
  function setProtection(uint256 protection) external;
}
