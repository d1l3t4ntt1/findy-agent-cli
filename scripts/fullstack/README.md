# Agency Setup for Local Development

## Description

In case you do not have cloud installation of Findy Agency available, you can
setup needed services locally and develop your application against a local Findy
Agency. This document describes how to set up Findy Agency service containers on
your local computer.

The setup uses agency internal file ledger, intended only for testing during
development. This setup does not suit for testing inter-agency communication
even though it is possible to set one up using a common indy-plenum ledger.

## Prerequisites

- [Docker](https://www.docker.com/products/docker-desktop)
- [findy-agent-installation](https://github.com/findy-network/findy-agent-cli#installation)

## Steps

1. Launch backend services with

   ```sh
   make pull-up
   ```

   This will pull the latest versions of the needed docker images. Later on when
   launching the backend you can use `make up` if there is no need to fetch the
   latest images.

   It will take a short while for all the services to start up. Logs from all of
   the started services are printed to the console. `<CTRL>+C` stops the
   containers.

   The script will create a folder called `.data` where all the data of the
   services are stored during execution. If there is no need for the test data
   anymore, `make clean` will remove all the generated data and allocated
   resources.

1. Build playground environment with CLI tool. It's usually
   good idea to have some test data at the backend before UI development or
   application logic itself. Now, when your whole stack is running thanks to
   step one you can easily play with it from the command line.

   To install `findy-agent-cli` execute the following:
   ```console
   /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/findy-network/findy-agent-cli/HEAD/install.sh)"
   ```
   It will install the one binary which is only that's needed in
   `./bin/findy-agent-cli` where you can move it to your path, create alias for
   it, setup auto-completion, etc. More information about it can be found from
   [here](https://github.com/findy-network/findy-agent-cli#installation).

   To make use of `findy-agent-cli` there is a helper script to setup the CLI
   environment. Enter the following command:
   ```console
   source ./setup-cli-env.sh 
   ```
   That will setup all the needed environment variables for CLI configuration
   for the currently running environment. Most importantly it creates a new
   master key for your CLI FIDO2 authenticator. If you want to keep your
   development environment between restarts you should persist that key by
   copying it to your environment variables. The key is in env `FCLI_KEY` after
   running the setup script. The setup generates the `use-key.sh` script for
   your convenient as well. Add `source use-key.sh` to your boot files for
   example.

   Next time you run the `./setup-cli-env.sh` it won't create a new key *if it
   founds the existing one* i.e. you have sourced the `use-key.sh` script.

   *Tip* Enter following commands:
   ```console
   alias cli=findy-agent-cli 
   . <(findy-agent-cli completion bash | sed 's/findy-agent-cli/cli/g')
   ```

   You should enter the following after you have installed the working
   `findy-agent-cli`:
   ```console
   export FCLI=<your-name-for-binary>
   ```
   That's for the helper scrips used in this directory and referenced here as
   well.

   **Admin Operations**

   After environment setup you can see what your configuration is by executing
   the following helper script:
   ```console
   admin/cli-env
   ```
   It will output all of the `findy-agent-cli` env configurations currently set.
   To check one specific variable enter: `admin/cli-env KEY` for example.

   To register your CLI authenticator for direct communication to Findy Agency
   enter the following commands:
   ```console
   source admin/register
   source admin/login
   ```
   Later the login is all what is needed. After successful login you can enter
   commands like:
   ```console
   $FCLI agency count         # get status of the clould agents
   $FCLI agency logging -L=5  # set login level of the core agency 
   ```

   **On-board Alice and Bob**
   ```console
   source alice/register
   source bob/register
   ```
   You can play each of them by entering for example following:
   ```console
   source alice/login
   $FCLI agent ping
   ```

   **Alice invites Bob to connect**

   ```console
   export FCLI_CONN_ID=`alice/invitation | bob/connect`
   ```
   Or on macOS could be convenient to have it in clipboard as well:
   ```console
   alice/invitation | bob/connect | pbcopy && export FCLI_CONN_ID=`pbpaste`
   ```

   Now you have the connection ID (pairwise ID) in the environment variable and
   you could test that with the commands:
   ```console
   source alice/login
   $FCLI connection trustping
   ```
   Which means that Alice's end of the connection calls Aries's trustping
   protocol and Bob's cloud agent responses it.

   Before entering previous commands you could open a second terminal window and
   execute following:
   ```console
   source ./use-key.sh
   source ./setup-cli-env.sh
   source bob/login
   export FClI_CONN_ID="<perviously defined conn id here>"
   $FCLI agent listen
   ```
   You should now receive a notification of the trustping protocol.

   **Alice sends text message to Bob**

   First in the Bob's terminal stop the previous listening with C-c and enter
   the following:
   ```console
   $FCLI bot read
   ```
   Go to the Alice's terminal and enter the commands:
   ```console
   source alice/login
   echo "Hello Bob! Alice here." | $FCLI bot chat
   ```
   The Bob's terminal should output Alice's welcoming messages. To stop Bob's
   `bot read` command just press C-c.

   More samples and guides can be found from
   [Findy Wallet](https://github.com/findy-network/findy-wallet-pwa).
