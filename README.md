# tfarm

tfarm consists of two parts: `tfarmd` and `tfarm`.

`tfarmd` is a client-side runtime that manages tunnels. It is a daemon that runs in the background and listens for incoming connections from `tfarm`. It is responsible for creating and managing tunnels.

`tfarm` is a cli for interacting with `tfarmd`. It is used to create and manage tunnels.
