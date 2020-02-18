const Registry = artifacts.require('Registry')
const KeepBonding = artifacts.require("KeepBonding")
const BondedECDSAKeepFactory = artifacts.require("BondedECDSAKeepFactory")
const BondedECDSAKeepVendor = artifacts.require("BondedECDSAKeepVendor")
const BondedECDSAKeepVendorImplV1 = artifacts.require("BondedECDSAKeepVendorImplV1")

const { deployBondedSortitionPoolFactory } = require('@keep-network/sortition-pools/migrations/scripts/deployContracts')
const BondedSortitionPoolFactory = artifacts.require("BondedSortitionPoolFactory")

let { RandomBeaconAddress, TokenStakingAddress, RegistryAddress } = require('./external-contracts')

module.exports = async function (deployer) {
    await deployBondedSortitionPoolFactory(artifacts, deployer)

    if (process.env.TEST) {
        TokenStakingStub = artifacts.require("TokenStakingStub")
        TokenStakingAddress = (await TokenStakingStub.new()).address

        RandomBeaconStub = artifacts.require("RandomBeaconStub")
        RandomBeaconAddress = (await RandomBeaconStub.new()).address

        RegistryAddress = (await deployer.deploy(Registry)).address
    } else {
        RegistryAddress = (await Registry.at(RegistryAddress)).address
    }

    await deployer.deploy(
        KeepBonding,
        RegistryAddress,
        TokenStakingAddress
    )

    await deployer.deploy(
        BondedECDSAKeepFactory,
        BondedSortitionPoolFactory.address,
        TokenStakingAddress,
        KeepBonding.address,
        RandomBeaconAddress
    )

    await deployer.deploy(BondedECDSAKeepVendorImplV1)
    await deployer.deploy(BondedECDSAKeepVendor, BondedECDSAKeepVendorImplV1.address)
}
