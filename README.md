# bonusly-review-bot
Simple bot that automatically rewards your colleagues Bonusly points for reviewing your pull requests.

### How does it work?
The bot scans your open pull requests, identifies recent reviewers, and automatically sends them Bonusly points on your behalf.

## Setup
- You will need to generate 2 Access tokens and set them in your environment variables.

  - **GitHub Personal Access Token**: Create a GitHub PAT from [here](https://github.com/settings/tokens). Make sure to grant **all** `read` permissions to access complete pull request information. _After token is created, click on 'Conifgure SSO' and authorize your organization._ 
  - **Bonusly Access Token**: Create a Bonusly access token from [here](https://bonus.ly/api_keys/new).

- Run `main.go`
