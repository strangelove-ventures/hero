# Stability Report


## Testnet commands

[owner]: ## "Owner"
<!-- cosmos1fmtkwd7awdyavzk78yh75wt0tl7vcme89f42xg --> [master_minter]: ## "Master Minter"
<!-- cosmos1erqarmzn2ae5vd4j6sk77hccs0yfs854jee8wd --> [mint_controller]: ## "Mint Controller"
<!-- cosmos1uydlnqjz5mjfgafs9gxm5lhww4hdffldmu2j6y --> [minter]: ## "Minter"
<!-- cosmos1w9v9nfa2lvfhx3u6g8l67g64deu9d4wvjav7vk --> [blacklister]: ## "Blacklister"
<!-- cosmos1earswrsr7xhl8mgkglvfvcxqxqsflpu9cvun93 --> [pauser]: ## "Pauser"
<!-- cosmos1qtw23j6y9758juk6mnm937uands7vyxsxkhptu --> [user1]: ## "User1"
<!-- cosmos1zaamktfya6ps3ektdpt7xgka5nu9hx2w8ruthx --> [user2]: ## "User2"

| Title | Summary | Command | Images |
|:---:|:---:|:---:|:---:|
| Update Master Minter from Non-Privleged Account | This tx correctly failed due to the tx being signed from an account without the “owner”  privilege. | herod tx tokenfactory update-master-minter [cosmos1fmtkwd7awdyavzk78yh75wt0tl7vcme89f42xg][master_minter] --from user1 --fees 500uhero | x |
| Update Master Minter | Successfully updated Master Minter account using privileged “owner” account | herod tx tokenfactory update-master-minter [cosmos1fmtkwd7awdyavzk78yh75wt0tl7vcme89f42xg][master_minter] --from owner --fees 500uhero | x |
| Configure Mint Controller | Successfully updated the “Mint Controller” and “Minter” account using the privileged “Master Minter” account. | herod tx tokenfactory configure-minter-controller [cosmos1erqarmzn2ae5vd4j6sk77hccs0yfs854jee8wd][mint_controller] [cosmos1uydlnqjz5mjfgafs9gxm5lhww4hdffldmu2j6y][minter] --fees 500uhero --from masterminter | x |
| Configure Minter | Successfully enabled Minter account to have an allowance of 1,000,000 udrachma. | herod tx tokenfactory configure-minter [cosmos1uydlnqjz5mjfgafs9gxm5lhww4hdffldmu2j6y][minter] 1000000udrachma --fees 500uhero --from minter-controller | x |
| Mint Asset | The Minter successfully minted 100 udrachma into a users account | herod tx tokenfactory mint [cosmos1qtw23j6y9758juk6mnm937uands7vyxsxkhptu][user1] 100udrachma --from minter --fees 500uhero | x |
| Mint Asset - Over Allowance | Attempted to mint an amount of asset over the allowance allocated by the Mint Controller. This transaction correctly failed because 9,000,000 was over the allotted 1,000,000 allocated by the Mint Controller. | herod tx tokenfactory mint [cosmos1qtw23j6y9758juk6mnm937uands7vyxsxkhptu][user1] 9000000udrachma --from minter --fees 500uhero | x |
| Update Blacklister | Successfully created the blacklist account using the privileged “Owner” account. | herod tx tokenfactory update-blacklister [cosmos1w9v9nfa2lvfhx3u6g8l67g64deu9d4wvjav7vk][blacklister] --fees 500uhero --from owner | x |
| Blacklist User | Successfully blacklisted a users account using the privileged “Blacklister” account. Note: when querying for blacklisted accounts, it is not possible to get a list of all blacklisted accounts. It is only possible to query whether a specified address is blacklisted. | herod tx tokenfactory blacklist [cosmos1qtw23j6y9758juk6mnm937uands7vyxsxkhptu][user1] --from black-lister --fees 500uhero | x |
| Mint Asset to Blacklisted User | The Minter was correctly not able to mint 100 udrachma into a blacklisted users account | herod tx tokenfactory mint [cosmos1qtw23j6y9758juk6mnm937uands7vyxsxkhptu][user1] 100udrachma --from minter --fees 500uhero | x |
| Send tokenfactory asset to Blacklisted User | It was incorrectly possible for a non-blacklisted user to send minted funds to a blacklisted user. | herod tx bank send user2 [cosmos1qtw23j6y9758juk6mnm937uands7vyxsxkhptu][user1] 100udrachma --fees 500uhero | x |
| Send Tokenfactory Asset from Blacklisted User | Attempted to send a minted asset using a blacklisted account. The client acted correctly and the tx was never broadcast. | herod tx bank send [cosmos1qtw23j6y9758juk6mnm937uands7vyxsxkhptu][user1] [cosmos1zaamktfya6ps3ektdpt7xgka5nu9hx2w8ruthx][user2]100udrachma --fees 500uhero | x |
| Un-Blacklist User | Successfully un-blacklisted a user account using the privileged “Blacklister” account. The Minter was then correctly able to mint assets into the previously blacklisted account.  | herod tx tokenfactory unblacklist [cosmos1qtw23j6y9758juk6mnm937uands7vyxsxkhptu][user1] --fees 500uhero --from black-lister | x |
| Update Pauser | Successfully updated Pauser account using privileged “owner” account. Note: when querying for the pauser address, the address itself is not needed.  | herod tx tokenfactory  update-pauser [cosmos1earswrsr7xhl8mgkglvfvcxqxqsflpu9cvun93][pauser] --from owner --fees 500uhero | x |
| Pause | Successfully paused chain using the privileged “pauser” account | herod tx tokenfactory pause --fees 500uhero --from pauser | x |
| Mint Asset While Paused | The Minter correctly failed to mint 100 udrachma into a user account because the chain was paused. | herod tx tokenfactory mint [cosmos1qtw23j6y9758juk6mnm937uands7vyxsxkhptu][user1] 100udrachma --from minter --fees 500uhero | x |
| Send Tokenfactory Asset While Paused | Attempted to send a minted asset while the chain was paused. The client acted correctly and the tx was never broadcast | herod tx bank send user1 [cosmos1zaamktfya6ps3ektdpt7xgka5nu9hx2w8ruthx][user2] 100udrachma --fees 500uhero | x |
| Send Non-Tokenfactory Asset While Paused | Attempted to send a non-tokenfactory asset while the chain was paused. The client acted correctly and the tx was never broadcast. | herod tx bank send user1 [cosmos1zaamktfya6ps3ektdpt7xgka5nu9hx2w8ruthx][user2] 100uhero --fees 500uhero | x |
| Un-Pause | Successfully un-paused chain using the privileged “pauser” account | herod tx tokenfactory unpause --fees 500uhero --from pauser | x |



