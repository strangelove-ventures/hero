# Stability Report

Written by: Strangelove Ventures

Published: January 6, 2023


## Abstract 
Hero chain is an application specific Cosmos chain developed and tested by the team at Strangelove designed to host the deployment of real world assets native to the IBC ecosystem. This chain is expected to be deployed using Interchain Security provided by the Cosmos Hub – a $3.5B market cap blockchain. A feature complete snapshot of the chain development was deployed as part of an incentivized public testnet that was stewarded by the Cosmos community, primarily by the teams at Informal Systems and Hypha Cooperative. The team found that in most critical ways both the application written in the Cosmos SDK and the Interchain Security mechanism behaved reliably. Exceptions and their remediations are noted in the document below as is the plan for further testing. 

The key properties related to compliance features, integration with IBC, governance were all tested in real world conditions. 


## About this Report

This stability report covers testing undertaken on the Hero chain by the Strangelove Team. The majority of the testing took place during the public incentivized testnet – Game of Chains [[Leaderboard](https://interchainsecurity.dev/game-of-chains-2022), [Announcement Blog Post](https://blog.cosmos.network/announcing-game-of-chains-open-for-registration-d1818662de8e)]. The testnet started on November 7, 2022 and finished on December 9, 2022. While this report covers some critical aspects of [Interchain Security](https://github.com/cosmos/interchain-security) (ICS) and the [Admin Module](https://github.com/Ethernal-Tech/admin-module/), this report will focus on the [Tokenfactory Module](https://github.com/strangelove-ventures/hero/tree/main/x/tokenfactory). This module provides the key functionality of Hero allowing privileged accounts to mint assets and blacklist users. The module was also the one piece that was directly developed by Strangelove. 


#### Approaches for the stability analysis
* Game of Chains monitoring
* IBC unit testing and CI/CD workflow
* Manual code reviews
* Audit of Tokenfactory module/middleware, Admin Module, Packet Forward Middleware (to begin Jan 17, 2023)


## Induction into Game Of Chains

On November 21, 2022, using [Interchain Security](https://github.com/cosmos/interchain-security), Hero was successfully added as a consumer chain to the testnet. In a real world-context, the provider chain is the high-market cap proof of stake chain which will be able to offer high security guarantees to Cosmos blockchains in return for a fee paid to validators. In this testnet, the provider chain was a [forked version of Gaia](https://github.com/cosmos/gaia/tree/goc-dec-7), Gaia being the chain binary that runs the Cosmos Hub. This fork was necessary because it included the [pre-release of Interchain Security](https://github.com/cosmos/interchain-security/releases/tag/v0.1). 

After being on-boarded as a consumer chain, Hero acted as expected:

* Hero properly adopted the provider chains validator set
* IBC transactions and IBC functionality worked as expected
* Chain output logs were stable


## Provider Chain Halt and Eventual Hero Halt

On December 8, 2022, the provider chain halted while the community tested out a “slash throttle” feature. At the time of the halt, research was being tracked [here](https://github.com/hyphacoop/ics-testnets/tree/goc-dec-halt-incident/game-of-chains-2022/incidents). The bug has since been found and [fixed](https://github.com/cosmos/interchain-security/pull/605). 

In the current version of ICS, if a provider chain halts, consumer chains continue to produce blocks and accept IBC transactions until the trusting period of the [light-client](https://github.com/cosmos/ibc/tree/main/spec/core/ics-002-client-semantics) expires. 

With a client trusting period of 1 week, Hero continued to run smoothly without the provider chain. 

On December 14, 2022, one week after the halt of the provider chain, the light-client between Hero and the provider expired, causing Hero to also halt.

This is extremely unlikely in a real world scenario as Hero is intended to launch as a consumer chain of the Cosmos Hub. To our knowledge, the Cosmos Hub has not had any unplanned downtime since launch. Other precautions can also be taken, such as specifying a longer client trusting period between consumer and provider.  

It is possible to recover a consumer chain in this situation, but requires code changes and cooperation from validators.



## Tokenfactory Module


The [Tokenfactory Module](https://github.com/strangelove-ventures/hero/tree/main/x/tokenfactory) in the Hero chain binary allows generic assets to be minted and controlled by privileged accounts. Privileged accounts can also place a pause on all transactions to and from the chain.

A list of all possible privledges tied to the Tokenfactory can be seen [here](../readme.md#access-control).

Throughout the duration of the Game of Chains Testnet, Strangelove manually tested many Tokenfactory commands. At a high level, we tested:

* Delegation and creation of all privileged accounts (Master Minter, Minter, Blacklist, etc,)
* Minting of assets
* Blacklisting functionality
* Pause functionality


A detailed list of all commands and their outcomes are available in this [table](#testnet-commands)

In addition to the Game of Chains tests, there is a CI/CD pipeline using [ibctest](https://github.com/strangelove-ventures/ibctest) built into the github repo to continuously validate Tokenfactory functionality as features are added and changed.


## Tokenfactory Module Discoveries


During the course of the testnet, we discovered and fixed three critical bugs:

* [[PR #1]](https://github.com/strangelove-ventures/hero/pull/1) It was impossible to burn assets minted by tokenfactory 
* [[PR #4]](https://github.com/strangelove-ventures/hero/pull/4) A non-blacklisted user was able to send tokenfactory assets to a blacklisted user
* [[PR #4]](https://github.com/strangelove-ventures/hero/pull/4) A blacklisted user was unable to send/receive assets not minted by the tokenfactory


Due to the Game of Chains halt, we were unable to incorporate these changes during the testnet. However, rigorous test cases were added to the [ibctest CI pipeline](https://github.com/strangelove-ventures/hero/blob/main/ibctest/ibctest_test.go) to validate the fixes and prevent future regressions. Since the complexity of these bugs were minimal, we found these ibctest cases sufficient in ensuring confidence.


## Testnet commands

Below are the Tokenfactory commands manually ran during the Game of Chains Testnet. 


For the below commands, "user1" and "user2" are normal unprivileged accounts. The remaining accounts should be self explanatory from their name via the [Access Control Table](../readme.md#access-control). For situations where the wallet address was needed instead of the key name, simply hover over the address to see the tool tip describing what privileged account the address is tied to.


The linked images are there to show what the output of the command looks like in either a blockchain explorer or the command line. 


<!-- cosmos1q4adyc9a75u6eclu6czqtr2vfyqf6v4svwwnyr -->
[owner]: ## "Owner"
<!-- cosmos1fmtkwd7awdyavzk78yh75wt0tl7vcme89f42xg --> 
[master_minter]: ## "Master Minter"
<!-- cosmos1erqarmzn2ae5vd4j6sk77hccs0yfs854jee8wd --> 
[mint_controller]: ## "Mint Controller"
<!-- cosmos1uydlnqjz5mjfgafs9gxm5lhww4hdffldmu2j6y --> 
[minter]: ## "Minter"
<!-- cosmos1w9v9nfa2lvfhx3u6g8l67g64deu9d4wvjav7vk --> 
[blacklister]: ## "Blacklister"
<!-- cosmos1earswrsr7xhl8mgkglvfvcxqxqsflpu9cvun93 --> 
[pauser]: ## "Pauser"
<!-- cosmos1qtw23j6y9758juk6mnm937uands7vyxsxkhptu --> 
[user1]: ## "User1"
<!-- cosmos1zaamktfya6ps3ektdpt7xgka5nu9hx2w8ruthx --> 
[user2]: ## "User2"

| Description | Summary | Pass/Fail | Command | Images |
|:---:|:---:|:---:|:---:|:---:|
| Update Master Minter from Non-Privleged Account | This tx correctly failed due to the tx being signed from an account without the “owner”  privilege. | 🟢 | herod tx tokenfactory update-master-minter [cosmos1fmtkwd7awdyavzk78yh75wt0tl7vcme89f42xg][master_minter] --from user1 --fees 500uhero | [x](./images/stability_report/update_master_minter_nonPrivlidged/) |
| Update Master Minter | Successfully updated Master Minter account using privileged “owner” account | 🟢 | herod tx tokenfactory update-master-minter [cosmos1fmtkwd7awdyavzk78yh75wt0tl7vcme89f42xg][master_minter] --from owner --fees 500uhero | [x](./images/stability_report/update_master_minter/) |
| Configure Mint Controller | Successfully updated the “Mint Controller” and “Minter” account using the privileged “Master Minter” account. | 🟢 | herod tx tokenfactory configure-minter-controller [cosmos1erqarmzn2ae5vd4j6sk77hccs0yfs854jee8wd][mint_controller] [cosmos1uydlnqjz5mjfgafs9gxm5lhww4hdffldmu2j6y][minter] --fees 500uhero --from masterminter | [x](./images/stability_report/configure_mint_controller/) |
| Configure Minter | Successfully enabled Minter account to have an allowance of 1,000,000 udrachma. | 🟢 | herod tx tokenfactory configure-minter [cosmos1uydlnqjz5mjfgafs9gxm5lhww4hdffldmu2j6y][minter] 1000000udrachma --fees 500uhero --from minter-controller | [x](./images/stability_report/configure_minter/) |
| Mint Asset | The Minter successfully minted 100 udrachma into a users account | 🟢 | herod tx tokenfactory mint [cosmos1qtw23j6y9758juk6mnm937uands7vyxsxkhptu][user1] 100udrachma --from minter --fees 500uhero | [x](./images/stability_report/mint_asset/) |
| Mint Asset - Over Allowance | Attempted to mint an amount of asset over the allowance allocated by the Mint Controller. This transaction correctly failed because 9,000,000 was over the allotted 1,000,000 allocated by the Mint Controller. | 🟢 | herod tx tokenfactory mint [cosmos1qtw23j6y9758juk6mnm937uands7vyxsxkhptu][user1] 9000000udrachma --from minter --fees 500uhero | [x](./images/stability_report/mint_asset_over_allowance/) |
| Update Blacklister | Successfully created the blacklist account using the privileged “Owner” account. | 🟢 | herod tx tokenfactory update-blacklister [cosmos1w9v9nfa2lvfhx3u6g8l67g64deu9d4wvjav7vk][blacklister] --fees 500uhero --from owner | [x](./images/stability_report/update_blacklister/) |
| Blacklist User | Successfully blacklisted a users account using the privileged “Blacklister” account. Note: when querying for blacklisted accounts, it is not possible to get a list of all blacklisted accounts. It is only possible to query whether a specified address is blacklisted. | 🟢 | herod tx tokenfactory blacklist [cosmos1qtw23j6y9758juk6mnm937uands7vyxsxkhptu][user1] --from black-lister --fees 500uhero | [x](./images/stability_report/blacklist_user/) |
| Mint Asset to Blacklisted User | The Minter was correctly not able to mint 100 udrachma into a blacklisted users account | 🟢 | herod tx tokenfactory mint [cosmos1qtw23j6y9758juk6mnm937uands7vyxsxkhptu][user1] 100udrachma --from minter --fees 500uhero | [x](./images/stability_report/mint_asset_to_blacklisted_user/) |
| Send tokenfactory asset to Blacklisted User | It was incorrectly possible for a non-blacklisted user to send minted funds to a blacklisted user. This was a bug discovered and [fixed](https://github.com/strangelove-ventures/hero/pull/4). | ❌ | herod tx bank send user2 [cosmos1qtw23j6y9758juk6mnm937uands7vyxsxkhptu][user1] 100udrachma --fees 500uhero | [x](./images/stability_report/send_tokenfactory_asset_to_bl_user/) |
| Send Tokenfactory Asset from Blacklisted User | Attempted to send a minted asset using a blacklisted account. The client acted correctly and the tx was never broadcast. | 🟢 | herod tx bank send [cosmos1qtw23j6y9758juk6mnm937uands7vyxsxkhptu][user1] [cosmos1zaamktfya6ps3ektdpt7xgka5nu9hx2w8ruthx][user2]100udrachma --fees 500uhero | [x](./images/stability_report/send_tokenfactory_asset_from_bl_user/) |
| Un-Blacklist User | Successfully un-blacklisted a user account using the privileged “Blacklister” account. The Minter was then correctly able to mint assets into the previously blacklisted account.  | 🟢 | herod tx tokenfactory unblacklist [cosmos1qtw23j6y9758juk6mnm937uands7vyxsxkhptu][user1] --fees 500uhero --from black-lister | [x](./images/stability_report/unblacklist_user/) |
| Update Pauser | Successfully updated Pauser account using privileged “owner” account. Note: when querying for the pauser address, the address itself is not needed.  | 🟢 | herod tx tokenfactory  update-pauser [cosmos1earswrsr7xhl8mgkglvfvcxqxqsflpu9cvun93][pauser] --from owner --fees 500uhero | [x](./images/stability_report/update_pauser/) |
| Pause | Successfully paused chain using the privileged “pauser” account | 🟢 | herod tx tokenfactory pause --fees 500uhero --from pauser | [x](./images/stability_report/pause/) |
| Mint Asset While Paused | The Minter correctly failed to mint 100 udrachma into a user account because the chain was paused. | 🟢 | herod tx tokenfactory mint [cosmos1qtw23j6y9758juk6mnm937uands7vyxsxkhptu][user1] 100udrachma --from minter --fees 500uhero | [x](./images/stability_report/mint_asset_while_paused/) |
| Send Tokenfactory Asset While Paused | Attempted to send a minted asset while the chain was paused. The client acted correctly and the tx was never broadcast | 🟢 | herod tx bank send user1 [cosmos1zaamktfya6ps3ektdpt7xgka5nu9hx2w8ruthx][user2] 100udrachma --fees 500uhero | [x](./images/stability_report/send_tokenfactory_asset_while_paused/) |
| Send Non-Tokenfactory Asset While Paused | Attempted to send a non-tokenfactory asset while the chain was paused. The client acted correctly and the tx was never broadcast. | 🟢 | herod tx bank send user1 [cosmos1zaamktfya6ps3ektdpt7xgka5nu9hx2w8ruthx][user2] 100uhero --fees 500uhero | [x](./images/stability_report/send_non_tokenfactory_asset_while_paused/) |
| Un-Pause | Successfully un-paused chain using the privileged “pauser” account | 🟢 | herod tx tokenfactory unpause --fees 500uhero --from pauser | [x](./images/stability_report/un_pause/) |



## Admin Module


The [Admin Module](https://github.com/Ethernal-Tech/admin-module/)  was built by the [Ethernal Team](https://github.com/Ethernal-Tech). This module allows the chain to be upgraded by bypassing the provider chain's validator set which would otherwise rely on a two week public governance process.
We were unable to test this in the Game of Chains Testnet due to the provider chain halt and subsequent Hero chain halt. However, we have confirmed this functionality works with an “[ibctest](https://github.com/strangelove-ventures/ibctest)” test case.

It should be noted that once [Cosmos SDK](https://github.com/cosmos/cosmos-sdk) 0.47 is released, this custom module will not be needed. 
The Gov module in this upgraded version of the SDK will have this functionality natively built in.



## Future Testing and Validation

* Lambda testnet (Set to launch mid to early January 2023)
* Integration tests with asset issuer(s)
* Continuing to build out more robust [CI/CD tests](https://github.com/strangelove-ventures/hero/tree/main/ibctest) using the [ibctest](https://github.com/strangelove-ventures/ibctest) test suite


## Conclusion


Besides the three bugs discovered in the Tokenfactory module, the Hero chain worked and functioned as intended. The halt of Game of Chains hindered some testing capacity but ibctest proved a reliable replacement.

While the chain functionality itself is easy to validate with the ibctest suite, the intricacies of Interchain Security along with more rigorous widespread testing of Tokenfactory behavior can be further explored during the [Cosmos SDK Rho](https://hub.cosmos.network/main/roadmap/cosmos-hub-roadmap-2.0.html#v8-rho-upgrade-expected-q1-2023) Testnet and via a third party audit conducted by [Oak Security](https://www.oaksecurity.io/). 