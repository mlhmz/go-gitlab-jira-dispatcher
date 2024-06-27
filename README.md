# Go GitLab Jira Dispatcher

This is a small project to synchronize GitLab Merge Requests with your JIRA Board.

The tool is configured with a webhook url from gitlab with the merge request permissions and a JIRA API Key. See [GitLab - Create a Webhook](https://docs.gitlab.com/ee/user/project/integrations/webhooks.html#create-a-webhook) for creating a webhook url.

## Technologies

* [Go](https://go.dev)
* [Fiber Web Framework](https://gofiber.io)
* [GORM ORM (with SQLite)](https://gorm.io)
* [urfave/cli](https://cli.urfave.org)
* [HTMX](https://htmx.org)
* [Bootstrap](https://getbootstrap.com/)

## Requirements

* In order to get any assignments working, any actors in the Merge Request need to set their Mail to public 

## Environment Variables

The applications urls and secrets are managed with a env file, it can be created in the application directory, get exported in e.g. the bash shell, or defined in the `docker-compose.yml`

|Identifier|Description|
|----------|-----------|
|JIRA_URL|The url of the jira instance|
|JIRA_API_TOKEN|The api token of the jira user that should be used|
|TOKEN_SIGNATURE|The signature of the jwt for user sessions|


## How to configure the tool

In order to configure this tool, you need a user that can be created with the application cli under the section "users".

After the user was created you can visit the web ui of this tool with the url `http://<YOUR_IP_ADDRESS/HOST>:<PORT default 3000>/`

You will be greeted by the login screen. After you log into the tool you will see this screen and are able to create a new config with "Create"

![Main Screen](https://github.com/mlhmz/go-gitlab-jira-dispatcher/assets/66556288/296812c2-3777-4f10-918b-2bdebe29a023)

On the Create Screen, you can define a Title (the title has no function except to help you recognize your config), the allowed JIRA-Projects (comma seperated) and the Transition IDs for the Events that GitLab is basically sending)

![Create screen](https://github.com/mlhmz/go-gitlab-jira-dispatcher/assets/66556288/3918a91d-77a7-4cb2-82bf-f095771d6713)

after creating the config you will see the uuid, that you can use the webhook with. You can also update or delete the new webhook

![Webhook detail screen](https://github.com/mlhmz/go-gitlab-jira-dispatcher/assets/66556288/dd09aee6-24a2-436c-8a9c-e39ba871e9be)

Finally you can add the webhook into GitLab with the url `http://<YOUR_IP_ADDRESS/HOST>:<PORT default 3000>/webhook/<CONFIG_UUID>`

## Workflow

w.i.p
