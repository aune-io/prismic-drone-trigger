# Prismic Drone trigger
Trigger a Drone pipeline when receiving a webhook from Prismic.
Useful if you are using Drone do build a static website (based for instance on GatsbyJS or Hugo) and want to trigger a build every time there is an update on Prismic.

## Usage
Run the process or container with the following environment variables:
* `DRONE_HOST`: The Drone host.
* `DRONE_TOKEN`: The Drone token.
* `REPO_OWNER`: The repository owner.
* `REPO_NAME`: The repository name.
* `REPO_BRANCH`: The repository branch.
* `PRISMIC_SECRET`: The webhook secret set on Prismic.
* `HTTP_ADDRESS`: The listening IP for the webserver, specify only to bind to a specific address (optional).
* `HTTP_PORT`: The listening port for the webserver, specify if different from 80 (optional).
* `HTTP_ROUTE`: The path for the webhook, `/handle` by default (optional).

Notes:
* A build for the specified repository and branch needs to exist, as Drone only support a build restart. This is checked at the process startup.
