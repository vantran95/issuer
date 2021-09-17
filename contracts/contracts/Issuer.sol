// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.7.0 <0.9.0;
/// @title Voting with delegation.

import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";

import "./Registry.sol";

/// @title Issuer contract
/// @notice Used to manage the request/approve workflow for issuing ERC-1888 certificates.
contract Issuer is Initializable, OwnableUpgradeable, UUPSUpgradeable {

    // Certificate topic - check ERC-1888 topic description
    uint256 public certificateTopic;

    // ERC-1888 contract to issue certificates to
    Registry public registry;

    // Optional: Private Issuance contract
    address public privateIssuer;

    // Storage for CertificationRequest structs
    mapping(uint256 => CertificationRequest) private _certificationRequests;

    // Mapping from request ID to certificate ID
    mapping(uint256 => uint256) private _requestToCertificate;

    // Incrementing nonce, used for generating certification request IDs
    uint256 private _latestCertificationRequestId;

    // Mapping to whether a certificate with a specific ID has been revoked by the Issuer
    mapping(uint256 => bool) private _revokedCertificates;

    struct CertificationRequest {
        address owner; // Owner of the request
        bytes data;
        bool approved;
        bool revoked;
        address sender;  // User that triggered the request creation
    }

    event CertificationRequested(address indexed _owner, uint256 indexed _id);
    event CertificationRequestApproved(address indexed _owner, uint256 indexed _id, uint256 indexed _certificateId);
    event CertificationRequestRevoked(address indexed _owner, uint256 indexed _id);

    event CertificateRevoked(uint256 indexed _certificateId);
    event CertificateVolumeMinted(address indexed _owner, uint256 indexed _certificateId, uint256 indexed _volume);

    /// @notice Contructor.
    /// @dev Uses the OpenZeppelin `initializer` for upgradeability.
    /// @dev `_registry` cannot be the zero address.
    function initialize(uint256 _certificateTopic, address _registry) public initializer {
        require(_registry != address(0), "Cannot use address 0x0");

        certificateTopic = _certificateTopic;

        registry = Registry(_registry);
        OwnableUpgradeable.__Ownable_init();
        UUPSUpgradeable.__UUPSUpgradeable_init();
    }

    /// @notice Attaches a private issuance contract to this issuance contract.
    /// @dev `_privateIssuer` cannot be the zero address.
    function setPrivateIssuer(address _privateIssuer) public onlyOwner {
        require(_privateIssuer != address(0), "Cannot use address 0x0");
        require(privateIssuer == address(0), "Private issuer already set");

        privateIssuer = _privateIssuer;
    }

    /*
        Certification requests
    */

    function getCertificationRequest(uint256 _requestId) public view returns (CertificationRequest memory) {
        return _certificationRequests[_requestId];
    }

    function requestCertificationFor(bytes memory _data, address _owner) public returns (uint256) {
        require(_owner != address(0), "Owner cannot be 0x0");

        uint256 id = ++_latestCertificationRequestId;

        _certificationRequests[id] = CertificationRequest({
        owner: _owner,
        data: _data,
        approved: false,
        revoked: false,
        sender: _msgSender()
        });

        emit CertificationRequested(_owner, id);

        return id;
    }

    function requestCertification(bytes calldata _data) external returns (uint256) {
        return requestCertificationFor(_data, _msgSender());
    }

    function tryGetCertId(uint256 _requestId) public view returns (uint256) {
        return _requestToCertificate[_requestId];
    }

    function revokeRequest(uint256 _requestId) external {
        CertificationRequest storage request = _certificationRequests[_requestId];

        require(_msgSender() == request.owner || _msgSender() == OwnableUpgradeable.owner(), "Unauthorized revoke request");
        require(!request.revoked, "Already revoked");
        require(!request.approved, "Can't revoke approved requests");

        request.revoked = true;

        emit CertificationRequestRevoked(request.owner, _requestId);
    }

    function revokeCertificate(uint256 _certificateId) external onlyOwner {
        require(!_revokedCertificates[_certificateId], "Already revoked");
        _revokedCertificates[_certificateId] = true;

        emit CertificateRevoked(_certificateId);
    }

    function approveCertificationRequest(
        uint256 _requestId,
        uint256 _value
    ) public returns (uint256) {
        require(_msgSender() == owner() || _msgSender() == privateIssuer, "caller not owner or issuer");
        require(_requestNotApprovedOrRevoked(_requestId), "already approved or revoked");

        CertificationRequest storage request = _certificationRequests[_requestId];
        request.approved = true;

        uint256 certificateId = registry.issue(
            request.owner,
            abi.encodeWithSignature("isRequestValid(uint256)",_requestId),
            certificateTopic,
            _value,
            request.data
        );

        _requestToCertificate[_requestId] = certificateId;

        emit CertificationRequestApproved(request.owner, _requestId, certificateId);

        return certificateId;
    }

    /// @notice Directly issue a certificate without going through the request/approve procedure manually.
    function issue(address _to, uint256 _value, bytes memory _data) public onlyOwner returns (uint256) {
        uint256 requestId = requestCertificationFor(_data, _to);

        return approveCertificationRequest(
            requestId,
            _value
        );
    }

    /// @notice Mint more volume to existing certificates
    /// @param _to To whom the volume should be minted.
    /// @param _certificateId The ID of the certificate.
    /// @param _volume Volume that should be minted.
    function mint(address _to, uint256 _certificateId, uint256 _volume) external onlyOwner {
        registry.mint(_certificateId, _to, _volume);

        emit CertificateVolumeMinted(_to, _certificateId, _volume);
    }

    /// @notice Validation for certification requests.
    /// @dev Used by other contracts to validate the token.
    /// @dev `_requestId` has to be an existing ID.
    function isRequestValid(uint256 _requestId) external view returns (bool) {
        require(_requestId <= _latestCertificationRequestId, "cert request ID out of bounds");

        CertificationRequest memory request = _certificationRequests[_requestId];
        uint256 certificateId = _requestToCertificate[_requestId];

        return request.approved
        && !request.revoked
        && !_revokedCertificates[certificateId];
    }

    /*
        Info
    */

    function getRegistryAddress() external view returns (address) {
        return address(registry);
    }

    function getPrivateIssuerAddress() external view returns (address) {
        return privateIssuer;
    }

    function version() external pure returns (string memory) {
        return "v0.1";
    }

    /*
        Private methods
    */

    function _requestNotApprovedOrRevoked(uint256 _requestId) internal view returns (bool) {
        CertificationRequest memory request = _certificationRequests[_requestId];

        return !request.approved && !request.revoked;
    }

    /// @notice Needed for OpenZeppelin contract upgradeability.
    /// @dev Allow only to the owner of the contract.
    function _authorizeUpgrade(address) internal override onlyOwner {
        // Allow only owner to authorize a smart contract upgrade
    }
}