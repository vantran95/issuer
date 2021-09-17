// SPDX-License-Identifier: GPL-3.0-or-later
pragma solidity 0.8.6;

import "@openzeppelin/contracts/token/ERC1155/ERC1155.sol";
import "./ERC1888/IERC1888.sol";

/// @title Implementation of the Transferable Certificate standard ERC-1888.
/// @dev Also complies to ERC-1155: https://eips.ethereum.org/EIPS/eip-1155.
/// @dev ** Data set to 0 because there is no meaningful check yet to be done on the data
contract Registry is ERC1155, ERC1888 {

	// Storage for the Certificate structs
	mapping(uint256 => Certificate) public certificateStorage;

	// Mapping from token ID to account balances
	mapping(uint256 => mapping(address => uint256)) public claimedBalances;

	// Incrementing nonce, used for generating certificate IDs
    uint256 internal _latestCertificateId;

	constructor(string memory _uri) ERC1155(_uri) {
		// Trigger ERC1155 constructor
	}

	/// @notice See {IERC1888-issue}.
    /// @dev `_to` cannot be the zero address.
	function issue(address _to, bytes calldata _validityData, uint256 _topic, uint256 _value, bytes calldata _data) external override returns (uint256 id) {
		require(_to != address(0x0), "_to must be non-zero.");
		
		_validate(_msgSender(), _validityData);

		id = ++_latestCertificateId;
		ERC1155._mint(_to, id, _value, new bytes(0)); // Check **

		certificateStorage[id] = Certificate({
			topic: _topic,
			issuer: _msgSender(),
			validityData: _validityData,
			data: _data
		});

		emit IssuanceSingle(_msgSender(), _topic, id, _value);
	}

	/// @notice Allows the issuer to mint more fungible tokens for existing ERC-188 certificates.
    /// @dev `_to` cannot be the zero address.
	function mint(uint256 _id, address _to, uint256 _quantity) external {
		require(_to != address(0x0), "_to must be non-zero.");
		require(_quantity > 0, "_quantity must be above 0.");

		Certificate memory cert = certificateStorage[_id];
		require(_msgSender() == cert.issuer, "Not original issuer");

		ERC1155._mint(_to, _id, _quantity, new bytes(0)); // Check **
	}

	/// @notice See {IERC1888-safeTransferAndClaimFrom}.
    /// @dev `_to` cannot be the zero address.
    /// @dev `_from` has to have a balance above or equal `_value`.
	function safeTransferAndClaimFrom(
		address _from,
		address _to,
		uint256 _id,
		uint256 _value,
		bytes calldata _data,
		bytes calldata _claimData
	) external override {
		Certificate memory cert = certificateStorage[_id];

		_validate(cert.issuer,  cert.validityData);

        require(_to != address(0x0), "_to must be non-zero.");
		require(_from != address(0x0), "_from address must be non-zero.");

        require(_from == _msgSender() || ERC1155.isApprovedForAll(_from, _msgSender()), "No operator approval");
        require(ERC1155.balanceOf(_from, _id) >= _value, "_from balance less than _value");

		if (_from != _to) {
			safeTransferFrom(_from, _to, _id, _value, _data);
		}

		_burn(_to, _id, _value);

		emit ClaimSingle(_from, _to, cert.topic, _id, _value, _claimData); //_claimSubject address ??
	}

	/// @notice See {IERC1888-claimedBalanceOf}.
	function claimedBalanceOf(address _owner, uint256 _id) external override view returns (uint256) {
		return claimedBalances[_id][_owner];
	}

	/// @notice See {IERC1888-getCertificate}.
	function getCertificate(uint256 _id) public view override returns (address issuer, uint256 topic, bytes memory validityCall, bytes memory data) {
		require(_id <= _latestCertificateId, "_id out of bounds");

		Certificate memory certificate = certificateStorage[_id];
		return (certificate.issuer, certificate.topic, certificate.validityData, certificate.data);
	}

	/// @notice Burn certificates after they've been claimed, and increase the claimed balance.
	function _burn(address _from, uint256 _id, uint256 _value) internal override {
		ERC1155._burn(_from, _id, _value);

		claimedBalances[_id][_from] = claimedBalances[_id][_from] + _value;
	}

	/// @notice Validate if the certificate is valid against an external `_verifier` contract.
	function _validate(address _verifier, bytes memory _validityData) internal view {
		(bool success, bytes memory result) = _verifier.staticcall(_validityData);

		require(success && abi.decode(result, (bool)), "Request/certificate invalid");
	}
}
