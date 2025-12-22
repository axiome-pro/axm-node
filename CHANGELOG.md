# Axiome Chain — Release Notes v2.0.0
*(compared to v1.0.4)*

**Release date:** 2025-12-22

This release introduces new protocol-level modules and functionality.  
The upgrade requires a planned **in-place state upgrade** and adds new KV stores.  
Nodes must prepare for the upgrade via a governance **Software Upgrade** proposal named **v2.0.0**.

---

## Summary
- Added the **CosmWasm** module (WASM smart-contract execution).
- Added **Authz** and **Feegrant** modules from the Cosmos SDK.
- Added new KV stores: **wasm**, **authz**, **feegrant** (store migration required at upgrade height).
- CosmWasm parameters are initialized with default values in the upgrade handler.
- New CLI commands for smart-contract interaction (`axmd tx/query wasm …`), authorization delegation (authz), and fee grants (feegrant).

---

## Major Changes

### 1) CosmWasm
- Integrated the **wasmd** module: imports, keeper, and `AppModule` wired into the application.
- Enabled the following CosmWasm features:
  - `iterator`
  - `staking`
  - `distribution`
  - `cosmwasm_1_1`
  - `cosmwasm_1_2`
  - `cosmwasm_1_3`
  - `cosmwasm_1_4`
- In the **v2.0.0 upgrade handler**, CosmWasm parameters are set to default values (`wasmtypes.DefaultParams()`).

**Files:**
- `app/app.go` — store and module wiring, `StoreUpgrades` configuration for v2.0.0.
- `app/upgrades.go` — v2.0.0 upgrade handler, CosmWasm parameter initialization.
- `x/wasm/*` — module and Keeper initialization.

---

### 2) Authz and Feegrant
- Added modules:
  - **Authz** — delegation of message execution permissions.
  - **Feegrant** — grants for transaction fee payments.
- KV stores added for both modules and migrations enabled.

---

### 3) Store Migrations and Load Order
- When an upgrade plan named **v2.0.0** is present, the application adds new KV stores:
  - `wasm`
  - `authz`
  - `feegrant`
- The upgraded store is loaded via `UpgradeStoreLoader`.

---

## Compatibility and Protocol Changes
- A coordinated network upgrade at a fixed height via governance is required.
- This is a **state-breaking upgrade**:
  - new KV stores are added,
  - new transaction execution logic is enabled.
- New RPC and CLI endpoints become available for CosmWasm, Authz, and Feegrant.

---

## Upgrade Procedure for Validators and Node Operators

### 1. Vote (Governance)
- Vote for and pass a **Software Upgrade** proposal with:
  - **Name:** `v2.0.0`
  - **Upgrade height:** specified in the proposal

### 2. Before Upgrade Height
- Build or install the **v2.0.0** binary.
- Stop the node at the upgrade height  
  (or let it halt automatically if auto-halt is enabled).

### 3. At Upgrade Height
- Start the new **v2.0.0** binary.
- The node will:
  - run migrations,
  - add new KV stores: `wasm`, `authz`, `feegrant`.

### 4. After Upgrade
- Wait for synchronization.
- Verify API and CLI functionality.

---

## Post-Upgrade Checks
```bash
axmd q wasm params
axmd q wasm libwasmvm-version
axmd q authz grants
axmd q feegrant grants
```

---

## New CLI Capabilities

### CosmWasm
```bash
axmd tx wasm store <file.wasm> ...
axmd tx wasm instantiate <code-id> '{...}' ...
axmd tx wasm execute <contract-address> '{...}' ...
axmd q wasm contract-state smart <contract-address> '{...}'
```

### Authz
```bash
axmd tx authz ...
axmd q authz ...
```

### Feegrant
```bash
axmd tx feegrant ...
axmd q feegrant ...
```

> **Note:** The exact set of commands may depend on the Cosmos SDK and wasmd versions shipped with this release.

---

## Changes for Developers and Integrators
- New CosmWasm protobuf types and services are enabled.
- OpenAPI / Swagger generation includes additional WASM-related schemas and endpoints.
- Application initialization (`app/app.go`) respects upgrade info.
- Local upgrades must use the upgrade name **v2.0.0**.

---

## Backward Compatibility
- Existing module behavior is preserved.
- New modules are added as extensions.
- A one-time store migration is required.

---

## Acknowledgements
Thanks to all community members and contributors who helped prepare the **v2.0.0** release.
