// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "@openzeppelin/contracts/token/ERC1155/IERC1155.sol";

interface ERC1888 is IERC1155 {

    struct Certificate {
        uint256 topic;
        address issuer;
        bytes validityData;
        bytes data;
    }

   event IssuanceSingle(address indexed _issuer, uint256 indexed _topic, uint256 _id, uint256 _value);

   event ClaimSingle(address indexed _claimIssuer, address indexed _claimSubject, uint256 indexed _topic, uint256 _id, uint256 _value, bytes _claimData);

   function issue(address _to, bytes calldata _validityData, uint256 _topic, uint256 _value, bytes calldata _data) external returns (uint256 id);

   function safeTransferAndClaimFrom(address _from, address _to, uint256 _id, uint256 _value, bytes calldata _data, bytes calldata _claimData) external;

   function claimedBalanceOf(address _owner, uint256 _id) external view returns (uint256);

   function getCertificate(uint256 _id) external view returns (address issuer, uint256 topic, bytes memory validityCall, bytes memory data);
}