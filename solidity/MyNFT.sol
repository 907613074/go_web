// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract MyNFT is ERC721, Ownable {
    uint256 private _tokenCounter;

    constructor() ERC721("MyNFT", "MNFT") Ownable(msg.sender) {}

    mapping(uint256 => string) private tokenIdToMetadataUri;

    function _baseURI() internal view virtual override returns (string memory) {
        return "https://ipfs.io/ipfs/";
    }

    function mintNFT(
        string memory _metadataUri
    ) public onlyOwner returns (uint256) {
        uint256 newTokenId = _tokenCounter + 1;
        tokenIdToMetadataUri[newTokenId] = _metadataUri;
        _mint(msg.sender, newTokenId);
        return newTokenId;
    }

    function tokenURI(
        uint256 tokenId
    ) public view override returns (string memory) {
        return string.concat(_baseURI(), tokenIdToMetadataUri[tokenId]);
    }
}
